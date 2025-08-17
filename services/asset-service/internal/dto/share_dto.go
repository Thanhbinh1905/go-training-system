package dto

import (
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
)

type CreateFolderShareInput struct {
	FolderID uuid.UUID
	UserID   uuid.UUID
	Access   *model.AccessLevel
}

type CreateNoteShareInput struct {
	NoteID uuid.UUID
	UserID uuid.UUID
	Access *model.AccessLevel
}
