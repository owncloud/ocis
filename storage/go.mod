module github.com/owncloud/ocis/storage

go 1.13

require (
	github.com/cs3org/reva v1.4.1-0.20201216095137-5207197217f1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.1.0
	github.com/owncloud/flaex v0.0.0-20200411150708-dce59891a203
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.7.0
)

replace (
	// pr https://github.com/cs3org/reva/pull/1368
	github.com/cs3org/reva => github.com/butonic/reva v0.0.0-20201221220454-a2a2f2bb2ab5
	github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
