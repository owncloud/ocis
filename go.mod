module github.com/owncloud/ocis

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/coreos/etcd v3.3.21+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-log/log v0.2.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/lucas-clemente/quic-go v0.15.7 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.7.0
	github.com/micro/micro v1.16.0
	github.com/micro/micro/v2 v2.7.0
	github.com/miekg/dns v1.1.29 // indirect
	github.com/nats-io/nats.go v1.10.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/flaex v0.2.0
	github.com/owncloud/ocis-accounts v0.1.1
	github.com/owncloud/ocis-glauth v0.4.0
	github.com/owncloud/ocis-graph v0.0.0-20200318175820-9a5a6e029db7
	github.com/owncloud/ocis-graph-explorer v0.0.0-20200210111049-017eeb40dc0c
	github.com/owncloud/ocis-hello v0.1.0-alpha1.0.20200207094758-c866cafca7e5
	github.com/owncloud/ocis-konnectd v0.3.1
	github.com/owncloud/ocis-migration v0.0.0-20200504185909-72274a4f1449
	github.com/owncloud/ocis-ocs v0.0.0-20200318181133-cc66a0531da7
	github.com/owncloud/ocis-phoenix v0.6.0
	github.com/owncloud/ocis-pkg/v2 v2.2.1
	github.com/owncloud/ocis-proxy v0.3.1
	github.com/owncloud/ocis-reva v0.2.2-0.20200513073117-ee9cd9b8d3ab
	github.com/owncloud/ocis-thumbnails v0.1.2-0.20200422124828-f92a40879feb
	github.com/owncloud/ocis-webdav v0.1.0
	github.com/refs/pman v0.0.0-20200520152433-d1823a649d98
	github.com/restic/calens v0.2.0
	go.opencensus.io v0.22.3
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/net v0.0.0-20200519113804-d87ec0cfa476 // indirect
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299 // indirect
	golang.org/x/tools v0.0.0-20200519205726-57a9e4404bf7 // indirect
	google.golang.org/genproto v0.0.0-20200519141106-08726f379972 // indirect
	gopkg.in/olivere/elastic.v5 v5.0.83 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/lucas-clemente/quic-go v0.15.7 => github.com/lucas-clemente/quic-go v0.14.1
