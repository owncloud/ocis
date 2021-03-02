SHELL := bash
NAME := ocis-pkg

############ tooling ############
include ../.bingo/Variables.mk

############ go tooling ############
include ../.make/go.mk

############ release ############
include ../.make/release.mk

############ docs generate ############
include ../.make/docs.mk

.PHONY: docs-generate
docs-generate:

.PHONY: config-docs-generate
config-docs-generate:
	@echo "ocis-pkg: nothing to do"

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target
	@go generate $(GENERATE)

.PHONY: ci-node-generate
ci-node-generate:
