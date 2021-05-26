module github.com/owncloud/ocis/ocis

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/asim/go-micro/plugins/logger/zerolog/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/olekukonko/tablewriter v0.0.5
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/ocis/accounts v0.5.3-0.20210216094451-dc73176dc62d
	github.com/owncloud/ocis/glauth v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/graph v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/graph-explorer v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/idp v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/ocs v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/onlyoffice v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/proxy v0.0.0-20210412105747-9b95e9b1191b
	github.com/owncloud/ocis/settings v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/storage v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/store v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/thumbnails v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/web v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/webdav v0.0.0-20210413063522-955bd60edf33
	github.com/rs/zerolog v1.22.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c // indirect
	github.com/thejerf/suture/v4 v4.0.1
	go.opencensus.io v0.23.0
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
)

replace (
	// broken dependency chain for konnect v0.34.0
	github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
	github.com/gomodule/redigo => github.com/gomodule/redigo v1.8.2
	github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1
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
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	// latest version compatible with etcd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
