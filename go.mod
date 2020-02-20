module github.com/owncloud/ocis-konnectd

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gorilla/mux v1.7.3
	github.com/micro/cli/v2 v2.1.1
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg v1.3.0
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/rs/zerolog v1.17.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
	go.opencensus.io v0.22.2
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	stash.kopano.io/kc/konnect v0.28.1
	stash.kopano.io/kgol/rndm v1.1.0
)

replace stash.kopano.io/kc/konnect => github.com/IljaN/konnect v0.29.0-alpha2
