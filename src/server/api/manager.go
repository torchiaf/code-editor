package api

import (
	"net/http"

	kube "server/kubernetes"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func StartEditor(c *gin.Context) {
	kube.StartEditor()
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
	})
}

func StopEditor(c *gin.Context) {
	kube.StopEditor()
	c.JSON(http.StatusOK, gin.H{
		"status": "stopped",
	})
}

func GetPods(c *gin.Context) {
	pods := kube.GetPods()
	c.JSON(200, pods)
}
