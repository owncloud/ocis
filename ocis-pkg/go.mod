module github.com/owncloud/ocis/ocis-pkg

go 1.16

require (
	github.com/CiscoM31/godata v0.0.0-20201003040028-eadcd34e7f06
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210408173139-0d57213d3f5c
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
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20210614143420-5ee2eb1e7887
	github.com/cs3org/reva v1.9.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/iancoleman/strcase v0.1.3
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/oklog/run v1.1.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/owncloud/ocis/accounts v0.5.3-0.20210216094451-dc73176dc62d
	github.com/owncloud/ocis/glauth v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/graph v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/graph-explorer v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/idp v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/ocs v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/onlyoffice v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/proxy v0.0.0-20210412105747-9b95e9b1191b
	github.com/owncloud/ocis/settings v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/storage v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/store v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/thumbnails v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/web v0.0.0-20210413063522-955bd60edf33
	github.com/owncloud/ocis/webdav v0.0.0-20210413063522-955bd60edf33
	github.com/prometheus/client_golang v1.10.0
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.etcd.io/etcd/api/v3 v3.5.0-pre // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.0-pre // indirect
	go.opencensus.io v0.23.0
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
	golang.org/x/sys v0.0.0-20210608053332-aa57babbf139 // indirect
	golang.org/x/tools v0.1.2 // indirect
	google.golang.org/grpc v1.38.0
)

replace (
	// broken dependency chain for konnect v0.34.0
	github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/glauth => ../glauth
	github.com/owncloud/ocis/graph => ../graph
	github.com/owncloud/ocis/graph-explorer => ../graph-explorer
	github.com/owncloud/ocis/idp => ../idp
	github.com/owncloud/ocis/ocs => ../ocs
	github.com/owncloud/ocis/onlyoffice => ../onlyoffice
	github.com/owncloud/ocis/proxy => ../proxy
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/storage => ../storage
	github.com/owncloud/ocis/store => ../store
	github.com/owncloud/ocis/thumbnails => ../thumbnails
	github.com/owncloud/ocis/web => ../web
	github.com/owncloud/ocis/webdav => ../webdav
	// taken from https://github.com/asim/go-micro/blob/master/plugins/registry/etcd/go.mod#L14-L16
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
)
