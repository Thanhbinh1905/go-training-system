package service

import (
	"context"
	"errors"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type FolderService interface {
	CreateFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderInput) error
	GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	UpdateFolder(ctx context.Context, userID, folderID uuid.UUID, input *dto.UpdateFolderInput) error
	DeleteFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) error
}

type folderService struct {
	folderRepo repository.FolderRepository
}

func NewFolderService(folderRepo repository.FolderRepository) FolderService {
	return &folderService{folderRepo}
}

func (s *folderService) CreateFolder(ctx context.Context, userID uuid.UUID, input *dto.CreateFolderInput) error {
	folder := &model.Folder{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
	}
	return s.folderRepo.Create(ctx, folder)
}

func (s *folderService) GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	// TODO: quyền xem
	return s.folderRepo.GetByID(ctx, id)
}

func (s *folderService) UpdateFolder(ctx context.Context, userID, folderID uuid.UUID, input *dto.UpdateFolderInput) error {
	existing, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return err
	}

	if existing.OwnerID != userID {
		return errors.New("forbidden: not the owner")
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Description != nil {
		existing.Description = input.Description
	}

	return s.folderRepo.Update(ctx, existing)
}

func (s *folderService) DeleteFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) error {
	existing, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return err
	}

	if existing.OwnerID != userID {
		return errors.New("forbidden: not the owner")
	}

	return s.folderRepo.Delete(ctx, folderID)
}
