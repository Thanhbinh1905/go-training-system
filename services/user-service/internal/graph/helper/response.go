package helper

import (
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/dto"
	gqlmodel "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/graph/model"
)

func NewUserMutationSuccess(user *gqlmodel.User) *gqlmodel.UserMutationResponse {
	msg := "User operation successful"
	return &gqlmodel.UserMutationResponse{
		Code:    "200",
		Success: true,
		Message: &msg,
		User:    user,
	}
}

func NewUserMutationError(code string, message *string, errors []*string) *gqlmodel.UserMutationResponse {
	return &gqlmodel.UserMutationResponse{
		Code:    code,
		Success: false,
		Message: message,
		Errors:  errors,
	}
}

func AuthMutationSuccess(authResponse *dto.AuthResponse) *gqlmodel.AuthMutationResponse {
	userRes := authResponse.User
	createdAtStr := userRes.CreatedAt.Format(time.RFC3339)
	user := &gqlmodel.User{
		ID:        userRes.ID,
		Username:  userRes.Username,
		Email:     userRes.Email,
		Role:      gqlmodel.UserType(userRes.Role),
		CreatedAt: &createdAtStr,
	}
	return &gqlmodel.AuthMutationResponse{
		Code:         "200",
		Success:      true,
		Message:      "Authentication successful",
		AccessToken:  &authResponse.AccessToken,
		RefreshToken: &authResponse.RefreshToken,
		User:         user,
	}
}

func AuthMutationError(code string, message string, errors []*string) *gqlmodel.AuthMutationResponse {
	return &gqlmodel.AuthMutationResponse{
		Code:    code,
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
