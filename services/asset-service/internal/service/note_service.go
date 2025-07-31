package service

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(ctx context.Context, note *model.Note) error
	GetNote(ctx context.Context, id uuid.UUID) (*model.Note, error)
	UpdateNote(ctx context.Context, note *model.Note) error
	DeleteNote(ctx context.Context, id uuid.UUID) error
	GetNotesByFolder(ctx context.Context, folderID uuid.UUID) ([]model.Note, error)
}

type noteService struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &noteService{noteRepo}
}

func (s *noteService) CreateNote(ctx context.Context, note *model.Note) error {
	return s.noteRepo.Create(ctx, note)
}

func (s *noteService) GetNote(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	return s.noteRepo.GetByID(ctx, id)
}

func (s *noteService) UpdateNote(ctx context.Context, note *model.Note) error {
	return s.noteRepo.Update(ctx, note)
}

func (s *noteService) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return s.noteRepo.Delete(ctx, id)
}

func (s *noteService) GetNotesByFolder(ctx context.Context, folderID uuid.UUID) ([]model.Note, error) {
	return s.noteRepo.GetNotesByFolderID(ctx, folderID)
}
