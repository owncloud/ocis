SHELL := bash
NAME := ocis-pkg

include ../.make/recursion.mk

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
endif

############ go tooling ############
include ../.make/go.mk

############ release ############
include ../.make/release.mk

############ docs generate ############
SKIP_CONFIG_DOCS_GENERATE = 1

include ../.make/docs.mk

.PHONY: docs-generate
docs-generate:

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target

.PHONY: ci-node-generate
ci-node-generate:
