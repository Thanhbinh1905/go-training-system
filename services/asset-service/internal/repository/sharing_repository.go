package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderShareRepository interface {
	ShareFolder(ctx context.Context, share *model.FolderShare) error
	RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error
}

type NoteShareRepository interface {
	ShareNote(ctx context.Context, share *model.NoteShare) error
	RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error
}

type shareRepo struct {
	db *gorm.DB
}

func NewFolderShareRepository(db *gorm.DB) FolderShareRepository {
	return &shareRepo{db}
}

func NewNoteShareRepository(db *gorm.DB) NoteShareRepository {
	return &shareRepo{db}
}

func (r *shareRepo) ShareFolder(ctx context.Context, share *model.FolderShare) error {
	return r.db.WithContext(ctx).Create(share).Error
}

func (r *shareRepo) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("folder_id = ? AND user_id = ?", folderID, userID).Delete(&model.FolderShare{}).Error
}

func (r *shareRepo) ShareNote(ctx context.Context, share *model.NoteShare) error {
	return r.db.WithContext(ctx).Create(share).Error
}

func (r *shareRepo) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("note_id = ? AND user_id = ?", noteID, userID).Delete(&model.NoteShare{}).Error
}
