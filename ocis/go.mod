module github.com/owncloud/ocis/ocis

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/go-test/deep v1.0.6 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.8.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis-graph v0.0.0-20200318175820-9a5a6e029db7
	github.com/owncloud/ocis-graph-explorer v0.0.0-20200210111049-017eeb40dc0c
	github.com/owncloud/ocis-hello v0.1.0-alpha1.0.20200828085053-37fcf3c8f853
	github.com/owncloud/ocis/accounts v0.5.3-0.20201103104733-ff2c41028d9b
	github.com/owncloud/ocis/glauth v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/konnectd v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/ocis-phoenix v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/ocis-pkg v0.1.0
	github.com/owncloud/ocis/ocs v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/onlyoffice v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/proxy v0.0.0-00010101000000-000000000000
	github.com/owncloud/ocis/settings v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/storage v0.0.0-20201015120921-38358ba4d4df
	github.com/owncloud/ocis/store v0.0.0-20200918125107-fcca9faa81c8
	github.com/owncloud/ocis/thumbnails v0.1.6
	github.com/owncloud/ocis/webdav v0.0.0-00010101000000-000000000000
	github.com/refs/pman v0.0.0-20200701173654-f05b8833071a
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.1
	go.opencensus.io v0.22.5
)

replace (
	github.com/cs3org/reva => github.com/labkode/reva v0.0.0-20201202134237-befa4a5708b6
	github.com/gomodule/redigo => github.com/gomodule/redigo v1.8.2
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/glauth => ../glauth
	github.com/owncloud/ocis/konnectd => ../konnectd
	github.com/owncloud/ocis/ocis-phoenix => ../ocis-phoenix
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/ocs => ../ocs
	github.com/owncloud/ocis/onlyoffice => ../onlyoffice
	github.com/owncloud/ocis/proxy => ../proxy
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/storage => ../storage
	github.com/owncloud/ocis/store => ../store
	github.com/owncloud/ocis/thumbnails => ../thumbnails
	github.com/owncloud/ocis/webdav => ../webdav
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
