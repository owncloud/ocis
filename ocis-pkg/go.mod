module github.com/owncloud/ocis/ocis-pkg

go 1.13

require (
	github.com/ascarter/requestid v0.0.0-20170313220838-5b76ab3d4aee
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/cs3org/reva v1.1.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/haya14busa/goverage v0.0.0-20180129164344-eec3514a20b5
	github.com/justinas/alice v1.2.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opencensus/v2 v2.9.1
	github.com/owncloud/ocis/settings v0.0.0-20200918114005-1a0ddd2190ee // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.19.0
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	go.opencensus.io v0.22.4
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	honnef.co/go/tools v0.0.1-2020.1.5
)

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	github.com/owncloud/ocis/settings => ../settings
	)
