module github.com/owncloud/ocis/accounts

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cs3org/go-cs3apis v0.0.0-20210527092509-2b828e94ed4c
	github.com/cs3org/reva v1.7.1-0.20210531093513-b74a2b156af6
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0
	github.com/mennanov/fieldmask-utils v0.3.3
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210412113235-982264811ecb
	github.com/owncloud/ocis/settings v0.0.0-20210413063522-955bd60edf33
	github.com/prometheus/client_golang v1.10.0
	github.com/rs/zerolog v1.22.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/thejerf/suture/v4 v4.0.0
	go.opencensus.io v0.23.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/cs3org/reva => ../../reva
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/store => ../store
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
)
