package models

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type User struct {
	Id       string `yaml:"id"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	IsAdmin  bool   `yaml:"isAdmin"`
}

type Users struct {
	Users []User
}

type View struct {
	Id             string `yaml:"id"`
	UserId         string `yaml:"userId"`
	Status         string `yaml:"status"`
	Path           string `yaml:"path"`
	Password       string `yaml:"password"`
	VScodeSettings string `yaml:"vscode-settings"`
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

type Extension struct {
	Id       string                 `yaml:"id"`
	Settings map[string]interface{} `yaml:"settings"`
}

type EnableConfig struct {
	Git            GitConfig              `json:"git"`
	Extensions     []Extension            `json:"extensions"`
	VscodeSettings map[string]interface{} `json:"vscode-settings"`
}
