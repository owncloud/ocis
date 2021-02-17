module github.com/owncloud/ocis/idp

go 1.15

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210210071449-f9f5e7422d8f // indirect
	github.com/asim/go-micro/v3 v3.5.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/mitchellh/copystructure v1.1.1 // indirect
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/pquerna/cachecontrol v0.0.0-20200921180117-858c6e7e6b7e // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.6
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	stash.kopano.io/kc/konnect v0.33.3
	stash.kopano.io/kgol/rndm v1.1.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	// latest version compatible with etcd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
