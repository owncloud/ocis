module github.com/owncloud/ocis/thumbnails

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cs3org/go-cs3apis v0.0.0-20210507060801-f176760d55f4
	github.com/cs3org/reva v1.7.1-0.20210507160327-e2c3841d0dbc
	github.com/disintegration/imaging v1.6.2
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/thejerf/suture/v4 v4.0.0
	go.opencensus.io v0.23.0
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/store => ../store
	github.com/cs3org/go-cs3apis => github.com/refs/go-cs3apis v0.0.0-20210520134557-97cdc479815b
    github.com/cs3org/reva => github.com/butonic/reva v0.0.0-20210526094447-9dfa8ec812ac
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	// latest version compatible with etcd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
