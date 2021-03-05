module github.com/owncloud/ocis/storage

go 1.15

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/cs3org/reva v1.6.1-0.20210305144241-2011cb557105
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/mitchellh/copystructure v1.1.1 // indirect
	github.com/oklog/run v1.1.0
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20210216094451-dc73176dc62d
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/mod v0.4.1 // indirect
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	// latest version compatible with etcd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
