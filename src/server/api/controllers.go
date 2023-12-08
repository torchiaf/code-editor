package api

import (
	"fmt"
	"net/http"
	"time"

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

	_, err := kube.ScaleCodeServer(user, 1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code-server - Cannot enable UI instance"})
		return
	}

	time.Sleep(2000 * time.Millisecond)

	session, err := authentication.EditorLogin(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Code-server - Unauthorized"})
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

	kube.ScaleCodeServer(user, 0)
	c.JSON(http.StatusOK, gin.H{
		"status": "disabled",
	})
}

func (vw View) Config(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	var vwConfig model.ViewConfig
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

	label := fmt.Sprintf("app.code-editor/path=%s", user.Path)

	err := kube.ExecCmdOnPod(label, gitCmd, nil, nil, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Pod Configuration failed; %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "config saved",
		"query-param": fmt.Sprintf("folder=/git/%s", git.Repo),
	})
}
