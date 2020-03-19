module github.com/owncloud/ocis-webdav

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/micro/cli/v2 v2.1.1
	github.com/micro/go-micro/v2 v2.0.0
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/owncloud/ocis-thumbnails v0.0.0-20200318131505-e0ab0b37a5a4
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.5.0
	go.opencensus.io v0.22.2
)

replace github.com/owncloud/ocis-pkg/v2 => /home/corby/Development/go/src/github.com/owncloud/ocis-pkg