package config

import "context"

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// Server defines the available server configuration.
type Server struct {
	Name      string
	Namespace string
	Address   string
	Version   string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Config combines all available configuration parts.
type Config struct {
	File      string
	Log       Log
	Debug     Debug
	Server    Server
	Tracing   Tracing
	Thumbnail Thumbnail

	Context    context.Context
	Supervised bool
}

// FileSystemStorage defines the available filesystem storage configuration.
type FileSystemStorage struct {
	RootDirectory string
}

// FileSystemSource defines the available filesystem source configuration.
type FileSystemSource struct {
	BasePath string
}

// Thumbnail defines the available thumbnail related configuration.
type Thumbnail struct {
	Resolutions         []string
	FileSystemStorage   FileSystemStorage
	WebdavAllowInsecure bool
	CS3AllowInsecure    bool
	RevaGateway         string
	WebdavNamespace     string
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
