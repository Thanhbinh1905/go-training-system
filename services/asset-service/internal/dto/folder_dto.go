package dto

type CreateFolderInput struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateFolderInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
