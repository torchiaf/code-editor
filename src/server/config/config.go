package config

import (
	"os"
	"server/utils"

	models "server/models"
)

type app struct {
	Name      string
	Namespace string
}

type config struct {
	IsDev bool
	App   app
	Users map[string]models.User
}

func isDevEnv() bool {
	env := os.Getenv("DEV_ENV")
	if len(env) > 0 {
		return true
	}
	return false
}

func getUsers() map[string]models.User {
	users := utils.ParseFile[models.Users]("assets/users/users.yaml").Users

	return utils.Map(users, func(user models.User) string { return user.Name })
}

func initConfig() config {

	c := config{
		IsDev: isDevEnv(),
		Users: getUsers(),
		App: app{
			Name:      utils.IfNull(os.Getenv("APP_NAME"), "code-editor"),
			Namespace: utils.IfNull(os.Getenv("APP_NAMESPACE"), "code-editor"),
		},
	}

	return c
}

var Config = initConfig()
