package authentication

import (
	"errors"
	"fmt"
	"os"
	"server/config"
	"server/models"
	utils "server/utils"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var tokenSecret = utils.IfNull(os.Getenv("API_SECRET"), "francesco")
var tokenExpiration = utils.IfNull(os.Getenv("API_TOKEN_EXPIRATION"), "24")

func GenerateToken(username string) (string, error) {

	token_lifespan, err := strconv.Atoi(tokenExpiration)

	if err != nil {
		return "", err
	}

	found, user := utils.Find(config.Config.Users, "Name", username)

	if !found {
		return "", errors.New("User not found")
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["path"] = user.Path
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))

}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Request.Header.Get("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractUser(c *gin.Context) (models.User, error) {
	var user models.User
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return user, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user.Name = claims["username"].(string)
		user.Path = claims["path"].(string)
	}
	return user, nil
}
