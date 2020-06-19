module github.com/owncloud/ocis-glauth

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/GeertJohan/yubigo v0.0.0-20190917122436-175bc097e60e
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/glauth/glauth v1.1.3-0.20200228160118-2d4f5d547682
	github.com/go-logr/logr v0.1.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.6.0
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484
	github.com/nmcclain/ldap v0.0.0-20191021200707-3b3b69a7e9e3
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-accounts v0.1.2-0.20200617152311-02e759f95e82
	github.com/owncloud/ocis-pkg/v2 v2.2.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.19.0
	github.com/spf13/viper v1.7.0
	go.opencensus.io v0.22.4
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
