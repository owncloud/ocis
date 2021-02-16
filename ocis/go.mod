module github.com/owncloud/ocis/ocis

go 1.15

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/asim/go-micro/plugins/registry/mdns/v3 v3.0.0-20210210071449-f9f5e7422d8f // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/go-test/deep v1.0.6 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.8.0
	github.com/mitchellh/copystructure v1.1.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/nsf/termbox-go v0.0.0-20210114135735-d04385b850e8 // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis-hello v0.1.0-alpha1.0.20210204050952-c291e4c5b73f
	github.com/owncloud/ocis/accounts v0.5.3-0.20210216094451-dc73176dc62d
	github.com/owncloud/ocis/glauth v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/graph v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/graph-explorer v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/idp v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/ocs v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/onlyoffice v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/proxy v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/settings v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/storage v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/store v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/thumbnails v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/web v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/webdav v0.0.0-20210216094451-dc73176dc62d
	github.com/restic/calens v0.2.0
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.opencensus.io v0.22.6
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	honnef.co/go/tools v0.1.1 // indirect
)

replace (
	github.com/gomodule/redigo => github.com/gomodule/redigo v1.8.2
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/glauth => ../glauth
	github.com/owncloud/ocis/graph => ../graph
	github.com/owncloud/ocis/graph-explorer => ../graph-explorer
	github.com/owncloud/ocis/idp => ../idp
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/ocs => ../ocs
	github.com/owncloud/ocis/onlyoffice => ../onlyoffice
	github.com/owncloud/ocis/proxy => ../proxy
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/storage => ../storage
	github.com/owncloud/ocis/store => ../store
	github.com/owncloud/ocis/thumbnails => ../thumbnails
	github.com/owncloud/ocis/web => ../web
	github.com/owncloud/ocis/webdav => ../webdav
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
