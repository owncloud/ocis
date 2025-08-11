//go:build tools
// +build tools

package main

// This file pins CLI tool dependencies in go.mod for reproducible installs.
// Tools are imported as blanks and guarded by the "tools" build tag so they
// are not included in normal builds.

import (
	_ "github.com/go-micro/generator/cmd/protoc-gen-micro"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
