// internal/service/asset_service.go
package service

import (
	"context"
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type AssetService interface {
	GetUserAssets(ctx context.Context, userID uuid.UUID) ([]model.Folder, []model.Note, error)
	GetTeamAssets(ctx context.Context, teamID uuid.UUID) ([]model.Folder, []model.Note, error)
}

type assetService struct {
	folderRepo      repository.FolderRepository
	noteRepo        repository.NoteRepository
	folderShareRepo repository.FolderShareRepository
	noteShareRepo   repository.NoteShareRepository
	userFetcher     func(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) // dùng gRPC or REST
}

func NewAssetService(
	folderRepo repository.FolderRepository,
	noteRepo repository.NoteRepository,
	folderShareRepo repository.FolderShareRepository,
	noteShareRepo repository.NoteShareRepository,
	userFetcher func(context.Context, uuid.UUID) ([]uuid.UUID, error),
) AssetService {
	return &assetService{
		folderRepo, noteRepo, folderShareRepo, noteShareRepo, userFetcher,
	}
}

func (s *assetService) GetUserAssets(ctx context.Context, userID uuid.UUID) ([]model.Folder, []model.Note, error) {
	ownedFolders, err := s.folderRepo.GetOwnedFolders(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get owned folders: %w", err)
	}

	sharedFolders, err := s.folderRepo.GetSharedFolders(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get shared folders: %w", err)
	}

	ownedNotes, err := s.noteRepo.GetOwnedNotes(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get owned notes: %w", err)
	}

	sharedNotes, err := s.noteRepo.GetSharedNotes(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get shared notes: %w", err)
	}

	// Merge folders
	allFolders := append(ownedFolders, sharedFolders...)
	folders := make([]model.Folder, 0, len(allFolders))
	for _, f := range allFolders {
		if f != nil {
			folders = append(folders, *f)
		}
	}

	// Merge notes
	allNotes := append(ownedNotes, sharedNotes...)
	notes := make([]model.Note, 0, len(allNotes))
	for _, n := range allNotes {
		if n != nil {
			notes = append(notes, *n)
		}
	}

	return folders, notes, nil
}

func (s *assetService) GetTeamAssets(ctx context.Context, teamID uuid.UUID) ([]model.Folder, []model.Note, error) {
	userIDs, err := s.userFetcher(ctx, teamID)
	if err != nil {
		return nil, nil, err
	}

	folders := []model.Folder{}
	notes := []model.Note{}
	for _, userID := range userIDs {
		userFolders, userNotes, _ := s.GetUserAssets(ctx, userID)
		folders = append(folders, userFolders...)
		notes = append(notes, userNotes...)
	}
	return folders, notes, nil
}
