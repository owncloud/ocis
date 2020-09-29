module github.com/owncloud/ocis/settings

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-ozzo/ozzo-validation/v4 v4.2.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/mitchellh/gox v1.0.1
	github.com/oklog/run v1.0.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.4
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884
	google.golang.org/protobuf v1.23.0
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
