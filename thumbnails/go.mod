module github.com/owncloud/ocis/thumbnails

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/Masterminds/sprig/v3 v3.1.0 // indirect
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cespare/reflex v0.2.0
	github.com/go-test/deep v1.0.2-0.20181118220953-042da051cf31 // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang/protobuf v1.4.2
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro v1.18.0 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/ogier/pflag v0.0.1
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.4
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	gopkg.in/square/go-jose.v2 v2.5.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
