package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	config "server/config"
	kube "server/kubernetes"

	"github.com/gin-gonic/gin"
)

var localConfig = config.GetConfig()

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "paolo",
	})
}

func Auth(c *gin.Context) {

	// HTTP endpoint
	loginUrl := "http://localhost/code-editor/login" //, localConfig.CodeServerHost, localConfig.CodeServerPort)

	// JSON body
	data := url.Values{}
	data.Set("password", "1bf6c9ffb4ff401d36bc9bf0")

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

	cookie := resp.Cookies()[0]

	c.JSON(http.StatusOK, gin.H{
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
