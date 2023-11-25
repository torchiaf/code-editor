package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func CreateEditor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"created": "ok",
	})
}
