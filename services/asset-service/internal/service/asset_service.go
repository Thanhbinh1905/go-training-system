package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

// Permission actions
type Action string

const (
	ActionView   Action = "view"
	ActionEdit   Action = "edit"
	ActionDelete Action = "delete"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type AssetService interface {
}

type assetService struct {
	assetRepo  repository.AssetRepo
	userClient client.UserGRPCClient
	teamClient client.TeamGRPCClient
}

func NewAssetService(assetRepo repository.AssetRepo, userClient client.UserGRPCClient, teamClient client.TeamGRPCClient) AssetService {
	return &assetService{assetRepo: assetRepo, userClient: userClient, teamClient: teamClient}
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

	// 2. Check FolderShare
	folderShare, err := s.assetRepo.GetFolderShare(ctx, userID, folder.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.AccessLevelNone, fmt.Errorf("get folder share: %w", err)
	}
	if folderShare != nil {
		if folderShare.Access == model.AccessLevelWrite {
			return model.AccessLevelWrite, nil
		}
		return model.AccessLevelRead, nil
	}

	// 3. Nếu không có quyền
	return model.AccessLevelNone, nil
}

func (s *assetService) checkNotePermission(ctx context.Context, userID uuid.UUID, note *model.Note) (model.AccessLevel, error) {
	// 1. Owner => full quyền
	if note.OwnerID == userID {
		return model.AccessLevelWrite, nil
	}

	// 2. Check NoteShare
	noteShare, err := s.assetRepo.GetNoteShare(ctx, userID, note.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.AccessLevelNone, fmt.Errorf("get note share: %w", err)
	}
	if noteShare != nil {
		if noteShare.Access == model.AccessLevelWrite {
			return model.AccessLevelWrite, nil
		}
		return model.AccessLevelRead, nil
	}

	// 3. Check FolderShare
	folderShare, err := s.assetRepo.GetFolderShare(ctx, userID, note.FolderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.AccessLevelNone, fmt.Errorf("get folder share: %w", err)
	}
	if folderShare != nil {
		if folderShare.Access == model.AccessLevelWrite {
			return model.AccessLevelWrite, nil
		}
		return model.AccessLevelRead, nil
	}

	// 4. Nếu không có quyền
	return model.AccessLevelNone, nil
}

// -------------------- Folder --------------------

func (s *assetService) CreateFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderInput) error {
	folder := &model.Folder{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
	}
	return s.assetRepo.CreateFolder(ctx, folder)
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
	if folderPerm == model.AccessLevelNone {
		return nil, ErrPermissionDenied
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
		return ErrPermissionDenied
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
		return ErrPermissionDenied
	}

	if err := s.assetRepo.DeleteFolder(ctx, folderID); err != nil {
		return fmt.Errorf("delete folder: %w", err)
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
	return s.assetRepo.CreateNote(ctx, note)
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
	if perm == model.AccessLevelNone {
		return nil, ErrPermissionDenied
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
		return ErrPermissionDenied
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
		return ErrPermissionDenied
	}

	if err := s.assetRepo.DeleteNote(ctx, noteID); err != nil {
		return fmt.Errorf("delete note: %w", err)
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
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

func (s *assetService) ShareFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderShareInput) error {
	access := model.AccessLevelRead
	if input.Access != nil {
		access = *input.Access
	}
	return s.assetRepo.ShareFolder(ctx, input.FolderID, userID, input.UserIDs, access)
}

func (s *assetService) ShareNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteShareInput) error {
	access := model.AccessLevelRead
	if input.Access != nil {
		access = *input.Access
	}
	return s.assetRepo.ShareNote(ctx, input.NoteID, userID, input.UserIDs, access)
}

func (s *assetService) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	return s.assetRepo.RevokeFolderShare(ctx, folderID, userID)
}

func (s *assetService) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	return s.assetRepo.RevokeNoteShare(ctx, noteID, userID)
}

// -------------------- Manage Team Asset --------------------

// GET	/teams/:teamId/assets	View all assets that team members own or can access
func (s *assetService) GetTeamAssets(ctx context.Context, userID, teamID uuid.UUID) (*dto.TeamAssetsResponse, error) {

}

// GET	/users/:userId/assets	View all assets owned by or shared with user

func (s *assetService) GetUserAssets(ctx context.Context, requestedID, userID uuid.UUID) (*dto.UserAssetsBlock, error) {

}
