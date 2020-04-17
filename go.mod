module github.com/owncloud/ocis-accounts

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/joho/godotenv v1.3.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.0.0
	github.com/oklog/run v1.1.0
	github.com/owncloud/ocis-hello v0.0.0-20200114105804-61741477dcec
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/owncloud/ocis-settings v0.0.0-20200407203258-bd5da39fe8c0
	github.com/restic/calens v0.2.0
	github.com/spf13/viper v1.6.1
)

replace github.com/owncloud/ocis-settings => ../ocis-settings
