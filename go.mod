module github.com/owncloud/ocis-proxy

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.5.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/micro/cli/v2 v2.1.2-0.20200203150404-894195727d9c
	github.com/micro/go-micro/v2 v2.0.1-0.20200212105717-d76baf59de2e
	github.com/micro/go-plugins/micro/router/v2 v2.0.3-0.20200221093116-8ed9b03043f0
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.1.6
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/prometheus/client_golang v1.2.1
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.1
	go.opencensus.io v0.22.2
)
