package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	"github.com/Thanhbinh1905/go-training-system/shared/errors"
	"github.com/Thanhbinh1905/go-training-system/shared/kafka"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	teampb "github.com/Thanhbinh1905/go-training-system/services/asset-service/pb/team"
)

type Action string

const (
	ActionView   Action = "view"
	ActionEdit   Action = "edit"
	ActionDelete Action = "delete"
)

const (
	DefaultConcurrencyLimit = 5
	MapLookupThreshold      = 10
)

type AssetService interface {
	CreateFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderInput) error
	GetFolderByID(ctx context.Context, userID, folderID uuid.UUID) (*dto.FolderBlock, error)
	UpdateFolder(ctx context.Context, userID, folderID uuid.UUID, input *dto.UpdateFolderInput) error
	DeleteFolder(ctx context.Context, userID, folderID uuid.UUID) error

	CreateNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteInput) error
	GetNoteByID(ctx context.Context, userID, noteID uuid.UUID) (*dto.NoteBlock, error)
	UpdateNote(ctx context.Context, userID, noteID uuid.UUID, input *dto.UpdateNoteInput) error
	DeleteNote(ctx context.Context, userID, noteID uuid.UUID) error

	ShareFolder(ctx context.Context, userID, folderID uuid.UUID, input *dto.CreateFolderShareInput) error
	ShareNote(ctx context.Context, userID, noteID uuid.UUID, input *dto.CreateNoteShareInput) error
	RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error
	RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error

	GetTeamAssets(ctx context.Context, requestedID, teamID uuid.UUID) (*dto.TeamAssetsResponse, error)
	GetUserAssets(ctx context.Context, requestedID, userID uuid.UUID) (*dto.UserAssetsBlock, error)
}

type assetService struct {
	assetRepo     repository.AssetRepo
	userClient    client.UserGRPCClient
	teamClient    client.TeamGRPCClient
	kafkaProducer kafka.Producer
}

func NewAssetService(assetRepo repository.AssetRepo, userClient client.UserGRPCClient, teamClient client.TeamGRPCClient, kafkaProducer kafka.Producer) AssetService {
	return &assetService{
		assetRepo:     assetRepo,
		userClient:    userClient,
		teamClient:    teamClient,
		kafkaProducer: kafkaProducer,
	}
}

// CheckFolderOwnership checks if user owns the folder
func (s *assetService) checkFolderOwnership(ctx context.Context, userID, folderID uuid.UUID) (bool, error) {
	folder, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return false, err
	}
	return folder.OwnerID == userID, nil
}

// CheckNoteOwnership checks if user owns the note
func (s *assetService) checkNoteOwnership(ctx context.Context, userID, noteID uuid.UUID) (bool, error) {
	note, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		return false, err
	}
	return note.OwnerID == userID, nil
}

func (s *assetService) checkFolderPermission(ctx context.Context, userID uuid.UUID, folder *model.Folder) (model.AccessLevel, error) {
	// 1. Owner => full quyền
	if folder.OwnerID == userID {
		return model.AccessLevelWrite, nil
	}

	// 1.5. Check Redis ACL first
	if perm, err := s.assetRepo.GetUserAccessFromACL(ctx, folder.ID, userID); err == nil && perm != "" {
		return perm, nil
	}

	// 2. Check FolderShare
	folderShare, err := s.assetRepo.GetFolderShare(ctx, userID, folder.ID)
	if err != nil {
		return "", fmt.Errorf("get folder share: %w", err)
	}
	if folderShare != nil {
		if folderShare.Access == model.AccessLevelWrite {
			return model.AccessLevelWrite, nil
		}
		return model.AccessLevelRead, nil
	}

	// 3. Nếu không có quyền
	return "", nil
}

