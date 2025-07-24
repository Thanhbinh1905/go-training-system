package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Thanhbinh1905/go-training-system/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func RequireManagerRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}

func OptionalAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.VerifyToken(tokenStr, secret)
			if err == nil {
				c.Set(ContextUserID, claims.UserID)
			}
		}
		c.Next()
	}
}

// AuthMiddleware extracts and verifies JWT token from Authorization header
func RequiredAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing or malformed"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.VerifyToken(tokenStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: " + err.Error(),
			})
			return
		}

		// Set thông tin người dùng vào context của Gin
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get(ContextUserID)
		role, _ := c.Get(ContextRole)

		// Truyền dữ liệu vào context chuẩn
		ctx := context.WithValue(c.Request.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextRole, role)

		// Gán lại context mới vào request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
