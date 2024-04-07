package config

import (
	"fmt"
	"os"

	"github.com/torchiaf/code-editor/server/utils"
)

const LOCAL string = "local"
const EXTERNAL string = "external"
const TOKEN_TYPE_HEADERS string = "headers"

type authentication struct {
	IsExternal bool
	Url        string
	TokenType  string
	TokenKey   string
	Query      string
}

type app struct {
	Name      string
	Namespace string
}

type resources struct {
	IngressName string
	ConfigName  string
	UsersName   string
}

type config struct {
	IsDev          bool
	Authentication authentication
	App            app
	Resources      resources
}

func isDevEnv() bool {
	env := os.Getenv("DEV_ENV")
	if len(env) > 0 {
		return true
	}
	return false
}

func initConfig() config {

	app := app{
		Name:      utils.IfNull(os.Getenv("APP_NAME"), "code-editor"),
		Namespace: utils.IfNull(os.Getenv("APP_NAMESPACE"), "code-editor"),
	}

	c := config{
		IsDev: isDevEnv(),
		App:   app,
		Authentication: authentication{
			IsExternal: os.Getenv("AUTH_TYPE") == EXTERNAL,
			Url:        os.Getenv("AUTH_URL"),
			TokenType:  os.Getenv("AUTH_TOKEN_TYPE"),
			TokenKey:   os.Getenv("AUTH_TOKEN_KEY"),
			Query:      os.Getenv("AUTH_QUERY"),
		},
		Resources: resources{
			IngressName: fmt.Sprintf("%s-ui", app.Name),
			ConfigName:  fmt.Sprintf("%s-config", app.Name),
			UsersName:   fmt.Sprintf("%s-users", app.Name),
		},
	}

	return c
}

var Config = initConfig()
