module github.com/owncloud/ocis-settings

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cespare/reflex v0.2.0 // indirect
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-ozzo/ozzo-validation/v4 v4.2.1
	github.com/golang/protobuf v1.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.4
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.6.0
	github.com/mitchellh/gox v1.0.1
	github.com/ogier/pflag v0.0.1 // indirect
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.2.2-0.20200527082518-5641fa4a4c8c
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.0
	go.opencensus.io v0.22.3
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/tools v0.0.0-20200526224456-8b020aee10d2
	google.golang.org/genproto v0.0.0-20200420144010-e5e8543f8aeb
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.21.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
