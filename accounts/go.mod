module github.com/owncloud/ocis/accounts

go 1.13

require (
	github.com/cs3org/go-cs3apis v0.0.0-20201118090759-87929f5bae21
	github.com/cs3org/reva v1.4.1-0.20201127111856-e6a6212c1b7b
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/mennanov/fieldmask-utils v0.3.3
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20201103111659-46bf133a3c63
	github.com/owncloud/ocis/settings v0.0.0-20200918114005-1a0ddd2190ee
	github.com/prometheus/client_golang v1.7.1
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	google.golang.org/genproto v0.0.0-20200624020401-64a14ca9d1ad
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)

replace (
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	github.com/owncloud/ocis/settings => ../settings
	github.com/owncloud/ocis/storage => ../storage
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
