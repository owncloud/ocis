package config

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint     string `yaml:"endpoint" env:"SEARCH_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster      string `yaml:"cluster" env:"SEARCH_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	AsyncUploads bool   `yaml:"async_uploads" env:"STORAGE_USERS_OCIS_ASYNC_UPLOADS;SEARCH_EVENTS_ASYNC_UPLOADS" desc:"Enable asynchronous file uploads."`
}
