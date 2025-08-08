# Greeter Example

This example demonstrates how to use `protoc-gen-microweb` to generate HTTP handlers from protobuf service definitions.

## Overview

The greeter example shows:
- Protobuf service definition with HTTP annotations
- Generated HTTP handlers using Chi router
- JSON ↔ protobuf conversion
- REST API endpoints

## Files Generated

From `greeter.proto`, the following files are generated in the `proto/` directory:

- `greeter.pb.go` - Protobuf message definitions
- `greeter_grpc.pb.go` - gRPC service definitions  
- `greeter.pb.micro.go` - Go-micro service definitions
- `greeter.pb.web.go` - HTTP handlers (protoc-gen-microweb)

## `greeter.pb.go` (protoc-gen-go)
- **Purpose**: Core protobuf message definitions and serialization
- **Contains**: 
  - Message structs (`SayRequest`, `SayResponse`) with protobuf tags
  - Serialization/deserialization methods (`XXX_Marshal`, `XXX_Unmarshal`)
  - Getter methods for fields
  - File descriptor for protobuf reflection
- **Usage**: Base layer for all protobuf operations

## `greeter.pb.micro.go` (protoc-gen-micro)
- **Purpose**: Go-micro framework client/server code
- **Contains**:
  - `GreeterService` interface for clients
  - `GreeterHandler` interface for servers
  - `RegisterGreeterHandler()` for micro service registration
  - Client wrapper with RPC calls using micro's client library
- **Usage**: For microservices using go-micro framework

## `greeter.pb.web.go` (protoc-gen-microweb)
- **Purpose**: HTTP REST API handlers
- **Contains**:
  - `RegisterGreeterWeb()` for Chi router registration
  - HTTP handlers that convert JSON ↔ protobuf
  - REST endpoints (`POST /api/say`, `POST /api/anything`)
  - JSON marshaling/unmarshaling for protobuf messages
- **Usage**: For HTTP APIs that expose protobuf services as REST

## Differences:
- **Protocol**: `pb.go` = protobuf binary, `micro.go` = micro RPC, `web.go` = HTTP/JSON
- **Transport**: `pb.go` = none, `micro.go` = micro client/server, `web.go` = Chi HTTP router
- **Serialization**: `pb.go` = protobuf binary, `micro.go` = protobuf binary, `web.go` = JSON
- **Use case**: `pb.go` = foundation, `micro.go` = microservices, `web.go` = REST APIs

**The three files work together: `pb.go` provides the data structures, `micro.go` enables microservice communication, and `web.go` exposes HTTP REST endpoints.**


## Step-by-Step Instructions

### 1. Install Dependencies

```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest
go install github.com/owncloud/protoc-gen-microweb@latest
```

### 2. Generate Code

```bash
# Generate all protobuf files
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

### 3. Build and Run

```bash
# Build the application
go build

# Run the server
go run main.go
```

The server will start on `http://localhost:8080`

### 4. Test the API

```bash
# Test the /api/say endpoint
curl -X POST http://localhost:8080/api/say \
  -H "Content-Type: application/json" \
  -d '{"name":"test"}'

# Expected response: {"message":"Hello test!"}

# Test the /api/anything endpoint  
curl -X POST http://localhost:8080/api/anything \
  -H "Content-Type: application/json" \
  -d '{}'

# Expected response: {"message":"Saying Anything!"}

# Test with empty name (uses default)
curl -X POST http://localhost:8080/api/say \
  -H "Content-Type: application/json" \
  -d '{"name":""}'

# Expected response: {"message":"Hello World!"}
```

## Troubleshooting

### Missing Dependencies
If you get missing go.sum entries, run:
```bash
go mod tidy
```

### gRPC Version Issues
If you get gRPC API errors, update to latest:
```bash
go get google.golang.org/grpc@latest
```

### Interface Conflicts
The plugin automatically avoids naming conflicts by using `GreeterWebHandler` instead of `GreeterHandler`. 
