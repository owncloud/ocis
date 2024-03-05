package config

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint         string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;SEARCH_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture." introductionVersion:"pre5.0"`
	Cluster          string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;SEARCH_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system." introductionVersion:"pre5.0"`
	AsyncUploads     bool   `yaml:"async_uploads" env:"OCIS_ASYNC_UPLOADS;SEARCH_EVENTS_ASYNC_UPLOADS" desc:"Enable asynchronous file uploads." introductionVersion:"pre5.0"`
	NumConsumers     int    `yaml:"num_consumers" env:"SEARCH_EVENTS_NUM_CONSUMERS" desc:"The amount of concurrent event consumers to start. Event consumers are used for searching files. Multiple consumers increase parallelisation, but will also increase CPU and memory demands. The default value is 0." introductionVersion:"pre5.0"`
	DebounceDuration int    `yaml:"debounce_duration" env:"SEARCH_EVENTS_REINDEX_DEBOUNCE_DURATION" desc:"The duration in milliseconds the reindex debouncer waits before triggering a reindex of a space that was modified." introductionVersion:"pre5.0"`

	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;SEARCH_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates." introductionVersion:"pre5.0"`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;SEARCH_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided SEARCH_EVENTS_TLS_INSECURE will be seen as false." introductionVersion:"pre5.0"`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;SEARCH_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"pre5.0"`
	AuthUsername         string `yaml:"username" env:"OCIS_EVENTS_AUTH_USERNAME;SEARCH_EVENTS_AUTH_USERNAME" desc:"The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"5.0"`
	AuthPassword         string `yaml:"password" env:"OCIS_EVENTS_AUTH_PASSWORD;SEARCH_EVENTS_AUTH_PASSWORD" desc:"The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"5.0"`
}
