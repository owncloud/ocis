module github.com/owncloud/ocis/glauth

go 1.15

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/GeertJohan/yubigo v0.0.0-20190917122436-175bc097e60e
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/asim/go-micro/v3 v3.5.0
	github.com/glauth/glauth v1.1.3-0.20201110124627-fd3ac7e4bbdc
	github.com/go-logr/logr v0.1.0
	github.com/micro/cli/v2 v2.1.2
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484
	github.com/nmcclain/ldap v0.0.0-20191021200707-3b3b69a7e9e3
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/ocis/accounts v0.5.3-0.20210216094451-dc73176dc62d
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/spf13/viper v1.7.1
	go.opencensus.io v0.22.6
)

replace (
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
