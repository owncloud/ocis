SHELL := bash
NAME := ocis-pkg
IMPORT := github.com/owncloud/ocis/$(NAME)
BIN := bin
DIST := dist

ifeq ($(OS), Windows_NT)
	UNAME := Windows
else
	UNAME := $(shell uname -s)
endif

ifeq ($(UNAME), Darwin)
	GOBUILD ?= go build -i
else
	GOBUILD ?= go build
endif

PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)
GENERATE ?= $(PACKAGES)

TAGS ?=

LDFLAGS += -s -w
GCFLAGS += all=-N -l

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean:
	go clean -i ./...
	rm -rf $(BIN) $(DIST)

.PHONY: go-mod-tidy
go-mod-tidy:
	@go mod tidy

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: staticcheck
staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck -tags '$(TAGS)' $(PACKAGES)

.PHONY: lint
lint:
	for PKG in $(PACKAGES); do go run golang.org/x/lint/golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	go run github.com/haya14busa/goverage -v -coverprofile coverage.out $(PACKAGES)

.PHONY: go-coverage
go-coverage:
	@if [ ! -f coverage.out ]; then $(MAKE) test  &>/dev/null; fi;
	@go tool cover -func coverage.out | tail -1 | grep -Eo "[0-9]+\.[0-9]+"

.PHONY: build
build:
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./...

############ docs generate ############

.PHONY: docs-generate
docs-generate:

############ generate ############
.PHONY: generate
generate: ci-node-generate ci-go-generate

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target
	go generate $(GENERATE)

.PHONY: ci-node-generate
ci-node-generate:
