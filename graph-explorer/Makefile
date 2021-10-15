SHELL := bash
NAME := graph-explorer

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
endif

############ go tooling ############
include ../.make/go.mk

############ release ############
include ../.make/release.mk

############ docs generate ############
include ../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: pull-assets # CI runs ci-node-generate automatically before this target
	#@go generate ./...

.PHONY: ci-node-generate
ci-node-generate:

.PHONY: pull-assets
pull-assets:
	@rm -Rf assets && mkdir assets/ && curl -slL -o- https://github.com/owncloud/graph-explorer/releases/download/v1.7.10.3/release.tar.gz | tar xvzf - -C assets/

