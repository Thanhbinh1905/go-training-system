package service

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type SharingService interface {
	ShareFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderShareInput) error
	RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error
	ShareNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteShareInput) error
	RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error
}

type sharingService struct {
	folderShareRepo repository.FolderShareRepository
	noteShareRepo   repository.NoteShareRepository
}

func NewSharingService(
	folderRepo repository.FolderShareRepository,
	noteRepo repository.NoteShareRepository,
) SharingService {
	return &sharingService{
		folderShareRepo: folderRepo,
		noteShareRepo:   noteRepo,
	}
}

func (s *sharingService) ShareFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderShareInput) error {
	share := &model.FolderShare{
		ID:         uuid.New(),
		FolderID:   input.FolderID,
		UserID:     input.FolderID,
		Access:     model.AccessLevelRead,
		SharedByID: userID,
	}
	if input.Access != nil {
		share.Access = *input.Access
	}
	return s.folderShareRepo.ShareFolder(ctx, share)
}

func (s *sharingService) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	return s.folderShareRepo.RevokeFolderShare(ctx, folderID, userID)
}

func (s *sharingService) ShareNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteShareInput) error {
	share := &model.NoteShare{
		ID:         uuid.New(),
		NoteID:     input.NoteID,
		UserID:     input.UserID,
		Access:     model.AccessLevelRead,
		SharedByID: userID,
	}
	if input.Access != nil {
		share.Access = *input.Access
	}
	return s.noteShareRepo.ShareNote(ctx, share)
}

func (s *sharingService) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	return s.noteShareRepo.RevokeNoteShare(ctx, noteID, userID)
}
