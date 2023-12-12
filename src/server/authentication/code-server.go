package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"server/config"
	"server/kubernetes"
	"server/models"
)

func EditorLogin(user models.User) (models.CodeServerSession, error) {

	var session models.CodeServerSession

	// code-server login endpoint
	loginUrl := ""

	if config.Config.IsDev {
		loginUrl = fmt.Sprintf("http://localhost/code-editor/%s/login", user.Path)
	} else {
		host := kubernetes.GetUserHost(user.Name)
		port := kubernetes.GetUserPort(user.Name)

		loginUrl = fmt.Sprintf("http://%s:%s/login", host, port)
	}

	// JSON body
	data := url.Values{}
	data.Set("password", user.Password)

	// Create a HTTP post request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return session, errors.New("Code-server, login request creation error")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Host", "localhost")

	resp, err := client.Do(req)

	if err != nil {
		return session, errors.New("Code-server, login response error")
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	if len(cookies) == 0 {
		return session, errors.New("Login failed, invalid User or Password")
	}

	cookie := cookies[0]

	session.Name = cookie.Name
	session.Value = cookie.Value

	return session, nil
}
