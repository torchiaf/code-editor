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

func ginSuccess(message string, data ...map[string]any) gin.H {
	ret := gin.H{
		"message": message,
	}

	if len(data) > 0 {
		ret["data"] = data[0]
	}

	return ret
}

func ginError(message string) gin.H {
	return gin.H{
		"error": message,
	}
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, ginSuccess("paolo"))
}

func Login(c *gin.Context) {

	var auth models.Auth

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, ginError(err.Error()))
		return
	}

	token, err := authentication.LoginCheck(auth)

	if err != nil {
		c.JSON(http.StatusBadRequest, ginError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, ginSuccess("Successful login", map[string]interface{}{
		"token": token,
	}))
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

	e := editor.New(user)

	store := e.Store()

	if (store != editor.StoreData{} && store.Status == editor.Enabled) {
		c.JSON(http.StatusNotFound, ginError("UI instance is already Enabled"))
		return
	}

	port, err := e.Create()
	if err != nil {
		c.JSON(http.StatusNotFound, ginError("Cannot enable UI instance"))
		return
	}

	time.Sleep(2000 * time.Millisecond)

	session, err := e.Login(port, e.Store().Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ginError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, ginSuccess("View enabled", map[string]interface{}{
		session.Name: session.Value,
		"path":       e.Store().Path,
	}))
}

func (vw View) Disable(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	e := editor.New(user)

	store := e.Store()

	if (store == editor.StoreData{} || store.Status == editor.Disabled) {
		c.JSON(http.StatusNotFound, ginError("UI instance is already Disabled"))
		return
	}

	e.Destroy(user)
	c.JSON(http.StatusOK, ginSuccess("View disabled"))
}

func (vw View) Config(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	editor := editor.New(user)

	var vwConfig models.ViewConfig
	if err := c.ShouldBindJSON(&vwConfig); err != nil {
		c.JSON(http.StatusBadRequest, ginError(err.Error()))
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
		c.JSON(http.StatusBadRequest, ginError(fmt.Sprintf("Pod Configuration failed; %s", err.Error())))
		return
	}

	queryParam := fmt.Sprintf("folder=/git/%s", git.Repo)

	c.JSON(http.StatusOK, ginSuccess("Configurations saved", map[string]interface{}{
		"query-param": queryParam,
	}))
}
