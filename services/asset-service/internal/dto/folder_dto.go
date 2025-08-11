package dto

import "github.com/google/uuid"

type CreateFolderInput struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id" binding:"required,uuid"`
}

type UpdateFolderInput struct {
	ID          uuid.UUID `json:"id" binding:"required,uuid"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
}
