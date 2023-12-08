package api

import (
	"fmt"
	"net/http"
	"time"

	authentication "server/authentication"
	"server/config"
	kube "server/kubernetes"
	model "server/models"
	"server/utils"

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

	v, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	u := v.(model.User)

	found, user := utils.Find(config.Config.Users, "Name", u.Name)
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Code-server - User not found"})
		return
	}

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

	v, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user := v.(model.User)

	kube.ScaleCodeServer(user, 0)
	c.JSON(http.StatusOK, gin.H{
		"status": "disabled",
	})
}

func (vw View) Config(c *gin.Context) {

	v, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user := v.(model.User)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pod Configuration failed"})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "config saved",
		"query-param": fmt.Sprintf("folder=/git/%s", git.Repo),
	})
}
