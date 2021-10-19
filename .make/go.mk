OCIS_REPO := github.com/owncloud/ocis
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

ifeq ($(UNAME), Darwin)
	GOBUILD ?= go build -i
else
	GOBUILD ?= go build
endif

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
		VERSION ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

LDFLAGS += -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -s -w -X "$(OCIS_REPO)/ocis-pkg/version.String=$(VERSION)" -X "$(OCIS_REPO)/ocis-pkg/version.Date=$(DATE)"
DEBUG_LDFLAGS += -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -X "$(OCIS_REPO)/ocis-pkg/version.String=$(VERSION)" -X "$(OCIS_REPO)/ocis-pkg/version.Date=$(DATE)"

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

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run -E gosec -E bodyclose -E dogsled -E durationcheck -E revive -E ifshort -E makezero -E prealloc -E predeclared --path-prefix $(NAME)

.PHONY: ci-golangci-lint
ci-golangci-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run -E gosec -E bodyclose -E dogsled -E durationcheck -E revive -E ifshort -E makezero -E prealloc -E predeclared --path-prefix $(NAME) --timeout 10m0s --issues-exit-code 0 --out-format checkstyle > checkstyle.xml

.PHONY: test
test: ci-node-generate
	@go test -v -coverprofile coverage.out ./...

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
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(DEBUG_LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

.PHONY: watch
watch: $(REFLEX)
	$(REFLEX) -c reflex.conf
