package service

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type SharingService interface {
	ShareFolder(ctx context.Context, share *model.FolderShare) error
	RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error
	ShareNote(ctx context.Context, share *model.NoteShare) error
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

func (s *sharingService) ShareFolder(ctx context.Context, share *model.FolderShare) error {
	return s.folderShareRepo.ShareFolder(ctx, share)
}

func (s *sharingService) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	return s.folderShareRepo.RevokeFolderShare(ctx, folderID, userID)
}

func (s *sharingService) ShareNote(ctx context.Context, share *model.NoteShare) error {
	return s.noteShareRepo.ShareNote(ctx, share)
}

func (s *sharingService) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	return s.noteShareRepo.RevokeNoteShare(ctx, noteID, userID)
}
