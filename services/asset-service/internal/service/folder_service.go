package service

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type FolderService interface {
	CreateFolder(ctx context.Context, folder *model.Folder) error
	GetFolder(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	UpdateFolder(ctx context.Context, folder *model.Folder) error
	DeleteFolder(ctx context.Context, id uuid.UUID) error
}

type folderService struct {
	folderRepo repository.FolderRepository
}

func NewFolderService(folderRepo repository.FolderRepository) FolderService {
	return &folderService{folderRepo}
}

func (s *folderService) CreateFolder(ctx context.Context, folder *model.Folder) error {
	return s.folderRepo.Create(ctx, folder)
}

func (s *folderService) GetFolder(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	return s.folderRepo.GetByID(ctx, id)
}

func (s *folderService) UpdateFolder(ctx context.Context, folder *model.Folder) error {
	return s.folderRepo.Update(ctx, folder)
}

func (s *folderService) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	return s.folderRepo.Delete(ctx, id)
}
