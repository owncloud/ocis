# protoc-gen-microweb

Protocol Buffers plugin that generates HTTP web handlers from protobuf service definitions. Converts gRPC HTTP annotations into Chi router handlers with JSON serialization, eliminating manual HTTP boilerplate code. Generates service interfaces, request/response marshaling, and proper error handling for REST APIs built on protobuf schemas.

The three files serve different purposes in a protobuf-based system:

## `greeter.pb.go` (protoc-gen-go)
- **Purpose**: Core protobuf message definitions and serialization
- **Usage**: Base layer for all protobuf operations

## `greeter.pb.micro.go` (protoc-gen-micro)
- **Purpose**: Go-micro framework client/server code
- **Usage**: For microservices using go-micro framework

## `greeter.pb.web.go` (protoc-gen-microweb)
- **Purpose**: HTTP REST API handlers
- **Usage**: For HTTP APIs that expose protobuf services as REST

## Differences:
- **Protocol**: `pb.go` = protobuf binary, `micro.go` = micro RPC, `web.go` = HTTP/JSON
- **Transport**: `pb.go` = none, `micro.go` = micro client/server, `web.go` = Chi HTTP router
- **Serialization**: `pb.go` = protobuf binary, `micro.go` = protobuf binary, `web.go` = JSON
- **Use case**: `pb.go` = foundation, `micro.go` = microservices, `web.go` = REST APIs

**The three files work together: `pb.go` provides the data structures, `micro.go` enables microservice communication, and `web.go` exposes HTTP REST endpoints.**


## Quick Example

```protobuf
service Greeter {
    rpc Say(SayRequest) returns (SayResponse) {
        option (google.api.http) = {
            post: "/api/say"
            body: "*"
        };
    }
}
```

```go
// Generated handler interface
type GreeterHandler interface {
    Say(ctx context.Context, in *SayRequest, out *SayResponse) error
}

// Usage
mux := chi.NewMux()
proto.RegisterGreeterWeb(mux, &Greeter{})
```

## Docs

### Run

```
# Install required tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/owncloud/protoc-gen-microweb@latest

# Generate Go code from protobuf
protoc \
	--proto_path=$(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
	--proto_path=proto/ \
	--go_out=proto/ \
	--go_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
	--go-grpc_out=proto/ \
	--go-grpc_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
	--micro_out=proto/ \
	--micro_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
	--microweb_out=proto/ \
	--microweb_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
	proto/greeter.proto
```

### Install

```
GO111MODULE=off go get -v github.com/owncloud/protoc-gen-microweb
```

### Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.12.

```bash
go get -d github.com/owncloud/protoc-gen-microweb
cd $GOPATH/src/github.com/owncloud/protoc-gen-microweb

go install
```

## Security

If you find a security issue please contact security@owncloud.com first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```
Copyright (c) 2021 ownCloud GmbH <https://owncloud.com>
```
