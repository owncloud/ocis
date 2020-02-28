module github.com/owncloud/ocis

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/micro/cli/v2 v2.1.2-0.20200203150404-894195727d9c
	github.com/micro/go-micro/v2 v2.0.1-0.20200207205803-ef537270add3
	github.com/micro/micro/v2 v2.0.1-0.20200210100719-f38a1d8d5348
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-devldap v0.0.0-20200210154914-6c796445c1bf
	github.com/owncloud/ocis-graph v0.0.0-20200217115956-172417259283
	github.com/owncloud/ocis-graph-explorer v0.0.0-20200210111049-017eeb40dc0c
	github.com/owncloud/ocis-hello v0.1.0-alpha1.0.20200207094758-c866cafca7e5
	github.com/owncloud/ocis-konnectd v0.0.0-20200218084322-104a3fcf99c7
	github.com/owncloud/ocis-ocs v0.0.0-20200207130609-800a64d45fac
	github.com/owncloud/ocis-phoenix v0.1.1-0.20200213204418-06f50c42c225
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/owncloud/ocis-reva v0.0.0-20200213202552-584d47daa8bc
	github.com/owncloud/ocis-webdav v0.0.0-20200210113150-6c4d498c38b0
	github.com/spf13/viper v1.6.1
	go.opencensus.io v0.22.2
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	stash.kopano.io/kc/konnect v0.29.0 // indirect
)

replace stash.kopano.io/kc/konnect => github.com/IljaN/konnect v0.30.0-alpha1
