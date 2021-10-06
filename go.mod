module github.com/owncloud/ocis

go 1.16

require (
	github.com/CiscoM31/godata v1.0.4
	github.com/GeertJohan/yubigo v0.0.0-20190917122436-175bc097e60e
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/logger/zerolog/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/kubernetes/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/mdns/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/registry/nats/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/server/http/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/plugins/wrapper/trace/opencensus/v3 v3.0.0-20210812172626-c7195aae9817
	github.com/asim/go-micro/v3 v3.6.1-0.20210924081004-8c39b1e1204d
	github.com/blevesearch/bleve/v2 v2.1.0
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/cs3org/go-cs3apis v0.0.0-20210922150613-cb9e3c99f8de
	github.com/cs3org/reva v1.13.1-0.20211006080436-67f39be571fa
	github.com/disintegration/imaging v1.6.2
	github.com/glauth/glauth v1.1.3-0.20210729125545-b9aecdfcac31
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/render v1.0.1
	github.com/go-logr/logr v0.4.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.5.2
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/iancoleman/strcase v0.2.0
	github.com/justinas/alice v1.2.0
	github.com/libregraph/lico v0.34.1-0.20210803054646-b584e0372224
	github.com/mennanov/fieldmask-utils v0.4.0
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484
	github.com/nmcclain/ldap v0.0.0-20210720162743-7f8d1e44eeba
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/owncloud/open-graph-api-go v0.0.0-20210511151655-57894f7d46fb
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.25.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/thejerf/suture/v4 v4.0.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/yaegashi/msgraph.go v0.1.4
	go.opencensus.io v0.23.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.24.0
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/jaeger v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/examples v0.0.0-20210802225658-edb9b3bc2266 // indirect
	google.golang.org/protobuf v1.27.1
	gotest.tools/v3 v3.0.3
	stash.kopano.io/kgol/rndm v1.1.1
)

// this is a transitive replace. See https://github.com/libregraph/lico/blob/master/go.mod#L38
replace github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
