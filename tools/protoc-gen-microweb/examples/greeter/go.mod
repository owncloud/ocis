module github.com/owncloud/protoc-gen-microweb/examples/greeter

go 1.23.0

toolchain go1.24.0

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/go-chi/render v1.0.3
	github.com/golang/protobuf v1.5.4
	go-micro.dev/v4 v4.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/api v0.0.0-20250528174236-200df99c418a
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/evanphx/json-patch/v5 v5.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/go-acme/lego/v4 v4.4.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)

replace go-micro.dev/v4 => github.com/micro/go-micro/v4 v4.11.0
