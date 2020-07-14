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
	github.com/micro/cli/v2 v2.1.1
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.0.3-0.20200309150924-5c659fd4b0ad
	github.com/restic/calens v0.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.5.0
	go.opencensus.io v0.22.3
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

