package dto

import (
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/model"
	"github.com/google/uuid"
)

type CreateUserInput struct {
	Username string
	Email    string
	Password string
	Role     model.UserRole
}

type UserResponse struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Role      model.UserRole
	CreatedAt time.Time
}

type UserPaginationInput struct {
	Limit  *int32
	Offset *int32
	Role   *model.UserRole
}

type PaginatedUsersResponse struct {
	Users  []*UserResponse
	Total  int32
	Limit  int32
	Offset int32
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type TokenVerifyInput struct {
	Token string `json:"token" binding:"required"`
}

type UserClaims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

type TokenVerifyResponse struct {
	IsValid     bool        `json:"is_valid"`
	User        *UserClaims `json:"user,omitempty"`
	ExpiredAt   *time.Time  `json:"expired_at,omitempty"`
	ErrorReason string      `json:"error_reason,omitempty"`
}
