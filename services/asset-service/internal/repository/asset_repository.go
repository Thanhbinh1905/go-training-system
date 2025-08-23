package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
)

type AssetRepo interface {
	// Folder
	CreateFolder(ctx context.Context, folder *model.Folder) error
	GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	UpdateFolder(ctx context.Context, folder *model.Folder) error
	DeleteFolder(ctx context.Context, id uuid.UUID) error
	GetOwnedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error)
	GetSharedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error)

	// Note
	CreateNote(ctx context.Context, note *model.Note) error
	GetNoteByID(ctx context.Context, id uuid.UUID) (*model.Note, error)
	UpdateNote(ctx context.Context, note *model.Note) error
	DeleteNote(ctx context.Context, id uuid.UUID) error
	GetNotesByFolderID(ctx context.Context, folderID uuid.UUID) ([]*model.Note, error)
	GetOwnedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error)
	GetSharedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error)

	// Sharing
	ShareFolder(ctx context.Context, folderID, sharedByID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error
	RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error

	ShareNote(ctx context.Context, noteID, sharedByID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error
	RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error

	// Get single share records (normalize not found -> nil, nil)
	GetNoteShare(ctx context.Context, userID, noteID uuid.UUID) (*model.NoteShare, error)
	GetFolderShare(ctx context.Context, userID, folderID uuid.UUID) (*model.FolderShare, error)
}

type CachedAssetRepo struct {
	dbRepo *AssetDBRepo
	cache  *AssetCache
}

func NewCachedAssetRepo(dbRepo *AssetDBRepo, cache *AssetCache) AssetRepo {
	return &CachedAssetRepo{
		dbRepo: dbRepo,
		cache:  cache,
	}
}

func (r *CachedAssetRepo) CreateFolder(ctx context.Context, folder *model.Folder) error {
	if err := r.dbRepo.CreateFolder(ctx, folder); err != nil {
		return err
	}
	// write-through cache
	return r.cache.SetFolder(ctx, folder)
}

func (r *CachedAssetRepo) GetFolderByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	// check cache first
	if folder, err := r.cache.GetFolder(ctx, id); err != nil {
		return nil, err
	} else if folder != nil {
		return folder, nil
	}

	// fallback DB
	folder, err := r.dbRepo.GetFolderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// update cache
	_ = r.cache.SetFolder(ctx, folder)
	return folder, nil
}

func (r *CachedAssetRepo) UpdateFolder(ctx context.Context, folder *model.Folder) error {
	if err := r.dbRepo.UpdateFolder(ctx, folder); err != nil {
		return err
	}
	// invalidate or update cache
	return r.cache.SetFolder(ctx, folder)
}

func (r *CachedAssetRepo) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	if err := r.dbRepo.DeleteFolder(ctx, id); err != nil {
		return err
	}
	return r.cache.InvalidateFolder(ctx, id)
}

// For list queries, we could skip caching for simplicity or implement a separate list cache
func (r *CachedAssetRepo) GetOwnedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error) {
	return r.dbRepo.GetOwnedFolders(ctx, userID)
}

func (r *CachedAssetRepo) GetSharedFolders(ctx context.Context, userID uuid.UUID) ([]*model.Folder, error) {
	return r.dbRepo.GetSharedFolders(ctx, userID)
}

// -------------------- Note --------------------

func (r *CachedAssetRepo) CreateNote(ctx context.Context, note *model.Note) error {
	if err := r.dbRepo.CreateNote(ctx, note); err != nil {
		return err
	}
	return r.cache.SetNote(ctx, note)
}

func (r *CachedAssetRepo) GetNoteByID(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	if note, err := r.cache.GetNote(ctx, id); err != nil {
		return nil, err
	} else if note != nil {
		return note, nil
	}

	note, err := r.dbRepo.GetNoteByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = r.cache.SetNote(ctx, note)
	return note, nil
}

func (r *CachedAssetRepo) UpdateNote(ctx context.Context, note *model.Note) error {
	if err := r.dbRepo.UpdateNote(ctx, note); err != nil {
		return err
	}
	return r.cache.SetNote(ctx, note)
}

func (r *CachedAssetRepo) DeleteNote(ctx context.Context, id uuid.UUID) error {
	if err := r.dbRepo.DeleteNote(ctx, id); err != nil {
		return err
	}
	return r.cache.InvalidateNote(ctx, id)
}

func (r *CachedAssetRepo) GetNotesByFolderID(ctx context.Context, folderID uuid.UUID) ([]*model.Note, error) {
	return r.dbRepo.GetNotesByFolderID(ctx, folderID)
}

func (r *CachedAssetRepo) GetOwnedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	return r.dbRepo.GetOwnedNotes(ctx, userID)
}

func (r *CachedAssetRepo) GetSharedNotes(ctx context.Context, userID uuid.UUID) ([]*model.Note, error) {
	return r.dbRepo.GetSharedNotes(ctx, userID)
}

// -------------------- Sharing --------------------

func (r *CachedAssetRepo) ShareFolder(ctx context.Context, folderID, sharedByID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error {
	return r.dbRepo.ShareFolder(ctx, folderID, sharedByID, userIDs, access)
}

func (r *CachedAssetRepo) RevokeFolderShare(ctx context.Context, folderID, userID uuid.UUID) error {
	return r.dbRepo.RevokeFolderShare(ctx, folderID, userID)
}

func (r *CachedAssetRepo) ShareNote(ctx context.Context, noteID, sharedByID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error {
	return r.dbRepo.ShareNote(ctx, noteID, sharedByID, userIDs, access)
}

func (r *CachedAssetRepo) RevokeNoteShare(ctx context.Context, noteID, userID uuid.UUID) error {
	return r.dbRepo.RevokeNoteShare(ctx, noteID, userID)
}

// -------------------- Single Share --------------------

func (r *CachedAssetRepo) GetNoteShare(ctx context.Context, userID, noteID uuid.UUID) (*model.NoteShare, error) {
	return r.dbRepo.GetNoteShare(ctx, userID, noteID)
}

func (r *CachedAssetRepo) GetFolderShare(ctx context.Context, userID, folderID uuid.UUID) (*model.FolderShare, error) {
	return r.dbRepo.GetFolderShare(ctx, userID, folderID)
}
