module github.com/owncloud/ocis-graph-explorer

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/micro/cli/v2 v2.1.1
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.6.1
	go.opencensus.io v0.22.2
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
)
