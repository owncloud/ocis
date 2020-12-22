module github.com/owncloud/ocis/webdav

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/thumbnails v0.0.0-20200918125107-fcca9faa81c8
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.5
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/thumbnails => ../thumbnails
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

replace github.com/cs3org/reva => github.com/butonic/reva v0.0.0-20201222133934-8e0b5aaf8aeb