func (s *assetService) checkNotePermission(ctx context.Context, userID uuid.UUID, note *model.Note) (model.AccessLevel, error) {
	// 1. Owner => full quyền
	if note.OwnerID == userID {
		return model.AccessLevelWrite, nil
	}

	// 2. Check NoteShare
	noteShare, err := s.assetRepo.GetNoteShare(ctx, userID, note.ID)
	if err != nil {
		return "", fmt.Errorf("get note share: %w", err)
	}
	if noteShare != nil {
		if noteShare.Access == model.AccessLevelWrite {
			return model.AccessLevelWrite, nil
		}
		return model.AccessLevelRead, nil
	}

	// 3. Check FolderShare
	folderShare, err := s.assetRepo.GetFolderShare(ctx, userID, note.FolderID)
	if err != nil {
		return "", fmt.Errorf("get folder share: %w", err)
	}
	if folderShare != nil {
		if folderShare.Access == model.AccessLevelWrite {
			return model.AccessLevelWrite, nil
		}
		return model.AccessLevelRead, nil
	}

	// 4. Nếu không có quyền
	return "", nil
}

// -------------------- Folder --------------------

func (s *assetService) CreateFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderInput) error {
	folder := &model.Folder{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
	}

	if err := s.assetRepo.CreateFolder(ctx, folder); err != nil {
		return err
	}

	// Emit FOLDER_CREATED event
	assetEvent := kafka.NewAssetEvent(kafka.AssetEventFolderCreated, "folder", folder.ID.String(), userID.String(), userID.String())
	if err := s.kafkaProducer.PublishAssetEvent(ctx, assetEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish folder created event: %v\n", err)
	}

	return nil
}

func (s *assetService) GetFolderByID(ctx context.Context, userID, folderID uuid.UUID) (*dto.FolderBlock, error) {
	folder, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// Check folder permission
	folderPerm, err := s.checkFolderPermission(ctx, userID, folder)
	if err != nil {
		return nil, err
	}
	if folderPerm == "" {
		return nil, errors.ErrForbidden
	}

	notes, err := s.getNotesByFolderID(ctx, userID, folderID, folderPerm)
	if err != nil {
		return nil, fmt.Errorf("get notes: %w", err)
	}

	res := &dto.FolderBlock{
		Folder:     *folder,
		Permission: folderPerm,
		Notes:      notes,
	}
	return res, nil
}

func (s *assetService) UpdateFolder(ctx context.Context, userID, folderID uuid.UUID, input *dto.UpdateFolderInput) error {
	existing, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	perm, err := s.checkFolderPermission(ctx, userID, existing)
	if err != nil {
		return fmt.Errorf("check permission: %w", err)
	}
	if perm != model.AccessLevelWrite {
		return errors.ErrForbidden
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Description != nil {
		existing.Description = input.Description
	}

	if err := s.assetRepo.UpdateFolder(ctx, existing); err != nil {
		return fmt.Errorf("update folder: %w", err)
	}

	// Emit FOLDER_UPDATED event
	assetEvent := kafka.NewAssetEvent(kafka.AssetEventFolderUpdated, "folder", folderID.String(), existing.OwnerID.String(), userID.String())
	if err := s.kafkaProducer.PublishAssetEvent(ctx, assetEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish folder updated event: %v\n", err)
	}

	return nil
}

func (s *assetService) DeleteFolder(ctx context.Context, userID, folderID uuid.UUID) error {
	existing, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	perm, err := s.checkFolderPermission(ctx, userID, existing)
	if err != nil {
		return fmt.Errorf("check permission: %w", err)
	}
	if perm != model.AccessLevelWrite {
		return errors.ErrForbidden
	}

	if err := s.assetRepo.DeleteFolder(ctx, folderID); err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	// Emit FOLDER_DELETED event
	assetEvent := kafka.NewAssetEvent(kafka.AssetEventFolderDeleted, "folder", folderID.String(), existing.OwnerID.String(), userID.String())
	if err := s.kafkaProducer.PublishAssetEvent(ctx, assetEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish folder deleted event: %v\n", err)
	}

	return nil
}

// -------------------- Note --------------------

func (s *assetService) CreateNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteInput) error {
	note := &model.Note{
		ID:       uuid.New(),
		Title:    input.Title,
		Body:     input.Body,
		FolderID: input.FolderID,
		OwnerID:  userID,
	}

	if err := s.assetRepo.CreateNote(ctx, note); err != nil {
		return err
	}

	// Emit NOTE_CREATED event
	assetEvent := kafka.NewAssetEvent(kafka.AssetEventNoteCreated, "note", note.ID.String(), userID.String(), userID.String())
	if err := s.kafkaProducer.PublishAssetEvent(ctx, assetEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish note created event: %v\n", err)
	}

	return nil
}

func (s *assetService) GetNoteByID(ctx context.Context, userID, noteID uuid.UUID) (*dto.NoteBlock, error) {
	note, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	// Check permission
	perm, err := s.checkNotePermission(ctx, userID, note)
	if err != nil {
		return nil, err
	}
	if perm == "" {
		return nil, errors.ErrForbidden
	}

	res := &dto.NoteBlock{
		Note:       *note,
		Permission: perm,
	}
	return res, nil
}

func (s *assetService) UpdateNote(ctx context.Context, userID, noteID uuid.UUID, input *dto.UpdateNoteInput) error {
	existing, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	// Check permission
	perm, err := s.checkNotePermission(ctx, userID, existing)
	if err != nil {
		return fmt.Errorf("check permission: %w", err)
	}
	if perm != model.AccessLevelWrite {
		return errors.ErrForbidden
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Body != nil {
		existing.Body = input.Body
	}

	if err := s.assetRepo.UpdateNote(ctx, existing); err != nil {
		return fmt.Errorf("update note: %w", err)
	}

	// Emit NOTE_UPDATED event
	assetEvent := kafka.NewAssetEvent(kafka.AssetEventNoteUpdated, "note", noteID.String(), existing.OwnerID.String(), userID.String())
	if err := s.kafkaProducer.PublishAssetEvent(ctx, assetEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish note updated event: %v\n", err)
	}

	return nil
}

func (s *assetService) DeleteNote(ctx context.Context, userID, noteID uuid.UUID) error {
	existing, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	// Check permission
	perm, err := s.checkNotePermission(ctx, userID, existing)
	if err != nil {
		return fmt.Errorf("check permission: %w", err)
	}
	if perm != model.AccessLevelWrite {
		return errors.ErrForbidden
	}

	if err := s.assetRepo.DeleteNote(ctx, noteID); err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	// Emit NOTE_DELETED event
	assetEvent := kafka.NewAssetEvent(kafka.AssetEventNoteDeleted, "note", noteID.String(), existing.OwnerID.String(), userID.String())
	if err := s.kafkaProducer.PublishAssetEvent(ctx, assetEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish note deleted event: %v\n", err)
	}

	return nil
}

func (s *assetService) getNotesByFolderID(ctx context.Context, userID, folderID uuid.UUID, folderPermission model.AccessLevel) ([]*dto.NoteBlock, error) {
	notes, err := s.assetRepo.GetNotesByFolderID(ctx, folderID)
	if err != nil {
		return nil, fmt.Errorf("get notes: %w", err)
	}

	notess := make([]*dto.NoteBlock, 0, len(notes))

	if folderPermission == model.AccessLevelWrite {
		// Folder có quyền write → tất cả notes trong folder write
		for _, note := range notes {
			notess = append(notess, &dto.NoteBlock{
				Note:       *note,
				Permission: model.AccessLevelWrite,
			})
		}
	} else {
		// Folder read-only: check to upgrate note permission
		for _, note := range notes {
			perm := model.AccessLevelRead

			noteShare, err := s.assetRepo.GetNoteShare(ctx, userID, note.ID)
			if err != nil {
				return nil, fmt.Errorf("get note share: %w", err)
			}
			if noteShare != nil && noteShare.Access == model.AccessLevelWrite {
				perm = model.AccessLevelWrite
			}

			notess = append(notess, &dto.NoteBlock{
				Note:       *note,
				Permission: perm,
			})
		}
	}

	return notess, nil
}

// -------------------- Sharing --------------------

func (s *assetService) ShareFolder(ctx context.Context, userID, folderID uuid.UUID, input *dto.CreateFolderShareInput) error {
	access := model.AccessLevelRead
	if input.Access != nil {
		access = *input.Access
	}

	if err := s.assetRepo.ShareFolder(ctx, folderID, userID, input.UserIDs, access); err != nil {
		return err
	}

	// Emit FOLDER_SHARED event for each user
	for _, targetUserID := range input.UserIDs {
		assetEvent := kafka.NewAssetShareEvent(kafka.AssetEventFolderShared, "folder", folderID.String(), userID.String(), userID.String(), targetUserID.String(), string(access))
		if err := s.kafkaProducer.PublishAssetShareEvent(ctx, assetEvent); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to publish folder shared event: %v\n", err)
		}
	}

	return nil
}

func (s *assetService) ShareNote(ctx context.Context, userID, noteID uuid.UUID, input *dto.CreateNoteShareInput) error {
	access := model.AccessLevelRead
	if input.Access != nil {
		access = *input.Access
	}

	if err := s.assetRepo.ShareNote(ctx, noteID, userID, input.UserIDs, access); err != nil {
		return err
	}

	// Emit NOTE_SHARED event for each user
	for _, targetUserID := range input.UserIDs {
		assetEvent := kafka.NewAssetShareEvent(kafka.AssetEventNoteShared, "note", noteID.String(), userID.String(), userID.String(), targetUserID.String(), string(access))
		if err := s.kafkaProducer.PublishAssetShareEvent(ctx, assetEvent); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to publish note shared event: %v\n", err)
		}
	}

	return nil
}

func (s *assetService) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	actorID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	folder, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	if err := s.assetRepo.RevokeFolderShare(ctx, folderID, userID); err != nil {
		return err
	}

	shareEvent := kafka.NewAssetShareEvent(
		kafka.AssetEventFolderUnshared,
		"folder",
		folderID.String(),
		folder.OwnerID.String(),
		actorID.String(),
		userID.String(),
		"",
	)
	if err := s.kafkaProducer.PublishAssetShareEvent(ctx, shareEvent); err != nil {
		fmt.Printf("Failed to publish folder unshared event: %v\n", err)
	}

	return nil
}

func (s *assetService) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	actorID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	note, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	if err := s.assetRepo.RevokeNoteShare(ctx, noteID, userID); err != nil {
		return err
	}

	shareEvent := kafka.NewAssetShareEvent(
		kafka.AssetEventNoteUnshared,
		"note",
		noteID.String(),
		note.OwnerID.String(),
		actorID.String(),
		userID.String(),
		"",
	)
	if err := s.kafkaProducer.PublishAssetShareEvent(ctx, shareEvent); err != nil {
		fmt.Printf("Failed to publish note unshared event: %v\n", err)
	}

	return nil
}

// -------------------- Manage Team Asset --------------------

