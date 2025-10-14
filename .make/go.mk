OCIS_REPO := github.com/owncloud/ocis/v2
IMPORT := ($OCIS_REPO)/$(NAME)
BIN := bin
DIST := dist

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(NAME).exe
	UNAME := Windows
else
	EXECUTABLE := $(NAME)
	UNAME := $(shell uname -s)
endif

GOBUILD ?= go build

SOURCES ?= $(shell find . -name "*.go" -type f -not -path "./node_modules/*")

TAGS ?=

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
		STRING ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

LDFLAGS += -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -s -w -X "$(OCIS_REPO)/ocis-pkg/version.String=$(STRING)" -X "$(OCIS_REPO)/ocis-pkg/version.Tag=$(VERSION)" -X "$(OCIS_REPO)/ocis-pkg/version.Date=$(DATE)"
DEBUG_LDFLAGS += -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -X "$(OCIS_REPO)/ocis-pkg/version.String=$(STRING)" -X "$(OCIS_REPO)/ocis-pkg/version.Tag=$(VERSION)" -X "$(OCIS_REPO)/ocis-pkg/version.Date=$(DATE)"

GCFLAGS += all=-N -l

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean:
	@echo "$(NAME): clean"
	go clean -i ./...
	rm -rf $(BIN) $(DIST)

.PHONY: go-mod-tidy
go-mod-tidy:
	@echo "$(NAME): go-mod-tidy"
	@go mod tidy

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: golangci-lint-fix
golangci-lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run $(LINTERS) --fix

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --path-prefix services/$(NAME)

.PHONY: ci-golangci-lint
ci-golangci-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --path-prefix services/$(NAME) --timeout 15m0s --issues-exit-code 0 --out-format checkstyle > checkstyle.xml

.PHONY: test
test:
	@go test -v -tags '$(TAGS)' -coverprofile coverage.out ./...

.PHONY: go-coverage
go-coverage:
	@if [ ! -f coverage.out ]; then $(MAKE) test  &>/dev/null; fi;
	@go tool cover -func coverage.out | tail -1 | grep -Eo "[0-9]+\.[0-9]+"

.PHONY: install
install: $(SOURCES)
	go install -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/$(NAME)

.PHONY: build-all
build-all: build build-debug

.PHONY: build
build: $(BIN)/$(EXECUTABLE)

.PHONY: build-debug
build-debug: $(BIN)/$(EXECUTABLE)-debug

$(BIN)/$(EXECUTABLE): $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(EXECUTABLE)-debug: $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(DEBUG_LDFLAGS)' -cover -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

.PHONY: watch
watch: $(REFLEX)
	$(REFLEX) -c reflex.conf

debug-linux-docker-amd64: release-dirs
	GOOS=linux \
	GOARCH=amd64 \
	go build \
        -gcflags="all=-N -l" \
		-tags 'netgo $(TAGS)' \
		-buildmode=exe \
		-trimpath \
		-ldflags '-extldflags "-static" $(DEBUG_LDFLAGS) $(DOCKER_LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-linux-amd64' \
		./cmd/$(NAME)

debug-linux-docker-arm64: release-dirs
	GOOS=linux \
	GOARCH=arm64 \
	go build \
        -gcflags="all=-N -l" \
		-tags 'netgo $(TAGS)' \
		-buildmode=exe \
		-trimpath \
		-ldflags '-extldflags "-static" $(DEBUG_LDFLAGS) $(DOCKER_LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-linux-arm64' \
		./cmd/$(NAME)
