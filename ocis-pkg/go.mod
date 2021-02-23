module github.com/owncloud/ocis/ocis-pkg

go 1.15

require (
	github.com/CiscoM31/godata v0.0.0-20201003040028-eadcd34e7f06
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/kubernetes/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/mdns/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/registry/nats/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/server/http/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/plugins/wrapper/trace/opencensus/v3 v3.0.0-20210217182006-0f0ace1a44a9
	github.com/asim/go-micro/v3 v3.5.1-0.20210217182006-0f0ace1a44a9
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20210209082852-35ace33082f5
	github.com/cs3org/reva v1.6.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5
	github.com/iancoleman/strcase v0.1.2
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/owncloud/ocis/accounts v0.5.3-0.20210216094451-dc73176dc62d
	github.com/owncloud/ocis/settings v0.0.0-20210216094451-dc73176dc62d
	github.com/owncloud/ocis/storage v0.0.0-20210216094451-dc73176dc62d
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/sony/gobreaker v0.4.1
	github.com/stretchr/testify v1.7.0
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.opencensus.io v0.22.6
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c
	google.golang.org/grpc v1.35.0
	honnef.co/go/tools v0.1.1
)

replace (
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/storage => ../storage
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	// latest version compatible with etcd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
