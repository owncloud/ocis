SHELL := bash
NAME := ocis

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

.PHONY: config-docs-generate
config-docs-generate: $(FLAEX)
	@echo "$(NAME): generating config docs"
	@$(FLAEX) >| ../docs/$(NAME)/configuration.md

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target
	@go generate ./...

.PHONY: ci-node-generate
ci-node-generate:
	@npm install --global yarn
	@make -C ../accounts ci-node-generate
	@make -C ../graph-explorer ci-node-generate
	@make -C ../settings ci-node-generate
	@make -C ../web ci-node-generate
