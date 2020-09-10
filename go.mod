module github.com/owncloud/ocis-reva

go 1.13

require (
	github.com/cs3org/reva v1.2.1-0.20200826162318-c0f54e1f37ea
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/micro/cli/v2 v2.1.1
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.0.0
	github.com/oklog/run v1.0.0
	github.com/owncloud/flaex v0.0.0-20200411150708-dce59891a203
	github.com/owncloud/ocis-pkg/v2 v2.2.1
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.1
	gopkg.in/ini.v1 v1.51.1 // indirect
)

// ocis-sharing branch
replace github.com/cs3org/reva => github.com/butonic/reva v0.0.0-20200910112438-dd43734ae8af
