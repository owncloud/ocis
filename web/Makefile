SHELL := bash
NAME := web

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # supress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
endif

############ go tooling ############
include ../.make/go.mk

############ release #########
include ../.make/release.mk

############ docs generate ############
include ../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: pull-assets # CI runs ci-node-generate automatically before this target
	@go generate $(GENERATE)

.PHONY: ci-node-generate
ci-node-generate:

WEB_ASSETS_VERSION = v3.2.0

.PHONY: pull-assets
pull-assets:
	@rm -Rf assets && mkdir assets/ && tar xzf /root/web/release/web.tar.gz -C assets/
