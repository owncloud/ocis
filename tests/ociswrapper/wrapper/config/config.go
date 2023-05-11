package config

var config = map[string]string{
	"port": "5000",
}

func Get(key string) string {
	return config[key]
}

func Set(key string, value string) {
	config[key] = value
}
