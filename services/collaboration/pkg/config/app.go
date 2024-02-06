package config

// App defines the available app configuration.
type App struct {
	Name        string `yaml:"name" env:"COLLABORATION_APP_NAME" desc:"The name of the app"`
	Description string `yaml:"description" env:"COLLABORATION_APP_DESCRIPTION" desc:"App description"`
	Icon        string `yaml:"icon" env:"COLLABORATION_APP_ICON" desc:"Icon for the app"`
	LockName    string `yaml:"lockname" env:"COLLABORATION_APP_LOCKNAME" desc:"Name for the app lock"`
}
