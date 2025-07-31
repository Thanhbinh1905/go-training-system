package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetNotesByFolderID(ctx context.Context, folderID uuid.UUID) ([]model.Note, error)
	GetOwnedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error)
	GetSharedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error)
}

type noteRepo struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepo{db}
}

func (r *noteRepo) Create(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *noteRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	var note model.Note
	if err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepo) Update(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *noteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Note{}).Error
}

func (r *noteRepo) GetNotesByFolderID(ctx context.Context, folderID uuid.UUID) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.WithContext(ctx).Where("folder_id = ?", folderID).Find(&notes).Error
	return notes, err
}

func (r *noteRepo) GetOwnedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Find(&notes).Error
	return notes, err
}

func (r *noteRepo) GetSharedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note
	err := r.db.WithContext(ctx).
		Joins("JOIN note_shares ON note_shares.note_id = notes.id").
		Where("note_shares.user_id = ?", userID).
		Find(&notes).Error
	return notes, err
}
