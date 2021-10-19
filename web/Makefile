SHELL := bash
NAME := web

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
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
ci-go-generate: # CI runs ci-node-generate automatically before this target

.PHONY: ci-node-generate
ci-node-generate: pull-assets

WEB_ASSETS_VERSION = v4.3.0

.PHONY: pull-assets
pull-assets:
	@rm -Rf assets && mkdir assets/ && curl -slL -o- https://github.com/owncloud/web/releases/download/$(WEB_ASSETS_VERSION)/web.tar.gz | tar xzf - -C assets/
