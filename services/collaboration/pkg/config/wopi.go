package config

// Wopi defines the available configuration for the WOPI endpoint.
type Wopi struct {
	WopiSrc string `yaml:"wopisrc" env:"COLLABORATION_WOPI_SRC" desc:"The WOPISrc base URL containing schema, host and port. Set this to the schema and domain where the collaboration service is reachable for the wopi app, such as https://office.owncloud.test." introductionVersion:"5.1"`
	Secret  string `yaml:"secret" env:"COLLABORATION_WOPI_SECRET" desc:"Used to mint and verify WOPI JWT tokens and encrypt and decrypt the REVA JWT token embedded in the WOPI JWT token." introductionVersion:"5.1"`
}
