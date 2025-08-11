package middleware

import (
	"context"
	"net/http"
	"strings"

	userpb "github.com/Thanhbinh1905/go-training-system/services/asset-service/pb/user"
	"github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userClient userpb.UserServiceClient, allowedRoles ...userpb.Role) gin.HandlerFunc {
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

		resp, err := userClient.VerifyAccessToken(c, &userpb.VerifyTokenRequest{AccessToken: token})
		if err != nil || !resp.IsValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role := resp.UserInfo.GetRole()
		userID := resp.UserInfo.GetUserId()

		// Check role if allowedRoles are specified
		if len(allowedRoles) > 0 {
			allowed := false
			for _, r := range allowedRoles {
				if r == role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized role"})
				return
			}
		}

		ctx := context.WithValue(c.Request.Context(), contextkey.CtxUserIDKey(), userID)
		ctx = context.WithValue(ctx, contextkey.CtxUserRoleKey(), role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
