package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *model.Folder) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	Update(ctx context.Context, folder *model.Folder) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetOwnedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error)
	GetSharedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error)
}

type folderRepo struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepo{db}
}

func (r *folderRepo) Create(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *folderRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	var folder model.Folder
	if err := r.db.WithContext(ctx).First(&folder, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepo) Update(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

func (r *folderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Folder{}).Error
}

func (r *folderRepo) GetOwnedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Find(&folders).Error
	return folders, err
}

func (r *folderRepo) GetSharedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.db.WithContext(ctx).
		Joins("JOIN folder_shares ON folder_shares.folder_id = folders.id").
		Where("folder_shares.user_id = ?", userID).
		Find(&folders).Error
	return folders, err
}
