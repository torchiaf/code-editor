package authentication

import (
	"errors"
	"server/config"
	"server/models"
	"server/users"

	"golang.org/x/crypto/bcrypt"
)

func verifyPassword(password, userPassword string) error {

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(auth models.Auth) (string, error) {

	if config.Config.Authentication.IsExternal {
		return "", errors.New("External authentication enabled, must register user first.")
	}

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
