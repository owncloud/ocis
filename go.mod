module github.com/owncloud/ocis-proxy

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20200306065539-29abc33f5be0
	github.com/cs3org/reva v0.1.0
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.6.0
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis-accounts v0.1.2-0.20200522102615-8c7da929195a
	github.com/owncloud/ocis-pkg/v2 v2.2.1
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	go.opencensus.io v0.22.3
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	gopkg.in/square/go-jose.v2 v2.4.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
