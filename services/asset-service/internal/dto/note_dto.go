package dto

import "github.com/google/uuid"

type CreateNoteInput struct {
	Title    string    `json:"title" gorm:"not null"`
	Body     *string   `json:"body"`
	FolderID uuid.UUID `json:"folder_id" gorm:"type:uuid;not null"`
}

type UpdateNoteInput struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}
