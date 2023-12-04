package models

type User struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Path     string `yaml:"path"`
}

type Users struct {
	Users []User
}

type CodeServerSession struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
