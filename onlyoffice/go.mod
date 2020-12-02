module github.com/owncloud/ocis/onlyoffice

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.5
	golang.org/x/net v0.0.0-20200822124328-c89045814202
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
