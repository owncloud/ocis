package config

import (
	"context"
)

// Config combines all available configuration parts.
type Config struct {
	File string
	Log  *Log

	Debug Debug `mask:"struct" yaml:"debug"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`

	InfectedFileHandling string `yaml:"infected-file-handling" env:"ANTIVIRUS_INFECTED_FILE_HANDLING" desc:"Defines the behaviour when a virus has been found. Supported options are: 'delete', 'continue' and 'abort '. Delete will delete the file. Continue will mark the file as infected but continues further processing. Abort will keep the file in the uploads folder for further admin inspection and will not move it to its final destination."`
	Events               Events
	Scanner              Scanner
	MaxScanSize          string `yaml:"max-scan-size" env:"ANTIVIRUS_MAX_SCAN_SIZE" desc:"The maximum scan size the virusscanner can handle. Only this many bytes of a file will be scanned. 0 means unlimited and is the default. Usable common abbreviations: [KB, KiB, GB, GiB, TB, TiB, PB, PiB, EB, EiB], example: 2GB."`

	Context context.Context `yaml:"-" json:"-"`

	DebugScanOutcome string `yaml:"-" env:"ANTIVIRUS_DEBUG_SCAN_OUTCOME" desc:"A predefined outcome for virus scanning, FOR DEBUG PURPOSES ONLY! (example values: 'found,infected')"`
}

// Service defines the available service configuration.
type Service struct {
	Name string `yaml:"-"`
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;ANTIVIRUS_LOG_LEVEL" desc:"The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;ANTIVIRUS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;ANTIVIRUS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;ANTIVIRUS_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"ANTIVIRUS_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"ANTIVIRUS_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"ANTIVIRUS_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"ANTIVIRUS_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;ANTIVIRUS_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;ANTIVIRUS_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;ANTIVIRUS_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;ANTIVIRUS_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided ANTIVIRUS_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;ANTIVIRUS_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services."`
}

// Scanner provides configuration options for the antivirusscanner
type Scanner struct {
	Type string `yaml:"type" env:"ANTIVIRUS_SCANNER_TYPE" desc:"The antivirus scanner to use. Supported values are 'clamav' and 'icap'."`

	ClamAV ClamAV // only if Type == clamav
	ICAP   ICAP   // only if Type == icap
}

// ClamAV provides configuration option for clamav
type ClamAV struct {
	Socket string `yaml:"socket" env:"ANTIVIRUS_CLAMAV_SOCKET" desc:"The socket clamav is running on. Note the default value is an example which needs adaption according your OS."`
}

// ICAP provides configuration option for ICAP
type ICAP struct {
	Timeout int64  `yaml:"timeout" env:"ANTIVIRUS_ICAP_TIMEOUT" desc:"Timeout for the ICAP client."`
	URL     string `yaml:"url" env:"ANTIVIRUS_ICAP_URL" desc:"URL of the ICAP server."`
	Service string `yaml:"service" env:"ANTIVIRUS_ICAP_SERVICE" desc:"The name of the ICAP service."`
}
