SHELL := bash
NAME := ocis

include ../.make/recursion.mk

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
GOARCH := $(shell go env GOARCH)
endif

############ go tooling ############
include ../.make/go.mk

############ release ############
include ../.make/release.mk

############ docs generate ############
CONFIG_DOCS_BASE_PATH := ../docs

include ../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate

############ dev-docker ############
.PHONY: dev-docker
dev-docker:
	$(MAKE) --no-print-directory release-linux-docker-$(GOARCH)
	docker build -f docker/Dockerfile.linux.$(GOARCH) -t owncloud/ocis:dev .

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target

.PHONY: ci-node-generate
ci-node-generate: # ocis needs assets of all other modules
	@if [ $(MAKE_DEPTH) -le 1 ]; then \
	$(MAKE) --no-print-directory -C .. ci-node-generate \
	; fi;
