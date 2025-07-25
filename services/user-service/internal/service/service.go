package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Thanhbinh1905/go-training-system/pkg/apperror"
	"github.com/Thanhbinh1905/go-training-system/pkg/hash"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/token"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, input *dto.CreateUserInput) (*dto.UserResponse, error)
	Login(ctx context.Context, input *dto.LoginInput) (*dto.AuthResponse, error)
	Logout(ctx context.Context)
	FetchUsers(ctx context.Context, input *dto.UserPaginationInput) (*dto.PaginatedUsersResponse, error)
	ValidateToken(input *dto.TokenVerifyInput) (*dto.TokenVerifyResponse, error)
	User(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}

type userService struct {
	repo         repository.UserRepository
	tokenManager token.TokenManager
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, input *dto.CreateUserInput) (*dto.UserResponse, error) {
	isTaken, err := s.repo.IsEmailTaken(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if isTaken {
		return nil, apperror.ErrEmailTaken
	}

	hashedPassword, err := hash.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	arg := &model.User{
		ID:           uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         input.Role,
		CreatedAt:    time.Now(),
	}

	err = s.repo.Create(ctx, arg)
	if err != nil {
		return nil, err
	}

	res := &dto.UserResponse{
		ID:        arg.ID,
		Username:  arg.Username,
		Email:     arg.Email,
		Role:      arg.Role,
		CreatedAt: arg.CreatedAt,
	}

	return res, err
}

func (s *userService) Login(ctx context.Context, input *dto.LoginInput) (*dto.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, apperror.ErrInvalidLogin
	}

	if !hash.CheckHash(input.Password, user.PasswordHash) {
		return nil, apperror.ErrInvalidLogin
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, apperror.ErrGenerateAccessTokenFail
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.ErrGenerateRefreshTokenFail
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *userService) Logout(ctx context.Context) {
	panic(fmt.Errorf("not implemented: Logout"))
}

func (s *userService) FetchUsers(ctx context.Context, input *dto.UserPaginationInput) (*dto.PaginatedUsersResponse, error) {
	limit := int32(10)
	offset := int32(0)
	if input.Limit != nil {
		limit = *input.Limit
	}
	if input.Offset != nil {
		offset = *input.Offset
	}

	paginatedUsers, err := s.repo.Fetch(ctx, input.Role, limit, offset)
	if err != nil {
		return nil, err
	}

	var userResponses []*dto.UserResponse
	for _, user := range paginatedUsers.Users {
		userResponses = append(userResponses, &dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}

	return &dto.PaginatedUsersResponse{
		Users:  userResponses,
		Total:  paginatedUsers.Total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *userService) User(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) ValidateToken(input *dto.TokenVerifyInput) (*dto.TokenVerifyResponse, error) {
	claims, err := s.tokenManager.VerifyAccessToken(input.Token)
	if err != nil {
		return &dto.TokenVerifyResponse{
			IsValid:     false,
			ErrorReason: err.Error(),
		}, nil
	}

	var exp *time.Time
	if claims.ExpiresAt != nil {
		exp = &claims.ExpiresAt.Time
	}

	return &dto.TokenVerifyResponse{
		IsValid: true,
		User: &dto.UserClaims{
			ID:   claims.UserID,
			Role: claims.Role,
		},
		ExpiredAt: exp,
	}, nil
}
