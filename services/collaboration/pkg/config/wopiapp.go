package config

// WopiApp defines the available configuration in order to connect to a WOPI app.
type WopiApp struct {
	Addr     string `yaml:"addr" env:"COLLABORATION_WOPIAPP_ADDR" desc:"The URL where the WOPI app is located, such as https://127.0.0.1:8080." introductionVersion:"5.1"`
	Insecure bool   `yaml:"insecure" env:"COLLABORATION_WOPIAPP_INSECURE" desc:"Skip TLS certificate verification when connecting to the WOPI app" introductionVersion:"5.1"`
	WopiSrc  string `yaml:"wopisrc" env:"COLLABORATION_WOPIAPP_WOPISRC" desc:"The WOPISrc base URL containing schema, host and port. Path will be set to /wopi/files/{fileid} if not given. {fileid} will be replaced with the WOPI file id. Set this to the schema and domain where the collaboration service is reachable for the wopi app, such as https://office.owncloud.test." introductionVersion:"5.1"`
	Secret   string `yaml:"secret" env:"COLLABORATION_WOPIAPP_SECRET" desc:"Used to mint and verify WOPI JWT tokens and encrypt and decrypt the REVA JWT token embedded in the WOPI JWT token." introductionVersion:"5.1"`
}
