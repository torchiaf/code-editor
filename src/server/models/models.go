package models

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type User struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Path     string `yaml:"path"`
}

type Users struct {
	Users []User
}

type Service struct {
	Id   string `yaml:"id"`
	Port int32  `yaml:"port"`
}

type CodeServerConfig struct {
	ServicePort int32  `yaml:"port"`
	Password    string `yaml:"password"`
}

type CodeServerSession struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Git struct {
	Type   string `yaml:"type"`
	Org    string `yaml:"org"`
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch"`
	Commit string `yaml:"commit"`
}

type ViewConfig struct {
	Git Git
}
