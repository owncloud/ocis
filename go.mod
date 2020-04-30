module github.com/owncloud/ocis-accounts

go 1.13

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/go-test/deep v1.0.6 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.1
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.4 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.0.0
	github.com/oklog/run v1.1.0
	github.com/owncloud/ocis-pkg/v2 v2.2.1
	github.com/restic/calens v0.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/viper v1.6.3
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84
	google.golang.org/protobuf v1.22.0
	honnef.co/go/tools v0.0.1-2020.1.0.20200427215036-cd1ad299aeab // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0