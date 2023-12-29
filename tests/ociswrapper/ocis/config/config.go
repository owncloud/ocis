package config

var config = map[string]string{
	"bin":   "/usr/bin/ocis",
	"url":   "https://localhost:9200",
	"retry": "5",
}

func Set(key string, value string) {
	config[key] = value
}

func Get(key string) string {
	return config[key]
}
