SHELL := bash
NAME := idp

include ../../.make/recursion.mk

############ tooling ############
ifneq (, $(shell command -v go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../../.bingo/Variables.mk
endif

############ go tooling ############
include ../../.make/go.mk

############ release ############
include ../../.make/release.mk

############ docs generate ############
include ../../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate

############ generate ############
include ../../.make/generate.mk

FAVICON_URL = https://raw.githubusercontent.com/owncloud/assets/main/favicon.ico

assets/identifier/static:
	mkdir -p assets/identifier/static

.PHONY: assets/identifier/static/favicon.ico # force overwrite
assets/identifier/static/favicon.ico:
	curl --fail -o assets/identifier/static/favicon.ico ${FAVICON_URL}

.PHONY: ci-go-generate
ci-go-generate: # CI runs ci-node-generate automatically before this target

.PHONY: assets
assets:
	@mkdir -p assets/identifier/static

############ testing ############
.PHONY: test-go
test-go: ## Run all tests (unit + integration)
	@go test -v -tags=integration ./pkg/service/v0/... -timeout 60s
