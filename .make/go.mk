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

#LINTERS += -E gofmt -E gci -E gofumpt -E goimports
#LINTERS += -E gosec -E bodyclose -E dogsled -E durationcheck
#LINTERS += -E revive -E ifshort -E makezero -E prealloc -E predeclared
#LINTERS += -E asciicheck -E bidichk -E cyclop

#Linters presets:
#bugs: asciicheck, bidichk, bodyclose, contextcheck, durationcheck, errcheck, errorlint, exhaustive, exportloopref, gosec, govet, makezero, nilerr, noctx, rowserrcheck, scopelint, sqlclosecheck, staticcheck, typecheck
#comment: godot, godox, misspell
#complexity: cyclop, funlen, gocognit, gocyclo, nestif
#error: errcheck, errorlint, goerr113, wrapcheck
#format: gci, gofmt, gofumpt, goimports
#import: depguard, gci, goimports, gomodguard
#metalinter: gocritic, govet, revive, staticcheck
#module: depguard, gomoddirectives, gomodguard
#performance: bodyclose, maligned, noctx, prealloc
#sql: rowserrcheck, sqlclosecheck
#style: asciicheck, depguard, dogsled, dupl, errname, exhaustivestruct, forbidigo, forcetypeassert, gochecknoglobals, gochecknoinits, goconst, gocritic, godot, godox, goerr113, goheader, golint, gomnd, gomoddirectives, gomodguard, goprintffuncname, gosimple, ifshort, importas, interfacer, ireturn, lll, makezero, misspell, nakedret, nilnil, nlreturn, nolintlint, paralleltest, predeclared, promlinter, revive, stylecheck, tagliatelle, tenv, testpackage, thelper, tparallel, unconvert, varnamelen, wastedassign, whitespace, wrapcheck, wsl
#test: exhaustivestruct, paralleltest, testpackage, tparallel
#unused: deadcode, ineffassign, structcheck, unparam, unused, varcheck

LINTERS += -p bugs
LINTERS += -p comment
LINTERS += -p error -D wrapcheck
LINTERS += -p format
LINTERS += -p complexity -D funlen -D cyclop
LINTERS += -p import
LINTERS += -p metalinter
LINTERS += -p module
LINTERS += -p performance
LINTERS += -p sql
LINTERS += -p style -D dupl -D goconst -D godox -D lll -D gomnd -D tagliatelle -D varnamelen -D ireturn -D nlreturn -D gochecknoglobals
LINTERS += -p test -D exhaustivestruct
LINTERS += -p unused
LINTERS += -D interfacer # depreciated
LINTERS += -D scopelint -E exportloopref # depreciated with replacement
LINTERS += -D golint -E revive # depreciated with replacement
LINTERS += -D maligned # depreciated and replaced by govet

.PHONY: golangci-lint-fix
golangci-lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run $(LINTERS) --fix

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run $(LINTERS) --path-prefix $(NAME)

.PHONY: ci-golangci-lint
ci-golangci-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run $(LINTERS) --path-prefix $(NAME) --timeout 10m0s --issues-exit-code 0 --out-format checkstyle > checkstyle.xml

.PHONY: test
test:
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
