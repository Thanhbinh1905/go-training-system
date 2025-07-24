package helper

import "github.com/gin-gonic/gin"

func respondError(c *gin.Context, status int, code string, message string, err error) {
	resp := gin.H{
		"code":    code,
		"success": false,
		"message": message,
	}
	if err != nil {
		resp["errors"] = []string{err.Error()}
	}
	c.JSON(status, resp)
}

func respondSuccess(c *gin.Context, status int, code string, message string, data interface{}) {
	resp := gin.H{
		"code":    code,
		"success": true,
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	c.JSON(status, resp)
}