// GetAllAssetsInTeam optimized version with concurrent processing and efficient DB queries
func (s *assetService) GetAllAssetsInTeam(ctx context.Context, requestedUserID uuid.UUID, teamUsers *teampb.GetUserIDsByTeamIDResponse) ([]*dto.UserAssetsBlock, error) {
	// Combine all user IDs
	allUserIDs := make([]string, 0, len(teamUsers.ManagerIds)+len(teamUsers.MemberIds))
	allUserIDs = append(allUserIDs, teamUsers.ManagerIds...)
	allUserIDs = append(allUserIDs, teamUsers.MemberIds...)

	if len(allUserIDs) == 0 {
		return []*dto.UserAssetsBlock{}, nil
	}

	// Pre-allocate result slice
	assets := make([]*dto.UserAssetsBlock, len(allUserIDs))

	// Check if requested user is manager
	isRequestedUserManager := s.isUserTeamManager(teamUsers.ManagerIds, requestedUserID)

	// Use errgroup for concurrent processing
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5) // Limit concurrent goroutines to avoid overwhelming DB

	for i, userID := range allUserIDs {
		i, userID := i, userID // Capture loop variables
		g.Go(func() error {
			targetUserUUID, err := uuid.Parse(userID)
			if err != nil {
				return fmt.Errorf("%w: failed to parse userID", errors.ErrBadRequest)
			}

			userAssets, err := s.processUserAssets(ctx, requestedUserID, targetUserUUID, isRequestedUserManager)
			if err != nil {
				return err
			}

			assets[i] = userAssets
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *assetService) isUserInTeam(teamUsers *teampb.GetUserIDsByTeamIDResponse, userID uuid.UUID) bool {
	userIDStr := userID.String()

	// Check in managers
	for _, id := range teamUsers.ManagerIds {
		if id == userIDStr {
			return true
		}
	}

	// Check in members
	for _, id := range teamUsers.MemberIds {
		if id == userIDStr {
			return true
		}
	}

	return false
}

// isUserTeamManager checks if user is a team manager using optimized lookup
func (s *assetService) isUserTeamManager(managerIDs []string, userID uuid.UUID) bool {
	if len(managerIDs) > MapLookupThreshold {
		managerMap := make(map[string]bool, len(managerIDs))
		for _, id := range managerIDs {
			managerMap[id] = true
		}
		return managerMap[userID.String()]
	}

	userIDStr := userID.String()
	for _, id := range managerIDs {
		if id == userIDStr {
			return true
		}
	}
	return false
}

// getTeamAssetsForManager gets all team assets when requested by a manager
func (s *assetService) getTeamAssetsForManager(ctx context.Context, teamUsers *teampb.GetUserIDsByTeamIDResponse) ([]*dto.UserAssetsBlock, error) {
	allUserIDs := make([]string, 0, len(teamUsers.ManagerIds)+len(teamUsers.MemberIds))
	allUserIDs = append(allUserIDs, teamUsers.ManagerIds...)
	allUserIDs = append(allUserIDs, teamUsers.MemberIds...)

	if len(allUserIDs) == 0 {
		return []*dto.UserAssetsBlock{}, nil
	}

	assets := make([]*dto.UserAssetsBlock, len(allUserIDs))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(DefaultConcurrencyLimit)

	for i, userIDStr := range allUserIDs {
		i, userIDStr := i, userIDStr
		g.Go(func() error {
			targetUserUUID, err := uuid.Parse(userIDStr)
			if err != nil {
				return fmt.Errorf("failed to parse user ID %s: %w", userIDStr, err)
			}

			folderBlocks, err := s.getUserFoldersByTeamManager(ctx, targetUserUUID)
			if err != nil {
				// Return empty folders thay vì error
				folderBlocks = []*dto.FolderBlock{}
			}

			assets[i] = &dto.UserAssetsBlock{
				UserID:  targetUserUUID,
				Folders: folderBlocks,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return assets, nil
}

// getTeamAssetsForMember gets permitted team assets when requested by a member
func (s *assetService) getTeamAssetsForMember(ctx context.Context, requestedID uuid.UUID, teamUsers *teampb.GetUserIDsByTeamIDResponse) ([]*dto.UserAssetsBlock, error) {
	allUserIDs := make([]string, 0, len(teamUsers.ManagerIds)+len(teamUsers.MemberIds))
	allUserIDs = append(allUserIDs, teamUsers.ManagerIds...)
	allUserIDs = append(allUserIDs, teamUsers.MemberIds...)

	if len(allUserIDs) == 0 {
		return []*dto.UserAssetsBlock{}, nil
	}

	assets := make([]*dto.UserAssetsBlock, 0, len(allUserIDs)) // Use 0 length since we'll filter

	// Use errgroup for concurrent processing
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5) // Limit concurrent operations

	// Use mutex to protect assets slice
	var mu sync.Mutex

	for _, userID := range allUserIDs {
		userID := userID
		g.Go(func() error {
			targetUserUUID, err := uuid.Parse(userID)
			if err != nil {
				return fmt.Errorf("%w: failed to parse userID", errors.ErrBadRequest)
			}

			// Member uses getUserFoldersByTeamMember with permission checks
			folderBlocks, err := s.getUserFoldersByTeamMember(ctx, requestedID, targetUserUUID)
			if err != nil {
				return fmt.Errorf("failed to get folders for user %s: %w", targetUserUUID, err)
			}

			// Only add if user has any accessible folders
			if len(folderBlocks) > 0 {
				userAsset := &dto.UserAssetsBlock{
					UserID:  targetUserUUID,
					Folders: folderBlocks,
				}

				mu.Lock()
				assets = append(assets, userAsset)
				mu.Unlock()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return assets, nil
}

// GetTeamAssets optimized main function
func (s *assetService) GetTeamAssets(ctx context.Context, requestedID, teamID uuid.UUID) (*dto.TeamAssetsResponse, error) {
	// 1. Check team existence
	teamResp, err := s.teamClient.Client.IsTeamExist(ctx, &teampb.GetTeamRequest{
		TeamId: teamID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if !teamResp.Base.Success || !teamResp.IsExist {
		return nil, errors.ErrNotFound
	}

	// 2. Get team users
	teamUsers, err := s.teamClient.Client.GetUserIDsByTeamID(ctx, &teampb.GetUserIDsByTeamIDRequest{
		TeamId: teamID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get team users: %w", err)
	}

	// 3. Check if user is manager and get assets accordingly
	var assets []*dto.UserAssetsBlock
	if s.isUserTeamManager(teamUsers.ManagerIds, requestedID) {
		assets, err = s.getTeamAssetsForManager(ctx, teamUsers)
	} else {
		assets, err = s.getTeamAssetsForMember(ctx, requestedID, teamUsers)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get team assets: %w", err)
	}

	return &dto.TeamAssetsResponse{
		TeamID: teamID,
		Assets: assets,
	}, nil
}

// processUserAssets processes a single user's assets with role-based permissions
func (s *assetService) processUserAssets(ctx context.Context, requestedUserID, targetUserID uuid.UUID, isManager bool) (*dto.UserAssetsBlock, error) {
	var folderBlocks []*dto.FolderBlock
	var err error

	if requestedUserID == targetUserID {
		folderBlocks, err = s.getSelfFolders(ctx, requestedUserID)
	} else if isManager {
		folderBlocks, err = s.getUserFoldersByTeamManager(ctx, targetUserID)
	} else {
		folderBlocks, err = s.getUserFoldersByTeamMember(ctx, requestedUserID, targetUserID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get folders for user %s: %w", targetUserID, err)
	}

	return &dto.UserAssetsBlock{
		UserID:  targetUserID,
		Folders: folderBlocks,
	}, nil
}

func (s *assetService) getSelfFolders(ctx context.Context, requestedID uuid.UUID) ([]*dto.FolderBlock, error) {
	folders, err := s.assetRepo.GetOwnedFolders(ctx, requestedID)
	if err != nil {
		return nil, err
	}

	if len(folders) == 0 {
		return []*dto.FolderBlock{}, nil
	}

	notes, err := s.assetRepo.GetOwnedNotes(ctx, requestedID)
	if err != nil {
		return nil, err
	}

	// Group notes by folderID
	notesByFolder := make(map[uuid.UUID][]*model.Note)
	for _, n := range notes {
		notesByFolder[n.FolderID] = append(notesByFolder[n.FolderID], n)
	}

	// Build DTO
	blocks := make([]*dto.FolderBlock, 0, len(folders))
	for _, f := range folders {
		noteBlocks := make([]*dto.NoteBlock, 0, len(notesByFolder[f.ID]))
		for _, n := range notesByFolder[f.ID] {
			noteBlocks = append(noteBlocks, &dto.NoteBlock{
				Note:       *n,
				Permission: model.AccessLevelWrite, // default permission
			})
		}
		blocks = append(blocks, &dto.FolderBlock{
			Folder:     *f,
			Permission: model.AccessLevelWrite,
			Notes:      noteBlocks,
		})
	}

	return blocks, nil
}

func (s *assetService) getUserFoldersByTeamManager(ctx context.Context, userID uuid.UUID) ([]*dto.FolderBlock, error) {
	folders, err := s.assetRepo.GetOwnedFolders(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(folders) == 0 {
		return []*dto.FolderBlock{}, nil
	}

	notes, err := s.assetRepo.GetOwnedNotes(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Group notes by folderID
	notesByFolder := make(map[uuid.UUID][]*model.Note)
	for _, n := range notes {
		notesByFolder[n.FolderID] = append(notesByFolder[n.FolderID], n)
	}

	// Build DTO
	blocks := make([]*dto.FolderBlock, 0, len(folders))
	for _, f := range folders {
		noteBlocks := make([]*dto.NoteBlock, 0, len(notesByFolder[f.ID]))
		for _, n := range notesByFolder[f.ID] {
			noteBlocks = append(noteBlocks, &dto.NoteBlock{
				Note:       *n,
				Permission: model.AccessLevelRead, // default permission
			})
		}
		blocks = append(blocks, &dto.FolderBlock{
			Folder:     *f,
			Permission: model.AccessLevelRead,
			Notes:      noteBlocks,
		})
	}

	return blocks, nil
}

func (s *assetService) getUserFoldersByTeamMember(ctx context.Context, requestedID, userID uuid.UUID) ([]*dto.FolderBlock, error) {
	folders, err := s.assetRepo.GetOwnedFolders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned folders for user %s: %w", userID, err)
	}

	if len(folders) == 0 {
		return []*dto.FolderBlock{}, nil
	}

	notes, err := s.assetRepo.GetOwnedNotes(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned notes for user %s: %w", userID, err)
	}

	// Group notes by folderID
	notesByFolder := make(map[uuid.UUID][]*model.Note)
	for _, n := range notes {
		notesByFolder[n.FolderID] = append(notesByFolder[n.FolderID], n)
	}

	// Build DTO với error handling tốt hơn
	blocks := make([]*dto.FolderBlock, 0, len(folders))
	for _, f := range folders {
		folderPerm, err := s.checkFolderPermission(ctx, requestedID, f)
		if err != nil {
			continue // Skip folder này thay vì return error
		}

		// Chỉ process nếu có permission
		if folderPerm == "" {
			continue
		}

		noteBlocks := make([]*dto.NoteBlock, 0, len(notesByFolder[f.ID]))
		for _, n := range notesByFolder[f.ID] {
			var notePerm model.AccessLevel
			if folderPerm == model.AccessLevelWrite {
				notePerm = model.AccessLevelWrite
			} else {
				notePerm, err = s.checkNotePermission(ctx, requestedID, n)
				if err != nil {
					continue // Skip note này
				}

				// Skip note nếu không có permission
				if notePerm == "" {
					continue
				}
			}

			noteBlocks = append(noteBlocks, &dto.NoteBlock{
				Note:       *n,
				Permission: notePerm,
			})
		}

		blocks = append(blocks, &dto.FolderBlock{
			Folder:     *f,
			Permission: folderPerm,
			Notes:      noteBlocks,
		})
	}

	return blocks, nil
}

// GET	/users/:userId/assets	View all assets owned by or shared with user
func (s *assetService) GetUserAssets(ctx context.Context, requestedID, userID uuid.UUID) (*dto.UserAssetsBlock, error) {
	var folderBlocks []*dto.FolderBlock
	var err error

	if requestedID == userID {
		folderBlocks, err = s.getSelfFolders(ctx, requestedID)
	} else {
		folderBlocks, err = s.getUserFoldersByTeamMember(ctx, requestedID, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get folders for user %s: %w", userID, err)
	}

	return &dto.UserAssetsBlock{
		UserID:  userID,
		Folders: folderBlocks,
	}, nil
}
