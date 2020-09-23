module github.com/owncloud/ocis/proxy

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20200810113633-b00aca449666
	github.com/cs3org/reva v1.2.1-0.20200826162318-c0f54e1f37ea
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis/accounts v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/settings v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/store v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.4
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.31.1
)

replace (
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/store => ../store
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
