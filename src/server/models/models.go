package models

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type User struct {
	Id       string `yaml:"id"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Users struct {
	Users []User
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

type ExternalUserLogin struct {
	Username string
	Password string
	Token    string
	Force    bool
}
