SHELL := bash
NAME := idp

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
ci-go-generate: # CI runs ci-node-generate automatically before this target
	@go generate ./...

.PHONY: ci-node-generate
ci-node-generate: assets

.PHONY: assets
assets: yarn-build \
		assets/identifier/static \
		assets/identifier/static/logo.svg \
		assets/identifier/static/favicon.ico

assets/identifier/static:
	mkdir -p assets/identifier/static

LOGO_URL = https://raw.githubusercontent.com/owncloud/assets/main/logo.svg

.PHONY: assets/identifier/static/logo.svg # force overwrite
assets/identifier/static/logo.svg:
	curl -o assets/identifier/static/logo.svg ${LOGO_URL}

FAVICON_URL = https://raw.githubusercontent.com/owncloud/assets/main/favicon.ico

.PHONY: assets/identifier/static/favicon.ico # force overwrite
assets/identifier/static/favicon.ico:
	curl -o assets/identifier/static/favicon.ico ${FAVICON_URL}

.PHONY: yarn-build
yarn-build: node_modules
	#yarn lint #TODO: activate
	#yarn test #TODO: activate
	yarn build

.PHONY: node_modules
node_modules:
	yarn install --immutable
