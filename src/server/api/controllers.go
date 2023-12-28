package api

import (
	"fmt"
	"net/http"
	"time"

	"server/authentication"
	"server/editor"
	"server/models"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func Login(c *gin.Context) {

	var auth models.Auth

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := authentication.LoginCheck(auth)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

type ViewI interface {
	Enable()
	Disable()
	Config()
}

type View struct {
}

func (vw View) Enable(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	editor := editor.New(user)

	port, password, err := editor.Create()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code-server - Cannot enable UI instance"})
		return
	}

	time.Sleep(2000 * time.Millisecond)

	session, err := editor.Login(port, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "enabled",
		session.Name: session.Value,
		"path":       user.Path,
	})
}

func (vw View) Disable(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	editor := editor.New(user)

	editor.Destroy(user)
	c.JSON(http.StatusOK, gin.H{
		"status": "disabled",
	})
}

func (vw View) Config(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	editor := editor.New(user)

	var vwConfig models.ViewConfig
	if err := c.ShouldBindJSON(&vwConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	git := vwConfig.Git

	gitCmd := fmt.Sprintf(
		"cd /git && rm -rf * && git clone git@github.com:%s/%s -b %s && cd %s && git checkout %s",
		git.Org,
		git.Repo,
		git.Branch,
		git.Repo,
		git.Commit,
	)

	err := editor.Config(gitCmd)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Pod Configuration failed; %s", err.Error())})
		return
	}

	queryParam := fmt.Sprintf("folder=/git/%s", git.Repo)

	c.JSON(http.StatusOK, gin.H{
		"status":      "config saved",
		"query-param": queryParam,
	})
}
