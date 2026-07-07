SHELL := bash
NAME := web
WEB_DIR := ../../web

include ../../.make/recursion.mk

############ tooling ############
ifneq (, $(shell command -v go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../../.bingo/Variables.mk
endif

############ go tooling ############
include ../../.make/go.mk

############ release #########
include ../../.make/release.mk

############ docs generate ############
include ../../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate

############ generate ############
include ../../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target

.PHONY: ci-node-generate
ci-node-generate: build-assets

.PHONY: build-assets
build-assets:
	cd $(WEB_DIR) && pnpm install && pnpm check:types && pnpm build
	cp $(WEB_DIR)/config/config.json.dist $(WEB_DIR)/dist/
	git clean -xfd assets
	cp -r $(WEB_DIR)/dist/. assets/core/

############ licenses ############
.PHONY: ci-node-check-licenses
ci-node-check-licenses:
	cd $(WEB_DIR) && pnpm licenses:check

.PHONY: ci-node-save-licenses
ci-node-save-licenses:
	@mkdir -p ../../third-party-licenses/node/web
	cd $(WEB_DIR) && pnpm licenses:csv && pnpm licenses:save
	cp -r $(WEB_DIR)/third-party-licenses/third-party-licenses/. ../../third-party-licenses/node/web/
