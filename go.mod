module github.com/owncloud/ocis-graph

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/cs3org/go-cs3apis v0.0.0-20191128165347-19746c015c83
	github.com/cs3org/reva v0.0.2-0.20200115110931-4c7513415ec5
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/micro/cli/v2 v2.1.1
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.5.0
	github.com/yaegashi/msgraph.go v0.0.0-20191104022859-3f9096c750b2
	go.opencensus.io v0.22.2
	google.golang.org/grpc v1.26.0
	gopkg.in/ldap.v3 v3.1.0
)
