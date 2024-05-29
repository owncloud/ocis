package config

// WopiApp defines the available configuration in order to connect to a WOPI app.
type WopiApp struct {
	Addr     string `yaml:"addr" env:"COLLABORATION_WOPIAPP_ADDR" desc:"The URL where the WOPI app is located, such as https://127.0.0.1:8080." introductionVersion:"5.1"`
	Insecure bool   `yaml:"insecure" env:"COLLABORATION_WOPIAPP_INSECURE" desc:"Skip TLS certificate verification when connecting to the WOPI app" introductionVersion:"5.1"`
	WopiSrc  string `yaml:"wopisrc" env:"OCIS_URL;COLLABORATION_WOPIAPP_WOPISRC" desc:"The WOPISrc base URL containing schema, host and port. Path will be set to /wopi/files/{fileid} if not given. {fileid} will be replaced with the WOPI file id. Set this to the schema and domain where the collaboration service is reachable for the wopi app, such as https://cloud.owncloud.test or, if you need to set up an additional WOPI server https://wopi-other.owncloud.test." introductionVersion:"5.1"`
}
