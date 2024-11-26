SHELL := bash
NAME := web
WEB_ASSETS_VERSION = v11.0.5

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
ci-node-generate: pull-assets

.PHONY: pull-assets
pull-assets:
	git clean -xfd assets
	curl --fail -slL -o- https://github.com/owncloud/web/releases/download/$(WEB_ASSETS_VERSION)/web.tar.gz | tar xzf - -C assets/core/

############ licenses ############
.PHONY: ci-node-check-licenses
ci-node-check-licenses: # ownCloud Web does the check on it's own

.PHONY: ci-node-save-licenses
ci-node-save-licenses:
	@mkdir -p ../../third-party-licenses/node/web
	@curl --fail -slL -o- https://github.com/owncloud/web/releases/download/$(WEB_ASSETS_VERSION)/third-party-licenses.tar.gz | tar xzf - -C ../../third-party-licenses/node/web
