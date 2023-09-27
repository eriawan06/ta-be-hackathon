package common

import "github.com/gin-gonic/gin"

func SendSuccess(c *gin.Context, status int, message string, data interface{}) {
	// Send Success Response
	c.JSON(status, gin.H{
		"code":    status,
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SendError(c *gin.Context, status int, message string, err interface{}) {
	// Send Error Response
	c.JSON(status, gin.H{
		"code":    status,
		"success": false,
		"message": message,
		"error":   err,
	})
}
