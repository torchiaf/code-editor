package authentication

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"server/config"
	"server/models"
	"server/users"
	utils "server/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func isExternal(user models.User) bool {
	return strings.HasPrefix(user.Id, "ext-")
}

func verifyPassword(password, userPassword string) error {

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(auth models.Auth) (string, error) {

	user, ok := users.Store.Get(auth.Username)

	if !ok {
		return "", errors.New("User not found")
	}

	err := verifyPassword(auth.Password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", errors.New("Password is not correct")
	}

	token, err := GenerateToken(user.Name)

	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyExternalUser(token string) (string, error) {

	// Disable tls check
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Create a HTTP post request
	client := &http.Client{}

	req, err := http.NewRequest("GET", config.Config.Authentication.Url, &strings.Reader{})
	if err != nil {
		return "", errors.New("External Login, request creation error")
	}

	if config.Config.Authentication.TokenType == config.TOKEN_TYPE_HEADERS {
		req.Header.Add(config.Config.Authentication.TokenKey, token)
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", errors.New("External Login, login response error")
	}
	defer resp.Body.Close()

	var v map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		panic(err)
	}

	// TODO expand error type based on response error: missing Url; missing token; incorrect token; etc.
	if resp.StatusCode != 200 {
		return "", errors.New("External Login check failed. Wrong Url or Token")
	}

	name, err := utils.JsonQuery[string](v, config.Config.Authentication.Query)
	if err != nil {
		return "", errors.New("External Login check failed. Cannot get username")
	}

	return name, nil
}
