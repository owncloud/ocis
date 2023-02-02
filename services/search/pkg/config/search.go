package config

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint         string `yaml:"endpoint" env:"SEARCH_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster          string `yaml:"cluster" env:"SEARCH_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	AsyncUploads     bool   `yaml:"async_uploads" env:"STORAGE_USERS_OCIS_ASYNC_UPLOADS;SEARCH_EVENTS_ASYNC_UPLOADS" desc:"Enable asynchronous file uploads."`
	NumConsumers     int    `yaml:"num_consumers" env:"SEARCH_EVENTS_NUM_CONSUMERS" desc:"The amount of concurrent event consumers to start. Event consumers are used for searching files. Multiple consumers increase parallelisation, but will also increase CPU and memory demands. The default value is 0."`
	DebounceDuration int    `yaml:"debounce_duration" env:"SEARCH_EVENTS_REINDEX_DEBOUNCE_DURATION" desc:"The duration in milliseconds the reindex debouncer waits before triggering a reindex of a space that was modified."`

	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;SEARCH_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"SEARCH_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided SEARCH_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;SEARCH_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services."`
}
