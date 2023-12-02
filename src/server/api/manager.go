package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	config "server/config"
	e "server/error"
	kube "server/kubernetes"
	model "server/models"
	routing "server/routing"
	utils "server/utils"

	"github.com/gin-gonic/gin"
)

var localConfig = config.GetConfig()

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func Login(c *gin.Context) {

	// Read request body
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	auth := model.Auth{}

	err1 := json.Unmarshal([]byte(reqBody), &auth)
	if err1 != nil {
		e.FailOnError(err, "Failed to get Deployment/Scale")
	}

	// Get user route
	found, route := utils.Find(localConfig.Routes, "Name", auth.User)

	if !found {
		c.JSON(http.StatusBadRequest, "Invalid User")
		return
	}

	// code-server login endpoint
	loginUrl := ""

	if localConfig.IsDev {
		loginUrl = fmt.Sprintf("http://localhost/code-editor/%s/login", route.Path)
	} else {
		host := routing.GetUserHost(route.Name)
		port := routing.GetUserPort(route.Name)

		loginUrl = fmt.Sprintf("http://%s:%s/login", host, port)
	}

	// JSON body
	data := url.Values{}
	data.Set("password", auth.Password)

	// Create a HTTP post request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(data.Encode()))
	if err != nil {
		e.FailOnError(err, "Code-server, login request creation error")
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Host", "localhost")

	resp, err := client.Do(req)

	if err != nil {
		e.FailOnError(err, "Code-server, login response error")
		return
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	if len(cookies) == 0 {
		c.JSON(http.StatusBadRequest, "Login failed, invalid User or Password")
		return
	}

	cookie := cookies[0]

	c.JSON(http.StatusOK, gin.H{
		"path":      fmt.Sprintf("code-editor/%s", route.Path),
		cookie.Name: cookie.Value,
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
