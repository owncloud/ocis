module github.com/owncloud/ocis/accounts

go 1.13

require (
	github.com/CiscoM31/godata v0.0.0-20191007193734-c2c4ebb1b415
	github.com/blevesearch/bleve v1.0.9
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/mennanov/fieldmask-utils v0.3.2
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/owncloud/ocis v1.0.0-rc1 // indirect
	github.com/owncloud/ocis-pkg/v2 v2.4.1-0.20200902134813-1e87c6173ada
	github.com/owncloud/ocis-settings v0.3.2-0.20200903035407-ad5de8264f91
	github.com/restic/calens v0.2.0
	github.com/rs/zerolog v1.19.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	google.golang.org/genproto v0.0.0-20200624020401-64a14ca9d1ad
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
