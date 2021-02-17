module github.com/owncloud/ocis/graph

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/cs3org/go-cs3apis v0.0.0-20210209082852-35ace33082f5
	github.com/cs3org/reva v1.6.0
	github.com/asim/go-micro/v3 v3.5.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/spf13/afero v1.3.4 // indirect
	github.com/spf13/viper v1.7.1
	github.com/yaegashi/msgraph.go v0.1.4
	go.opencensus.io v0.22.6
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	google.golang.org/grpc v1.35.0
	honnef.co/go/tools v0.1.1 // indirect
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
