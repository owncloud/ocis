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
	yarn install --immutable

############ protobuf ############
include ../.make/protobuf.mk

.PHONY: protobuf
protobuf: buf-generate

############ licenses ############
.PHONY: ci-node-licenses-check
ci-node-licenses-check:
	yarn licenses:check

.PHONY: ci-node-licenses-save
ci-node-licenses-save:
	yarn licenses:csv
	yarn licenses:save
