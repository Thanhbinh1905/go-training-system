package service

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type FolderService interface {
	CreateFolder(ctx context.Context, input *dto.CreateFolderInput) error
	GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	UpdateFolder(ctx context.Context, input *dto.UpdateFolderInput) error
	DeleteFolder(ctx context.Context, id uuid.UUID) error
}

type folderService struct {
	folderRepo repository.FolderRepository
}

func NewFolderService(folderRepo repository.FolderRepository) FolderService {
	return &folderService{folderRepo}
}

func (s *folderService) CreateFolder(ctx context.Context, input *dto.CreateFolderInput) error {
	folder := &model.Folder{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     input.OwnerID,
	}
	return s.folderRepo.Create(ctx, folder)
}

func (s *folderService) GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	return s.folderRepo.GetByID(ctx, id)
}

func (s *folderService) UpdateFolder(ctx context.Context, input *dto.UpdateFolderInput) error {
	folder := &model.Folder{ID: input.ID}
	if input.Name != nil {
		folder.Name = *input.Name
	}
	if input.Description != nil {
		folder.Description = input.Description
	}
	return s.folderRepo.Update(ctx, folder)
}

func (s *folderService) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	return s.folderRepo.Delete(ctx, id)
}
