package config

import "os"

type Config struct {
	IsDev          bool
	App            string
	Namespace      string
	CodeServerHost string
	CodeServerPort string
}

func isDevEnv() bool {
	env := os.Getenv("DEV_ENV")
	if len(env) > 0 {
		return true
	}
	return false
}

func GetConfig() Config {

	c := Config{
		IsDev:          isDevEnv(),
		App:            "code-editor",
		Namespace:      os.Getenv("POD_NAMESPACE"),
		CodeServerHost: os.Getenv("CODE_EDITOR_SERVICE_HOST"),
		CodeServerPort: os.Getenv("CODE_EDITOR_SERVICE_PORT=8080"),
	}

	return c
}
