module github.com/owncloud/ocis-proxy

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20200730121022-c4f3d4f7ddfd
	github.com/cs3org/reva v1.1.0
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis-accounts v0.1.2-0.20200618163128-aa8ae58dd95e
	github.com/owncloud/ocis-pkg/v2 v2.4.0
	github.com/owncloud/ocis-settings v0.3.2-0.20200827134219-0f9e141699e5
	github.com/owncloud/ocis-store v0.0.0-20200716140351-f9670592fb7b
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.4
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.31.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
