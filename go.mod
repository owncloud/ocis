module github.com/owncloud/ocis-settings

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.4
	github.com/micro/cli/v2 v2.1.1
	github.com/micro/go-micro/v2 v2.0.0
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-hello v0.1.0-alpha1
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.1
	go.opencensus.io v0.22.2
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	google.golang.org/genproto v0.0.0-20200420144010-e5e8543f8aeb
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
