package middleware

import (
	"context"
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/pb"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	ctxUserIDKey   ctxKey = "userID"
	ctxUserRoleKey ctxKey = "userRole"
)

func AuthMiddleware(userClient pb.UserServiceClient, requireManager bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		resp, err := userClient.VerifyAccessToken(c, &pb.VerifyTokenRequest{AccessToken: token})
		if err != nil || !resp.IsValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if requireManager && resp.UserInfo.GetRole() != pb.Role_MANAGER {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not a manager"})
			return
		}

		// Inject vào context
		ctx := context.WithValue(c.Request.Context(), ctxUserIDKey, resp.UserInfo.GetUserId())
		ctx = context.WithValue(ctx, ctxUserRoleKey, resp.UserInfo.GetRole())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
