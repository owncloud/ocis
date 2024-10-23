package config

// Wopi defines the available configuration for the WOPI endpoint.
type Wopi struct {
	WopiSrc     string `yaml:"wopisrc" env:"COLLABORATION_WOPI_SRC" desc:"The WOPI source base URL containing schema, host and port. Set this to the schema and domain where the collaboration service is reachable for the wopi app, such as https://office.owncloud.test." introductionVersion:"6.0.0"`
	Secret      string `yaml:"secret" env:"COLLABORATION_WOPI_SECRET" desc:"Used to mint and verify WOPI JWT tokens and encrypt and decrypt the REVA JWT token embedded in the WOPI JWT token." introductionVersion:"6.0.0"`
	DisableChat bool   `yaml:"disable_chat" env:"COLLABORATION_WOPI_DISABLE_CHAT;OCIS_WOPI_DISABLE_CHAT" desc:"Disable chat in the office web frontend. This feature applies to OnlyOffice and Microsoft." introductionVersion:"7.0.0"`
	ProxyURL    string `yaml:"proxy_url" env:"COLLABORATION_WOPI_PROXY_URL" desc:"The URL to the ownCloud Office365 WOPI proxy. Optional. To use this feature, you need an office365 proxy subscription. If you become part of the Microsoft CSP program (https://learn.microsoft.com/en-us/partner-center/enroll/csp-overview), you can use WebOffice without a proxy." introductionVersion:"7.0.0"`
	ProxySecret string `yaml:"proxy_secret" env:"COLLABORATION_WOPI_PROXY_SECRET" desc:"Optional, the secret to authenticate against the ownCloud Office365 WOPI proxy. This secret can be obtained from ownCloud via the office365 proxy subscription." introductionVersion:"7.0.0"`
}
