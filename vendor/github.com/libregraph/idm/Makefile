PACKAGE  = github.com/libregraph/idm
PACKAGE_NAME = libregraph-$(shell basename $(PACKAGE))

# Tools

GO      ?= go
GOFMT   ?= gofmt
GOLINT  ?= golangci-lint

GO2XUNIT ?= go2xunit

CHGLOG  ?= git-chglog
CURL    ?= curl

# Cgo

CGO_ENABLED ?= 1

# Go modules

GO111MODULE ?= on

# Variables

export CGO_ENABLED GO111MODULE
unexport GOPATH

ARGS    ?=
PWD     := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null | sed 's/^v//' || \
			cat $(CURDIR)/.version 2> /dev/null || echo 0.0.0-unreleased)
PKGS     = $(or $(PKG),$(shell $(GO) list -mod=readonly ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell $(GO) list -mod=readonly -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS) 2>/dev/null)
CMDS     = $(or $(CMD),$(addprefix cmd/,$(notdir $(shell find "$(PWD)/cmd/" -maxdepth 1 -type d))))
TIMEOUT  = 30
BUILD_TAGS ?=

# Build

.PHONY: all
all: fmt | $(CMDS) $(PLUGINS)

plugins: fmt | $(PLUGINS)

.PHONY: $(CMDS)
$(CMDS): vendor ; $(info building $@ ...) @
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
		-mod=vendor \
		-trimpath \
		-tags "release $(BUILD_TAGS)" \
		-buildmode=exe \
		-ldflags '-s -w -buildid=reproducible/$(VERSION) -X $(PACKAGE)/version.Version=$(VERSION) -X $(PACKAGE)/version.BuildDate=$(DATE) -extldflags -static' \
		-o bin/$(notdir $@) ./$@

# Helpers

.PHONY: lint
lint: vendor ; $(info running $(GOLINT) ...)	@
	$(GOLINT) run

.PHONY: lint-checkstyle
lint-checkstyle: vendor ; $(info running $(GOLINT) checkstyle ...)     @
	@mkdir -p test
	$(GOLINT) run --out-format checkstyle --issues-exit-code 0 > test/tests.lint.xml

.PHONY: fmt
fmt: ; $(info running gofmt ...)	@
	@ret=0 && for d in $$($(GO) list -mod=readonly -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

.PHONY: check
check: ; $(info checking dependencies ...) @
	@$(GO) mod verify && echo OK

# Tests

TEST_TARGETS := test-default test-bench test-short test-race test-verbose
.PHONY: $(TEST_TARGETS)
test-bench:   ARGS=-run=_Bench* -test.benchmem -bench=.
test-short:   ARGS=-short
test-race:    ARGS=-race
test-race:    CGO_ENABLED=1
test-verbose: ARGS=-v
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test

.PHONY: test
test: ; $(info running $(NAME:%=% )tests ...)	@
	@CGO_ENABLED=$(CGO_ENABLED) $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

TEST_XML_TARGETS := test-xml-default test-xml-short test-xml-race
.PHONY: $(TEST_XML_TARGETS)
test-xml-short: ARGS=-short
test-xml-race:  ARGS=-race
test-xml-race:  CGO_ENABLED=1
$(TEST_XML_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_XML_TARGETS): test-xml

.PHONY: test-xml
test-xml: ; $(info running $(NAME:%=% )tests ...)	@
	@mkdir -p test
	2>&1 CGO_ENABLED=$(CGO_ENABLED) $(GO) test -timeout $(TIMEOUT)s $(ARGS) -v $(TESTPKGS) | tee test/tests.output
	$(shell test -s test/tests.output && $(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml)

# Mod

go.sum: go.mod ; $(info updating dependencies ...)
	@$(GO) mod tidy -v
	@touch $@

.PHONY: vendor
vendor: go.sum ; $(info retrieving dependencies ...)
	@$(GO) mod vendor -v
	@touch $@

# Dist

.PHONY: licenses
licenses: vendor ; $(info building licenses files ...)
	$(CURDIR)/scripts/go-license-ranger.py > $(CURDIR)/3rdparty-LICENSES.md

3rdparty-LICENSES.md: licenses

.PHONY: dist
dist: 3rdparty-LICENSES.md ; $(info building dist tarball ...)
	@rm -rf "dist/${PACKAGE_NAME}-${VERSION}"
	@mkdir -p "dist/${PACKAGE_NAME}-${VERSION}"
	@mkdir -p "dist/${PACKAGE_NAME}-${VERSION}/scripts"
	@mkdir -p "dist/${PACKAGE_NAME}-${VERSION}/docs"
	@cd dist && \
	cp -avf ../LICENSE.txt "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../README.md "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../3rdparty-LICENSES.md "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../bin/* "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../docs/example-template.ldif "${PACKAGE_NAME}-${VERSION}/docs" && \
	cp -avf ../scripts/libregraph-idmd.binscript "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/libregraph-idmd.service "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/idmd.cfg "${PACKAGE_NAME}-${VERSION}/scripts" && \
	cp -avf ../scripts/*.ldif.in "${PACKAGE_NAME}-${VERSION}/scripts" && \
	tar --owner=0 --group=0 -czvf ${PACKAGE_NAME}-${VERSION}.tar.gz "${PACKAGE_NAME}-${VERSION}" && \
	cd ..

.PHONE: changelog
changelog: ; $(info updating changelog ...)
	$(CHGLOG) --output CHANGELOG.md $(ARGS) v0.1.0..

# Rest

.PHONY: clean
clean: ; $(info cleaning ...)	@
	@rm -rf bin
	@rm -rf test/test.*

.PHONY: version
version:
	@echo $(VERSION)
