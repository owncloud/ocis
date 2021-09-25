# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.4.0. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for bingo variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BINGO)
#	@echo "Running bingo"
#	@$(BINGO) <flags/args..>
#
BINGO := $(GOBIN)/bingo-v0.5.1
$(BINGO): $(BINGO_DIR)/bingo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/bingo-v0.5.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=bingo.mod -o=$(GOBIN)/bingo-v0.5.1 "github.com/bwplotka/bingo"

BUF := $(GOBIN)/buf-v0.56.0
$(BUF): $(BINGO_DIR)/buf.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/buf-v0.56.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=buf.mod -o=$(GOBIN)/buf-v0.56.0 "github.com/bufbuild/buf/cmd/buf"

BUILDIFIER := $(GOBIN)/buildifier-v0.0.0-20210920153738-d6daef01a1a2
$(BUILDIFIER): $(BINGO_DIR)/buildifier.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/buildifier-v0.0.0-20210920153738-d6daef01a1a2"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=buildifier.mod -o=$(GOBIN)/buildifier-v0.0.0-20210920153738-d6daef01a1a2 "github.com/bazelbuild/buildtools/buildifier"

CALENS := $(GOBIN)/calens-v0.2.0
$(CALENS): $(BINGO_DIR)/calens.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/calens-v0.2.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=calens.mod -o=$(GOBIN)/calens-v0.2.0 "github.com/restic/calens"

FILEB0X := $(GOBIN)/fileb0x-v1.1.4
$(FILEB0X): $(BINGO_DIR)/fileb0x.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/fileb0x-v1.1.4"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=fileb0x.mod -o=$(GOBIN)/fileb0x-v1.1.4 "github.com/UnnoTed/fileb0x"

FLAEX := $(GOBIN)/flaex-v0.2.0
$(FLAEX): $(BINGO_DIR)/flaex.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/flaex-v0.2.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=flaex.mod -o=$(GOBIN)/flaex-v0.2.0 "github.com/owncloud/flaex"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v1.42.1
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v1.42.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v1.42.1 "github.com/golangci/golangci-lint/cmd/golangci-lint"

GOX := $(GOBIN)/gox-v1.0.1
$(GOX): $(BINGO_DIR)/gox.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/gox-v1.0.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=gox.mod -o=$(GOBIN)/gox-v1.0.1 "github.com/mitchellh/gox"

HUGO := $(GOBIN)/hugo-v0.88.1
$(HUGO): $(BINGO_DIR)/hugo.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/hugo-v0.88.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=hugo.mod -o=$(GOBIN)/hugo-v0.88.1 "github.com/gohugoio/hugo"

MUTAGEN := $(GOBIN)/mutagen-v0.11.8
$(MUTAGEN): $(BINGO_DIR)/mutagen.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/mutagen-v0.11.8"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=mutagen.mod -o=$(GOBIN)/mutagen-v0.11.8 "github.com/mutagen-io/mutagen/cmd/mutagen"

OAPI_CODEGEN := $(GOBIN)/oapi-codegen-v1.8.2
$(OAPI_CODEGEN): $(BINGO_DIR)/oapi-codegen.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/oapi-codegen-v1.8.2"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=oapi-codegen.mod -o=$(GOBIN)/oapi-codegen-v1.8.2 "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"

PROTOC_GEN_DOC := $(GOBIN)/protoc-gen-doc-v1.5.0
$(PROTOC_GEN_DOC): $(BINGO_DIR)/protoc-gen-doc.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-doc-v1.5.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-doc.mod -o=$(GOBIN)/protoc-gen-doc-v1.5.0 "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go-v1.27.1
$(PROTOC_GEN_GO): $(BINGO_DIR)/protoc-gen-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-v1.27.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-go.mod -o=$(GOBIN)/protoc-gen-go-v1.27.1 "google.golang.org/protobuf/cmd/protoc-gen-go"

PROTOC_GEN_MICRO := $(GOBIN)/protoc-gen-micro-v3.0.0-20210924081004-8c39b1e1204d
$(PROTOC_GEN_MICRO): $(BINGO_DIR)/protoc-gen-micro.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-micro-v3.0.0-20210924081004-8c39b1e1204d"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-micro.mod -o=$(GOBIN)/protoc-gen-micro-v3.0.0-20210924081004-8c39b1e1204d "github.com/asim/go-micro/cmd/protoc-gen-micro/v3"

PROTOC_GEN_MICROWEB := $(GOBIN)/protoc-gen-microweb-v0.0.0-20210824101557-828409dbfbf9
$(PROTOC_GEN_MICROWEB): $(BINGO_DIR)/protoc-gen-microweb.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-microweb-v0.0.0-20210824101557-828409dbfbf9"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-microweb.mod -o=$(GOBIN)/protoc-gen-microweb-v0.0.0-20210824101557-828409dbfbf9 "github.com/owncloud/protoc-gen-microweb"

PROTOC_GEN_OPENAPIV2 := $(GOBIN)/protoc-gen-openapiv2-v2.6.0
$(PROTOC_GEN_OPENAPIV2): $(BINGO_DIR)/protoc-gen-openapiv2.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-openapiv2-v2.6.0"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=protoc-gen-openapiv2.mod -o=$(GOBIN)/protoc-gen-openapiv2-v2.6.0 "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"

REFLEX := $(GOBIN)/reflex-v0.3.1
$(REFLEX): $(BINGO_DIR)/reflex.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/reflex-v0.3.1"
	@cd $(BINGO_DIR) && $(GO) build -mod=mod -modfile=reflex.mod -o=$(GOBIN)/reflex-v0.3.1 "github.com/cespare/reflex"

