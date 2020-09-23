module github.com/owncloud/ocis/ocis-reva

go 1.13

require (
	github.com/cs3org/reva v1.2.1-0.20200911111727-51649e37df2d
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/oklog/run v1.0.0
	github.com/owncloud/flaex v0.0.0-20200411150708-dce59891a203
	github.com/owncloud/ocis/ocis-pkg v0.0.0-20200918114005-1a0ddd2190ee
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.3
	gopkg.in/ini.v1 v1.51.1 // indirect
)

replace github.com/owncloud/ocis/ocis-pkg => ../ocis-pkg
