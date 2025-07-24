package helper

import (
	gqlmodel "github.com/Thanhbinh1905/go-training-system/internal/graph/model"
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

func AuthMutationSuccess(accessToken string, refreshToken string, user *gqlmodel.User) *gqlmodel.AuthMutationResponse {
	return &gqlmodel.AuthMutationResponse{
		Code:         "200",
		Success:      true,
		Message:      "Authentication successful",
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken, // Assuming no refresh token is used
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
