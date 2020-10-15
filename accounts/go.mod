module github.com/owncloud/ocis/accounts

go 1.13

require (
	github.com/cs3org/go-cs3apis v0.0.0-20200810113633-b00aca449666
	github.com/cs3org/reva v1.2.2-0.20200924071957-e6676516e61e
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/gopherjs/gopherjs v0.0.0-20190910122728-9d188e94fb99 // indirect
	github.com/mennanov/fieldmask-utils v0.3.2
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/settings v0.0.0-20200918114005-1a0ddd2190ee
	github.com/owncloud/ocis/storage v0.0.0-20201015120921-38358ba4d4df
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tredoe/osutil v1.0.5
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	google.golang.org/genproto v0.0.0-20200624020401-64a14ca9d1ad
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/settings => ../settings
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
