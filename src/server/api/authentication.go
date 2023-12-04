package api

import (
	utils "server/utils"

	"golang.org/x/crypto/bcrypt"

	cfg "server/config"
)

var config = cfg.Config

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {

	found, user := utils.Find(config.Users, "Name", username)

	if !found {
		return err ...continue
	}

	err := VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}
