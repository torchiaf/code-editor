package models

type Route struct {
	Name string
	Path string
}

type Routes struct {
	Routes []Route
}

type Auth struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
