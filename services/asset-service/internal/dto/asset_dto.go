package dto

import (
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
)

type CreateFolderInput struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateFolderInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CreateNoteInput struct {
	Title    string    `json:"title" gorm:"not null"`
	Body     *string   `json:"body"`
	FolderID uuid.UUID `json:"folder_id" gorm:"type:uuid;not null"`
}

type UpdateNoteInput struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type CreateFolderShareInput struct {
	UserIDs []uuid.UUID        `json:"user_ids"`
	Access  *model.AccessLevel `json:"access_level"`
}

type CreateNoteShareInput struct {
	UserIDs []uuid.UUID        `json:"user_ids"`
	Access  *model.AccessLevel `json:"access_level"`
}

type TeamAssetsResponse struct {
	TeamID uuid.UUID          `json:"teamId"`
	Assets []*UserAssetsBlock `json:"assets"`
}

type UserAssetsBlock struct {
	UserID  uuid.UUID      `json:"userId"`
	Folders []*FolderBlock `json:"folders"`
}

type FolderBlock struct {
	model.Folder
	Permission model.AccessLevel `json:"permission"`
	Notes      []*NoteBlock      `json:"notes"`
}

type NoteBlock struct {
	model.Note
	Permission model.AccessLevel `json:"permission"`
}
