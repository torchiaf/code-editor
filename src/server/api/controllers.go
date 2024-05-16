package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/torchiaf/code-editor/server/authentication"
	"github.com/torchiaf/code-editor/server/config"
	"github.com/torchiaf/code-editor/server/editor"
	"github.com/torchiaf/code-editor/server/models"
	"github.com/torchiaf/code-editor/server/users"
	"github.com/torchiaf/code-editor/server/utils"

	"github.com/gin-gonic/gin"
)

func ginSuccess(message string, data ...any) gin.H {
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

type UserI interface {
	List()
	Register()
	Unregister()
}

type User struct {
}

type ViewI interface {
	Create()
	Config()
	Destroy()
}

type View struct {
}

func (user User) List(c *gin.Context) {
	u, _ := authentication.GetUser(c)

	if !u.IsAdmin {
		c.JSON(http.StatusConflict, ginError("User unauthorized"))
		return
	}

	var list []models.User

	for _, u := range users.Store.List() {
		if !u.IsAdmin {
			list = append(list, models.User{
				Id:   u.Id,
				Name: u.Name,
			})
		}
	}

	c.JSON(http.StatusOK, ginSuccess("Users list", list))
}

func (user User) Register(c *gin.Context) {

	if !config.Config.Authentication.IsExternal {
		c.JSON(http.StatusBadGateway, ginError("External authentication is not enabled"))
		return
	}

	var ext models.ExternalUserLogin
	if err := c.ShouldBindJSON(&ext); err != nil {
		c.JSON(http.StatusBadRequest, ginError("Missing configs"))
		return
	}

	// TODO add password constraints
	if len(ext.Password) == 0 {
		c.JSON(http.StatusBadRequest, ginError("Password is missing"))
		return
	}

	username, err := authentication.VerifyExternalUser(ext.Token)
	if err != nil {
		c.JSON(http.StatusForbidden, ginError(err.Error()))
		return
	}

	_, ok := users.Store.Get(username)
	if ok {
		c.JSON(http.StatusConflict, ginError(fmt.Sprintf("External Login, user [%s] is already registered", username)))
		return
	}

	u := models.User{
		// TODO generate in helm chart
		Id:       fmt.Sprintf("ext-%s", utils.RandomString(10, "0123456789")),
		Name:     username,
		Password: ext.Password,
	}

	users.Store.Set(u)

	c.JSON(http.StatusOK, ginSuccess("User successfully registered", map[string]interface{}{
		"username": username,
	}))
}

func (user User) Unregister(c *gin.Context) {
	if !config.Config.Authentication.IsExternal {
		c.JSON(http.StatusBadGateway, ginError("External authentication is not enabled"))
		return
	}

	var ext models.ExternalUserLogin
	if err := c.ShouldBindJSON(&ext); err != nil {
		c.JSON(http.StatusBadRequest, ginError("Missing configs"))
		return
	}

	if len(ext.Username) == 0 {
		c.JSON(http.StatusBadRequest, ginError("Username is missing"))
		return
	}

	username, err := authentication.VerifyExternalUser(ext.Token)
	if err != nil {
		c.JSON(http.StatusForbidden, ginError(err.Error()))
		return
	}

	u, ok := users.Store.Get(ext.Username)
	if !ok {
		c.JSON(http.StatusNotFound, ginError("User not found"))
		return
	}

	e := editor.New(c, u)

	store := e.Store()

	details := ""
	if (store != editor.StoreData{} && store.Status == editor.Enabled) {
		if ext.Force {
			// TODO add error handling, destroy could fail
			e.Destroy(u)
			details = ", UI instance destroyed"
		} else {
			c.JSON(http.StatusConflict, ginError(fmt.Sprintf("UI instance is Enabled for user [%s], cannot unregister", ext.Username)))
			return
		}
	}

	users.Store.Del(username)

	c.JSON(http.StatusOK, ginSuccess("User successfully unregistered"+details))
}

func (vw View) List(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	if !user.IsAdmin {
		c.JSON(http.StatusConflict, ginError("User unauthorized"))
		return
	}

	var list []models.View

	for _, user := range users.Store.List() {
		e := editor.New(c, user)

		store := e.Store()

		if len(store.Status) > 0 {
			list = append(list, models.View{
				Id:             e.Id,
				UserId:         user.Id,
				Status:         store.Status,
				Path:           store.Path,
				Password:       store.Password,
				VScodeSettings: store.VScodeSettings,
			})
		}
	}

	c.JSON(http.StatusOK, ginSuccess("View list", list))
}

func (vw View) Get(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	if !user.IsAdmin {
		c.JSON(http.StatusConflict, ginError("User unauthorized"))
		return
	}

	// TODO
	c.JSON(http.StatusOK, ginSuccess("View get", ""))
}

func (vw View) Create(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	e := editor.New(c, user)

	var username string
	if user.IsAdmin {
		username = c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, ginError("Missing 'username' param"))
			return
		}

		user, ok := users.Store.Get(username)
		if !ok {
			c.JSON(http.StatusNotFound, ginError("User not found"))
			return
		}

		e = editor.New(c, user)
	}

	store := e.Store()
	if (store != editor.StoreData{} && store.Status == editor.Enabled) {
		c.JSON(http.StatusForbidden, ginError("View instance already exists"))
		return
	}

	var enableConfig models.EnableConfig
	if err := c.ShouldBindJSON(&enableConfig); err != nil {
		c.JSON(http.StatusBadRequest, ginError("Missing configs"))
		return
	}

	port, err := e.Create(enableConfig)
	if err != nil {
		c.JSON(http.StatusConflict, ginError("Cannot create View instance"))
		return
	}

	time.Sleep(2000 * time.Millisecond)

	session, err := e.Login(port, e.Store().Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ginError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, ginSuccess("View created", map[string]interface{}{
		"viewId":     e.Id,
		session.Name: session.Value,
		"path":       fmt.Sprintf("/code-editor/%s/", e.Store().Path),
	}))
}

func (vw View) Destroy(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	e := editor.New(c, user)

	store := e.Store()

	if (store == editor.StoreData{} || store.Status == editor.Disabled) {
		c.JSON(http.StatusNotFound, ginError("View instance not found"))
		return
	}

	e.Destroy(user)
	c.JSON(http.StatusOK, ginSuccess("View destroyed"))
}

func (vw View) Config(c *gin.Context) {

	user, _ := authentication.GetUser(c)

	editor := editor.New(c, user)

	var vwConfig models.ViewConfig
	if err := c.ShouldBindJSON(&vwConfig); err != nil {
		c.JSON(http.StatusBadRequest, ginError("Missing configs"))
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
		c.JSON(http.StatusForbidden, ginError(fmt.Sprintf("code-server pod configuration failed; %s", err.Error())))
		return
	}

	queryParam := fmt.Sprintf("folder=/git/%s", git.Repo)

	c.JSON(http.StatusOK, ginSuccess("Configurations saved", map[string]interface{}{
		"query-param": queryParam,
	}))
}
