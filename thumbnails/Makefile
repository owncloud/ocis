SHELL := bash
NAME := thumbnails
IMPORT := github.com/owncloud/ocis/$(NAME)
BIN := bin
DIST := dist
HUGO := hugo
PROTO_VERSION := v0
PROTO_SRC := pkg/proto/$(PROTO_VERSION)

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(NAME).exe
	UNAME := Windows
else
	EXECUTABLE := $(NAME)
	UNAME := $(shell uname -s)
endif

ifeq ($(UNAME), Darwin)
	GOBUILD ?= go build -i
else
	GOBUILD ?= go build
endif

PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f -not -path "./node_modules/*")
GENERATE ?= $(PACKAGES)

TAGS ?=

ifndef GOPATH
	export GOPATH := $(shell go env GOPATH)
endif
export PATH := $(PATH):$(GOPATH)/bin

ifndef OUTPUT
	ifneq ($(DRONE_TAG),)
		OUTPUT ?= $(subst v,,$(DRONE_TAG))
	else
		OUTPUT ?= testing
	endif
endif

ifndef VERSION
	ifneq ($(DRONE_TAG),)
		VERSION ?= $(subst v,,$(DRONE_TAG))
	else
		VERSION ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

LDFLAGS += -s -w -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"
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

.PHONY: changelog
changelog:
	go run github.com/restic/calens >| CHANGELOG.md

.PHONY: test
test:
	go run github.com/haya14busa/goverage -v -coverprofile coverage.out $(PACKAGES)

.PHONY: go-coverage
go-coverage:
	@if [ ! -f coverage.out ]; then $(MAKE) test  &>/dev/null; fi;
	@go tool cover -func coverage.out | tail -1 | grep -Eo "[0-9]+\.[0-9]+"

.PHONY: install
install: $(SOURCES)
	go install -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/$(NAME)

.PHONY: build
build: $(BIN)/$(EXECUTABLE) $(BIN)/$(EXECUTABLE)-debug

$(BIN)/$(EXECUTABLE): $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(EXECUTABLE)-debug: $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

.PHONY: release
release: release-dirs release-linux release-windows release-darwin release-copy release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-linux
release-linux: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'linux' -arch 'amd64 386 arm64 arm' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-windows
release-windows: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'windows' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-darwin
release-darwin: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '$(LDFLAGS)' -os 'darwin' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

.PHONY: release-finish
release-finish: release-copy release-check

.PHONY: watch
watch:
	go run github.com/cespare/reflex -c reflex.conf

############ docs generate ############

.PHONY: docs-generate
docs-generate: config-docs-generate \
			   grpc-docs-generate

.PHONY: config-docs-generate
config-docs-generate:
	go run github.com/owncloud/flaex >| ../docs/extensions/$(NAME)/configuration.md

.PHONY: grpc-docs-generate
grpc-docs-generate: ../docs/extensions/${NAME}/grpc.md

############ generate ############
.PHONY: generate
generate: ci-node-generate ci-go-generate

.PHONY: ci-go-generate
ci-go-generate: protobuf # CI runs ci-node-generate automatically before this target
	go generate $(GENERATE)

.PHONY: ci-node-generate
ci-node-generate:

############ protobuf ############
$(GOPATH)/bin/protoc-gen-go:
	go get -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0

$(GOPATH)/bin/protoc-gen-micro:
	GO111MODULE=on go get -v github.com/asim/go-micro/cmd/protoc-gen-micro/v3

$(GOPATH)/bin/protoc-gen-microweb:
	GO111MODULE=off go get -v github.com/owncloud/protoc-gen-microweb

$(GOPATH)/bin/protoc-gen-openapiv2:
	GO111MODULE=off go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2

$(GOPATH)/bin/protoc-gen-doc:
	GO111MODULE=off go get -v github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

.PHONY: $(PROTO_SRC)/${NAME}.pb.go
$(PROTO_SRC)/${NAME}.pb.go: $(GOPATH)/bin/protoc-gen-openapiv2 $(GOPATH)/bin/protoc-gen-go
	protoc \
		-I=../third_party/ \
		-I=$(PROTO_SRC)/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
		--go_out=. ${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.pb.micro.go
$(PROTO_SRC)/${NAME}.pb.micro.go: $(GOPATH)/bin/protoc-gen-openapiv2 $(GOPATH)/bin/protoc-gen-micro
	protoc \
		-I=../third_party/ \
		-I=$(PROTO_SRC)/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
		--micro_out=. ${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.pb.web.go
$(PROTO_SRC)/${NAME}.pb.web.go: $(GOPATH)/bin/protoc-gen-openapiv2 $(GOPATH)/bin/protoc-gen-microweb
	protoc \
		-I=../third_party/ \
		-I=$(PROTO_SRC)/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
		--microweb_out=. ${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.swagger.json
$(PROTO_SRC)/${NAME}.swagger.json: $(GOPATH)/bin/protoc-gen-openapiv2
	protoc \
		-I=../third_party/ \
		-I=$(PROTO_SRC)/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
		--openapiv2_out=$(PROTO_SRC)/ ${NAME}.proto

.PHONY: ../docs/extensions/${NAME}/grpc.md
../docs/extensions/${NAME}/grpc.md: $(GOPATH)/bin/protoc-gen-openapiv2 $(GOPATH)/bin/protoc-gen-doc
	protoc \
		-I=../third_party/ \
		-I=$(PROTO_SRC)/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
		--doc_opt=./templates/GRPC.tmpl,grpc.md \
		--doc_out=../docs/extensions/${NAME} $(PROTO_SRC)/${NAME}.proto

.PHONY: protobuf
protobuf: $(PROTO_SRC)/${NAME}.pb.go \
		  $(PROTO_SRC)/${NAME}.pb.micro.go \
		  $(PROTO_SRC)/${NAME}.swagger.json \
		  ../docs/extensions/${NAME}/grpc.md
		  #$(PROTO_SRC)/${NAME}.pb.web.go \ # not needed for thumbnails
