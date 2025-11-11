package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonHttpResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
