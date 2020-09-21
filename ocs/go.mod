module github.com/owncloud/ocis/ocs

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cs3org/reva v1.2.1-0.20200826162318-c0f54e1f37ea
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.1.2 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis/accounts v0.0.0-20200918125107-fcca9faa81c8
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/store v0.0.0-20200918125107-fcca9faa81c8
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.4
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	google.golang.org/protobuf v1.25.0
)

replace (
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/store => ../store
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
