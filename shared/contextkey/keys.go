package contextkey

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ctxKey string

const (
	userIDKey   ctxKey = "userID"
	userRoleKey ctxKey = "userRole"
)

func UserIDKey() ctxKey   { return userIDKey }
func UserRoleKey() ctxKey { return userRoleKey }

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if userIDStr, ok := ctx.Value(userIDKey).(string); ok && userIDStr != "" {
		return uuid.Parse(userIDStr)
	}
	return uuid.Nil, fmt.Errorf("userID not found in context")
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	if roleStr, ok := ctx.Value(userRoleKey).(string); ok && roleStr != "" {
		return roleStr, nil
	}
	return "", fmt.Errorf("user role not found in context")
}
