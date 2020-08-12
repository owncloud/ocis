module github.com/owncloud/ocis-settings

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/Masterminds/sprig/v3 v3.1.0 // indirect
	github.com/UnnoTed/fileb0x v1.1.4
	github.com/cespare/reflex v0.2.0 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-ozzo/ozzo-validation/v4 v4.2.1
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.4
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mholt/certmagic v0.9.3 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro v1.18.0 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/ogier/pflag v0.0.1 // indirect
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.2.2-0.20200812103920-db41b5a3d14d
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/zenazn/goji v0.9.0 // indirect
	go.opencensus.io v0.22.4
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	golang.org/x/tools v0.0.0-20200811215021-48a8ffc5b207
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200420144010-e5e8543f8aeb
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.23.0
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
