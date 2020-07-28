module github.com/owncloud/ocis-ocs

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cs3org/go-cs3apis v0.0.0-20200611124600-7a1be2026543 // indirect
	github.com/cs3org/reva v0.1.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.6.0
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-accounts v0.1.2-0.20200727195215-6816703df41d
	github.com/owncloud/ocis-pkg/v2 v2.2.2-0.20200602070144-cd0620668170
	github.com/owncloud/ocis-store v0.0.0-20200716140351-f9670592fb7b
	github.com/prometheus/client_golang v1.7.0 // indirect
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.19.0 // indirect
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.4
	google.golang.org/protobuf v1.25.0
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
