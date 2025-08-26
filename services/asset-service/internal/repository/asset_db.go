package repository

import (
	"context"
	"errors"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetDBRepo struct {
	db *gorm.DB
}

func NewAssetDBRepo(db *gorm.DB) *AssetDBRepo {
	return &AssetDBRepo{db: db}
}

// -------------------- Folder --------------------

func (r *AssetDBRepo) CreateFolder(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *AssetDBRepo) GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	var folder model.Folder
	if err := r.db.WithContext(ctx).First(&folder, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *AssetDBRepo) UpdateFolder(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).
		Model(&model.Folder{}).
		Where("id = ?", folder.ID).
		Updates(folder).Error
}

func (r *AssetDBRepo) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Folder{}, id).Error
}

func (r *AssetDBRepo) GetOwnedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Find(&folders).Error
	return folders, err
}

func (r *AssetDBRepo) GetSharedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.db.WithContext(ctx).
		Joins("JOIN folder_shares ON folder_shares.folder_id = folders.id").
		Where("folder_shares.user_id = ?", userID).
		Where("folders.deleted_at IS NULL").
		Find(&folders).Error
	return folders, err
}

// -------------------- Note --------------------

func (r *AssetDBRepo) CreateNote(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *AssetDBRepo) GetNoteByID(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	var note model.Note
	if err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *AssetDBRepo) UpdateNote(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).
		Model(&model.Note{}).
		Where("id = ?", note.ID).
		Updates(note).Error
}

func (r *AssetDBRepo) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Note{}, id).Error
}

func (r *AssetDBRepo) GetNotesByFolderID(ctx context.Context, folderID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note
	err := r.db.WithContext(ctx).
		Where("folder_id = ?", folderID).
		Where("notes.deleted_at IS NULL").
		Find(&notes).Error
	return notes, err
}

func (r *AssetDBRepo) GetOwnedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", userID).
		Find(&notes).Error
	return notes, err
}

func (r *AssetDBRepo) GetSharedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	var notes []*model.Note
	err := r.db.WithContext(ctx).
		Joins("JOIN note_shares ON note_shares.note_id = notes.id").
		Where("note_shares.user_id = ?", userID).
		Where("notes.deleted_at IS NULL").
		Find(&notes).Error
	return notes, err
}

// -------------------- Sharing --------------------

func (r *AssetDBRepo) ShareFolder(ctx context.Context, folderID, sharedByID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error {
	shares := make([]*model.FolderShare, 0, len(userIDs))

	for _, uid := range userIDs {
		shares = append(shares, &model.FolderShare{
			ID:         uuid.New(),
			FolderID:   folderID,
			UserID:     uid,
			Access:     access,
			SharedByID: sharedByID,
		})
	}

	// Bulk insert
	return r.db.WithContext(ctx).Create(&shares).Error
}

func (r *AssetDBRepo) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("folder_id = ? AND user_id = ?", folderID, userID).Delete(&model.FolderShare{}).Error
}

func (r *AssetDBRepo) ShareNote(ctx context.Context, noteID, sharedByID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error {
	shares := make([]*model.NoteShare, 0, len(userIDs))

	for _, uid := range userIDs {
		shares = append(shares, &model.NoteShare{
			ID:         uuid.New(),
			NoteID:     noteID,
			UserID:     uid,
			Access:     access,
			SharedByID: sharedByID,
		})
	}

	// Bulk insert
	return r.db.WithContext(ctx).Create(&shares).Error
}

func (r *AssetDBRepo) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("note_id = ? AND user_id = ?", noteID, userID).Delete(&model.NoteShare{}).Error
}

// -------------------- Get single share --------------------

func (r *AssetDBRepo) GetNoteShare(ctx context.Context, userID, noteID uuid.UUID) (*model.NoteShare, error) {
	var noteShare model.NoteShare
	err := r.db.WithContext(ctx).
		Where("note_id = ? AND user_id = ?", noteID, userID).
		First(&noteShare).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &noteShare, nil
}

func (r *AssetDBRepo) GetFolderShare(ctx context.Context, userID, folderID uuid.UUID) (*model.FolderShare, error) {
	var folderShare model.FolderShare
	err := r.db.WithContext(ctx).
		Where("folder_id = ? AND user_id = ?", folderID, userID).
		First(&folderShare).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &folderShare, nil
}
