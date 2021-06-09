module github.com/owncloud/ocis/graph

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cs3org/go-cs3apis v0.0.0-20210527092509-2b828e94ed4c
	github.com/cs3org/reva v1.7.1-0.20210531093513-b74a2b156af6
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-ldap/ldap/v3 v3.3.0
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/prometheus/client_golang v1.10.0
	github.com/spf13/viper v1.7.1
	github.com/thejerf/suture/v4 v4.0.0
	github.com/yaegashi/msgraph.go v0.1.4
	go.opencensus.io v0.23.0
	google.golang.org/grpc v1.38.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/store => ../store
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
)
