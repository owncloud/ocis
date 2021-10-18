SHELL := bash
NAME := store

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
endif

############ go tooling ############
include ../.make/go.mk

############ release ############
include ../.make/release.mk

############ docs generate ############
include ../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate \
			   grpc-docs-generate

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: protobuf # CI runs ci-node-generate automatically before this target
	@go generate ./...

.PHONY: ci-node-generate
ci-node-generate:


############ protobuf ############
PROTO_VERSION := v0
PROTO_SRC := pkg/proto/$(PROTO_VERSION)

include ../.make/protobuf.mk

.PHONY: protobuf
protobuf: buf-generate
