module github.com/owncloud/ocis/ocis-pkg

go 1.13

require (
	github.com/CiscoM31/godata v0.0.0-20201003040028-eadcd34e7f06
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/go-cs3apis v0.0.0-20201118090759-87929f5bae21
	github.com/cs3org/reva v1.4.1-0.20201216095137-5207197217f1
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5
	github.com/iancoleman/strcase v0.1.2
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opencensus/v2 v2.9.1
	github.com/owncloud/ocis/accounts v0.5.3-0.20201103104733-ff2c41028d9b
	github.com/owncloud/ocis/settings v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/storage v0.0.0-20201015120921-38358ba4d4df
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.6.1
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.opencensus.io v0.22.5
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.33.2
	honnef.co/go/tools v0.0.1-2020.1.5
)

replace (
	// pr https://github.com/cs3org/reva/pull/1368
	github.com/cs3org/reva => github.com/butonic/reva v0.0.0-20201214153639-b976591161c9
	github.com/owncloud/ocis/accounts => ../accounts
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/storage => ../storage
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
