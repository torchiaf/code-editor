package config

import (
	"os"
	"server/utils"

	models "server/models"
)

type config struct {
	IsDev     bool
	App       string
	Namespace string
	Users     map[string]models.User
}

func isDevEnv() bool {
	env := os.Getenv("DEV_ENV")
	if len(env) > 0 {
		return true
	}
	return false
}

func getConfig() config {

	arr := utils.ParseFile[models.Users]("assets/users").Users
	userMap := utils.Map(arr, func(user models.User) string { return user.Name })

	c := config{
		IsDev:     isDevEnv(),
		App:       "code-editor",
		Namespace: utils.IfNull(os.Getenv("POD_NAMESPACE"), "code-editor"),
		Users:     userMap,
	}

	return c
}

var Config = getConfig()
