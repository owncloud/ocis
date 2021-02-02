SHELL := bash
NAME := idp
IMPORT := github.com/owncloud/ocis/$(NAME)
BIN := bin
DIST := dist

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

# Assets
LOGO_URL = https://raw.githubusercontent.com/owncloud/assets/main/logo.svg
FAVICON_URL = https://raw.githubusercontent.com/owncloud/assets/main/favicon.ico

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
DEBUG_LDFLAGS += -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"
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
	rm -rf node_modules
	rm -rf assets/identifier

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

.PHONY: generate
generate: assets
	go generate $(GENERATE)

# TODO find a docker container with go and yarn so we can properly build assets in ci
#assets: assets/identifier/static assets/identifier/index.html
.PHONY: assets
assets: assets/identifier/static assets/identifier/static/logo.svg assets/identifier/static/favicon.ico

assets/identifier/static:
	mkdir -p assets/identifier/static

assets/identifier/static/logo.svg:
	curl -o assets/identifier/static/logo.svg ${LOGO_URL}

assets/identifier/static/favicon.ico:
	curl -o assets/identifier/static/favicon.ico ${FAVICON_URL}

assets/identifier/index.html: node_modules
	yarn build

node_modules:
	yarn install

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
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(DEBUG_LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

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

.PHONY: config-docs-generate
config-docs-generate:
	go run github.com/owncloud/flaex >| ../docs/extensions/$(NAME)/configuration.md

.PHONY: watch
watch:
	go run github.com/cespare/reflex -c reflex.conf

# $(GOPATH)/bin/protoc-gen-go:
# 	GO111MODULE=off go get -v github.com/golang/protobuf/protoc-gen-go

# $(GOPATH)/bin/protoc-gen-micro:
# 	GO111MODULE=off go get -v github.com/micro/protoc-gen-micro

# $(GOPATH)/bin/protoc-gen-microweb:
# 	GO111MODULE=off go get -v github.com/webhippie/protoc-gen-microweb

# $(GOPATH)/bin/protoc-gen-swagger:
# 	GO111MODULE=off go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

# pkg/proto/v0/example.pb.go: pkg/proto/v0/example.proto
# 	protoc \
# 		-I=third_party/ \
# 		-I=pkg/proto/v0/ \
# 		--go_out=pkg/proto/v0 example.proto

# pkg/proto/v0/example.pb.micro.go: pkg/proto/v0/example.proto
# 	protoc \
# 		-I=third_party/ \
# 		-I=pkg/proto/v0/ \
# 		--micro_out=pkg/proto/v0 example.proto

# pkg/proto/v0/example.pb.web.go: pkg/proto/v0/example.proto
# 	protoc \
# 		-I=third_party/ \
# 		-I=pkg/proto/v0/ \
# 		--microweb_out=pkg/proto/v0 example.proto

# pkg/proto/v0/example.swagger.json: pkg/proto/v0/example.proto
# 	protoc \
# 		-I=third_party/ \
# 		-I=pkg/proto/v0/ \
# 		--swagger_out=pkg/proto/v0 example.proto

# .PHONY: protobuf
# protobuf:  $(GOPATH)/bin/protoc-gen-go $(GOPATH)/bin/protoc-gen-micro $(GOPATH)/bin/protoc-gen-microweb $(GOPATH)/bin/protoc-gen-swagger pkg/proto/v0/example.pb.go pkg/proto/v0/example.pb.micro.go pkg/proto/v0/example.pb.web.go pkg/proto/v0/example.swagger.json
