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

	"github.com/gin-gonic/gin"
)

var localConfig = config.GetConfig()

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func Auth(c *gin.Context) {

	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	auth := model.Auth{}

	err1 := json.Unmarshal([]byte(reqBody), &auth)
	if err1 != nil {
		e.FailOnError(err, "Failed to get Deployment/Scale")
	}

	path := ""
	for _, route := range localConfig.Routes {
		if auth.User == route.Name {
			path = route.Path
		}
	}

	if path == "" {
		// TODO return , user not found
	}

	// HTTP endpoint
	loginUrl := fmt.Sprintf("http://localhost/code-editor/%s/login", path)

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
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Host", "localhost")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	if len(cookies) == 0 {
		c.JSON(http.StatusOK, "Invalid User or Password")
	}

	cookie := cookies[0]

	c.JSON(http.StatusOK, gin.H{
		"path":      fmt.Sprintf("code-editor/%s", path),
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

func GetPods(c *gin.Context) {
	pods := kube.GetPods()
	c.JSON(200, pods)
}
