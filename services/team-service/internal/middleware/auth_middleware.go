package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/pb"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/pkg/contextkey"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userClient pb.UserServiceClient, requireManager bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		resp, err := userClient.VerifyAccessToken(c, &pb.VerifyTokenRequest{AccessToken: token})
		if err != nil || !resp.IsValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if requireManager && resp.UserInfo.GetRole() != pb.Role_MANAGER {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not a manager"})
			return
		}
		ctx := context.WithValue(c.Request.Context(), contextkey.CtxUserIDKey(), resp.UserInfo.GetUserId())
		ctx = context.WithValue(ctx, contextkey.CtxUserRoleKey(), resp.UserInfo.GetRole())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
