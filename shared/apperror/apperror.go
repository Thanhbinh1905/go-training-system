package apperror

import "errors"

var (
	ErrEmailTaken   = errors.New("email is already taken")
	ErrInvalidLogin = errors.New("invalid email or password")
	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("unauthorized access")
)
