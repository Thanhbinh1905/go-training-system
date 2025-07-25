package service

import (
	"context"

	gqlmodel "github.com/Thanhbinh1905/go-training-system/internal/graph/model"
	"github.com/Thanhbinh1905/go-training-system/internal/model"
	"github.com/Thanhbinh1905/go-training-system/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/pkg/apperror"
	"github.com/Thanhbinh1905/go-training-system/pkg/hash"
)

type UserService interface {
	Register(ctx context.Context, input *gqlmodel.CreateUserInput) (*model.User, error)
	Update(ctx context.Context, userID string, input *gqlmodel.UpdateUserInput) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context, role *gqlmodel.UserType) ([]*model.User, error)
	Login(ctx context.Context, input *gqlmodel.UserInput) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register creates a new user with hashed password
func (s *userService) Register(ctx context.Context, input *gqlmodel.CreateUserInput) (*model.User, error) {
	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Role:     model.UserRole(input.Role),
	}

	// Check if email is already taken
	isTaken, err := s.repo.IsEmailTaken(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if isTaken {
		return nil, apperror.ErrEmailTaken
	}

	// Hash the password before saving
	hashedPassword, err := hash.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashedPassword

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates an existing user's fields
func (s *userService) Update(ctx context.Context, userID string, input *gqlmodel.UpdateUserInput) (*model.User, error) {
	return s.repo.Update(ctx, userID, input)
}

// GetByID returns a user by ID
func (s *userService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.FindByID(ctx, userID)
}

// GetByEmail returns a user by email
func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

// GetAll returns all users (optionally filtered by role)
func (s *userService) GetAll(ctx context.Context, role *gqlmodel.UserType) ([]*model.User, error) {
	var modelRole *model.UserRole
	if role != nil {
		r := model.UserRole(*role)
		modelRole = &r
	}
	return s.repo.FetchAll(ctx, modelRole)
}

// Login authenticates a user by email + password
func (s *userService) Login(ctx context.Context, input *gqlmodel.UserInput) (*model.User, error) {
	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, apperror.ErrInvalidLogin
	}

	if !hash.CheckHash(input.Password, user.PasswordHash) {
		return nil, apperror.ErrInvalidLogin
	}
	return user, nil
}
