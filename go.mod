module github.com/owncloud/ocis-ocs

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cs3org/reva v0.1.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.6.0
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.2.2-0.20200527082518-5641fa4a4c8c
	github.com/owncloud/ocis-store v0.0.0-20200716140351-f9670592fb7b
	github.com/restic/calens v0.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.6.3
	go.opencensus.io v0.22.3
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
