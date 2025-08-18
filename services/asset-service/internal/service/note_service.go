package service

import (
	"context"
	"errors"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteInput) error
	GetNote(ctx context.Context, id uuid.UUID) (*model.Note, error)
	UpdateNote(ctx context.Context, userID, noteID uuid.UUID, input *dto.UpdateNoteInput) error
	DeleteNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	GetNotesByFolder(ctx context.Context, folderID uuid.UUID) ([]model.Note, error)
}

type noteService struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &noteService{noteRepo}
}

func (s *noteService) CreateNote(ctx context.Context, userID uuid.UUID, input *dto.CreateNoteInput) error {
	note := &model.Note{
		ID:       uuid.New(),
		Title:    input.Title,
		Body:     input.Body,
		FolderID: input.FolderID,
		OwnerID:  userID,
	}
	return s.noteRepo.Create(ctx, note)
}

func (s *noteService) GetNote(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	// TODO: quyền xem
	return s.noteRepo.GetByID(ctx, id)
}

func (s *noteService) UpdateNote(ctx context.Context, userID, noteID uuid.UUID, input *dto.UpdateNoteInput) error {
	existing, err := s.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}

	if existing.OwnerID != userID {
		return errors.New("forbidden: not the owner")
	}

	if input.Body != nil {
		existing.Body = input.Body
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}

	return s.noteRepo.Update(ctx, existing)
}

func (s *noteService) DeleteNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	existing, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existing.OwnerID != userID {
		return errors.New("forbidden: not the owner")
	}

	return s.noteRepo.Delete(ctx, id)
}

func (s *noteService) GetNotesByFolder(ctx context.Context, folderID uuid.UUID) ([]model.Note, error) {
	return s.noteRepo.GetNotesByFolderID(ctx, folderID)
}
