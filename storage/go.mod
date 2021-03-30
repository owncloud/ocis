module github.com/owncloud/ocis/storage

go 1.16

require (
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cs3org/reva v1.6.1-0.20210326165326-e8a00d9b2368
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210329140112-169962b4cfc1
	github.com/spf13/viper v1.7.1
	github.com/thejerf/suture/v4 v4.0.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/store => ../store
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	// latest version compatible with etcd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
