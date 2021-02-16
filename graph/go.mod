module github.com/owncloud/ocis/graph

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/cs3org/go-cs3apis v0.0.0-20200810113633-b00aca449666
	github.com/cs3org/reva v1.1.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-ldap/ldap/v3 v3.2.3
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.2.1
	github.com/spf13/afero v1.3.4 // indirect
	github.com/spf13/viper v1.7.1
	github.com/yaegashi/msgraph.go v0.1.4
	go.opencensus.io v0.22.4
	google.golang.org/grpc v1.31.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
