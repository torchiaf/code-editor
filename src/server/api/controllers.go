package api

import (
	"net/http"

	authentication "server/authentication"
	kube "server/kubernetes"
	model "server/models"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func Login(c *gin.Context) {

	var auth model.Auth

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, path, err := authentication.LoginCheck(auth)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"path":  path,
	})
}

func Create(c *gin.Context) {
	_, err := kube.StartEditor() // TODO await
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code-server - Cannot create UI instance"})
		return
	}

	username, err := authentication.ExtractUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code-server - User not found"})
		return
	}

	session, err := authentication.EditorLogin(username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Code-server - Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "created",
		session.Name: session.Value,
	})
}

func Delete(c *gin.Context) {
	kube.StopEditor()
	c.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}
