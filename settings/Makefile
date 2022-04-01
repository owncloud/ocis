SHELL := bash
NAME := settings

include ../.make/recursion.mk

.PHONY: test-acceptance-webui
test-acceptance-webui:
	./ui/tests/run-acceptance-test.sh $(FEATURE_PATH)


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

############ l10n ############
include ../.make/l10n.mk

.PHONY: docs-generate
docs-generate: config-docs-generate \
			   grpc-docs-generate

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: protobuf # CI runs ci-node-generate automatically before this target

.PHONY: ci-node-generate
ci-node-generate: yarn-build

.PHONY: yarn-build
yarn-build: node_modules
	yarn lint
	yarn test
	yarn build

.PHONY: node_modules
node_modules:
	@yarn install --immutable 2>&1 >/dev/null

############ protobuf ############
include ../.make/protobuf.mk

.PHONY: protobuf
protobuf: buf-generate

############ licenses ############
.PHONY: ci-node-check-licenses
ci-node-check-licenses: node_modules
	yarn licenses:check

.PHONY: ci-node-save-licenses
ci-node-save-licenses: node_modules
	yarn licenses:csv
	yarn licenses:save
