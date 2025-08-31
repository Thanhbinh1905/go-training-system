package httpresponse

import (
	"github.com/Thanhbinh1905/go-training-system/shared/errors"
	"github.com/gin-gonic/gin"
)

func RespondError(c *gin.Context, err error) {
	status := errors.StatusCode(err)
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
	c.Abort()
}

func RespondMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"message": msg})
}
