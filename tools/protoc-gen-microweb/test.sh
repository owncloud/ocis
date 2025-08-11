#!/bin/bash

# Test for protoc-gen-microweb
# Tests complete generation and usage pipeline

set -e

assert_file_exists() {
    if [ ! -f "$1" ]; then
        echo "File $1 does not exist"
        exit 1
    fi
}

setup_prerequisites() {
  echo "0. Checking and installing protoc tools..."
  which protoc || (echo "protoc not found" && exit 1)

  # Install compatible protoc toolchain (pinned)
  echo "Installing compatible protoc-gen-go version..."
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
  which protoc-gen-go-grpc || (echo "Installing protoc-gen-go-grpc..." && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1)
  which protoc-gen-micro || (echo "Installing protoc-gen-micro..." && go install github.com/go-micro/generator/cmd/protoc-gen-micro@v1.0.0)
  which protoc-gen-go || (echo "Installing protoc-gen-go..." && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6)
  export PATH="$PATH:$(go env GOPATH)/bin"

  # Build the protoc-gen-microweb tool
  echo "Building protoc-gen-microweb tool..."
  GOWORK=off go build -o protoc-gen-microweb .

  # Clone googleapis (unpinned HEAD). For full determinism, pin a commit.
  echo "    Setting up googleapis..."
  if [ ! -d "/tmp/googleapis" ]; then
    echo "Cloning googleapis to /tmp/googleapis..."
    git clone https://github.com/googleapis/googleapis.git /tmp/googleapis
  else
    echo "googleapis already exists at /tmp/googleapis"
  fi
}

main() {
  echo "=== protoc-gen-microweb Test ==="
  
  # Ensure we're in the protoc-gen-microweb directory
  cd "$(dirname "$0")"
  echo "Working directory: $(pwd)"
  
  # Prerequisites and environment
  setup_prerequisites
  
  echo "2. Assert initial example state"
  # Use greeter example
  cd examples/greeter
  
  # Clean up any existing generated files
  rm -f proto/greeter.pb.go
  rm -f proto/greeter.pb.web.go
  rm -f proto/greeter.pb.micro.go
  
  assert_file_exists proto/greeter.proto
  assert_file_exists main.go
  assert_file_exists go.mod
  assert_file_exists go.sum

  # Ensure we're using the correct protobuf version (no tidy/download for determinism)
  echo "3. Ensuring correct protobuf versions... (no go mod tidy/download)"

  # Generate code - microweb creates flat structure regardless of module parameter
  echo "4. Generate code..."
  protoc \
    --proto_path=/tmp/googleapis \
    --proto_path=proto/ \
    --go_out=proto/ \
    --go_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
    --go-grpc_out=proto/ \
    --go-grpc_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
    --micro_out=proto/ \
    --micro_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
    --microweb_out=proto/ \
    --microweb_opt=module=github.com/owncloud/protoc-gen-microweb/examples/greeter/proto \
    --plugin=protoc-gen-microweb=../../protoc-gen-microweb \
    proto/greeter.proto

  # 4.5. Skipping post-generation tidy/download; rely on pinned go.mod/go.sum

  # Assert generated files are in the correct location
  assert_file_exists proto/greeter.proto
  assert_file_exists proto/greeter.pb.go
  assert_file_exists proto/greeter.pb.web.go
  assert_file_exists proto/greeter.pb.micro.go

  echo "5. Run and test server..."
  # Use GOWORK=off to run in standalone mode; -mod=readonly forbids go.mod edits
  GOWORK=off go build -mod=readonly -o greeter .
  GOWORK=off ./greeter &
  SERVER_PID=$!

  # Wait for server to start
  sleep 2
  
  echo "6. Test POST /api/say..."
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/api/say \
    -H "Content-Type: application/json" \
    -d '{"name":"test"}')

  # Extract response body and status code
  RESPONSE_BODY=$(echo "$RESPONSE" | head -n 1)
  STATUS_CODE=$(echo "$RESPONSE" | tail -n 1)

  kill $SERVER_PID 2>/dev/null || true

  echo "7. Validate results..."
  echo "Response: $RESPONSE_BODY"
  echo "Status: $STATUS_CODE"

  if [ "$RESPONSE_BODY" = '{"message":"Hello test!"}' ] && [ "$STATUS_CODE" = "201" ]; then
      echo "PASS"
      exit 0
  else
      echo "Expected: {\"message\":\"Hello test!\"} with status 201"
      echo "Got: $RESPONSE_BODY with status $STATUS_CODE"
      echo "FAIL"
      exit 1
  fi
}

main

