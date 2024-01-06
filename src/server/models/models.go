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

type ExternalUserLogin struct {
	Username string
	Password string
	Token    string
	Force    bool
}

type CodeServerSession struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type GitSource struct {
	Type   string `yaml:"type"`
	Org    string `yaml:"org"`
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch"`
	Commit string `yaml:"commit"`
}

type ViewConfig struct {
	Git GitSource
}

type GitConfig struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type EnableConfig struct {
	Git        GitConfig
	Extensions []string
}
