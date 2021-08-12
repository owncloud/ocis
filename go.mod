module github.com/owncloud/ocis

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/CiscoM31/godata v0.0.0-20201003040028-eadcd34e7f06
	github.com/GeertJohan/yubigo v0.0.0-20190917122436-175bc097e60e
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/logger/zerolog/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/registry/kubernetes/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/registry/mdns/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/registry/nats/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/server/http/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/plugins/wrapper/trace/opencensus/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/blevesearch/bleve v1.0.9
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20210802070913-970eec344e59
	github.com/cs3org/reva v1.11.1-0.20210812105259-756bdced1d22
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/disintegration/imaging v1.6.2
	github.com/glauth/glauth v1.1.3-0.20210729125545-b9aecdfcac31
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-logr/logr v0.4.0
	github.com/go-ozzo/ozzo-validation/v4 v4.2.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.5.2
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0
	github.com/iancoleman/strcase v0.1.3
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/justinas/alice v1.2.0
	github.com/libregraph/lico v0.34.1-0.20210803054646-b584e0372224
	github.com/mennanov/fieldmask-utils v0.3.3
	github.com/micro/cli/v2 v2.1.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484
	github.com/nmcclain/ldap v0.0.0-20210720162743-7f8d1e44eeba
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/open-graph-api-go v0.0.0-20210511151655-57894f7d46fb
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/rs/zerolog v1.23.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c // indirect
	github.com/thejerf/suture/v4 v4.0.1
	github.com/yaegashi/msgraph.go v0.1.4
	go.etcd.io/etcd/pkg/v3 v3.5.0-pre // indirect
	go.opencensus.io v0.23.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	golang.org/x/tools v0.1.2 // indirect
	google.golang.org/genproto v0.0.0-20210402141018-6c239bbf2bb1
	google.golang.org/grpc v1.39.1
	google.golang.org/grpc/examples v0.0.0-20210802225658-edb9b3bc2266 // indirect
	google.golang.org/protobuf v1.27.1
	gotest.tools v2.2.0+incompatible
	stash.kopano.io/kgol/rndm v1.1.0
)

replace (
	github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
)
