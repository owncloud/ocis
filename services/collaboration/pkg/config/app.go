package config

// App defines the available app configuration.
type App struct {
	Name        string `yaml:"name" env:"COLLABORATION_APP_NAME" desc:"The name of the app, either Collabora, OnlyOffice, Microsoft365 or MicrosoftOfficeOnline" introductionVersion:"6.0.0"`
	Description string `yaml:"description" env:"COLLABORATION_APP_DESCRIPTION" desc:"App description" introductionVersion:"6.0.0"`
	Icon        string `yaml:"icon" env:"COLLABORATION_APP_ICON" desc:"Icon for the app" introductionVersion:"6.0.0"`
	LockName    string `yaml:"lockname" env:"COLLABORATION_APP_LOCKNAME" desc:"Name for the app lock" introductionVersion:"6.0.0"`

	Addr     string `yaml:"addr" env:"COLLABORATION_APP_ADDR" desc:"The URL where the WOPI app is located, such as https://127.0.0.1:8080." introductionVersion:"6.0.0"`
	Insecure bool   `yaml:"insecure" env:"COLLABORATION_APP_INSECURE" desc:"Skip TLS certificate verification when connecting to the WOPI app" introductionVersion:"6.0.0"`

	ProofKeys          ProofKeys `yaml:"proofkeys"`
	LicenseCheckEnable bool      `yaml:"licensecheckenable" env:"COLLABORATION_APP_LICENSE_CHECK_ENABLE" desc:"Enable license check for edit. Needs to be enabled when using Microsoft365 with the business flow." introductionVersion:"%%NEXT%%"`
}

type ProofKeys struct {
	Disable  bool   `yaml:"disable" env:"COLLABORATION_APP_PROOF_DISABLE" desc:"Disable the proof keys verification" introductionVersion:"6.0.0"`
	Duration string `yaml:"duration" env:"COLLABORATION_APP_PROOF_DURATION" desc:"Duration for the proof keys to be cached in memory, using time.ParseDuration format. If the duration can't be parsed, we'll use the default 12h as duration" introductionVersion:"6.0.0"`
}
