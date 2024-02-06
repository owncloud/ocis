package config

// WopiApp defines the available configuration in order to connect to a WOPI app.
type WopiApp struct {
	Addr     string `yaml:"addr" env:"COLLABORATION_WOPIAPP_ADDR" desc:"The URL where the WOPI app is located, such as https://127.0.0.1:8080."`
	Insecure bool   `yaml:"insecure" env:"COLLABORATION_WOPIAPP_INSECURE" desc:"Connect insecurely"`
}
