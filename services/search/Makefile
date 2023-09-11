SHELL := bash
NAME := search

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

.PHONY: ci-go-generate
ci-go-generate: $(PIGEON) $(MOCKERY) # CI runs ci-node-generate automatically before this target
	$(MOCKERY) --dir pkg/engine --output pkg/engine/mocks --case underscore --name Engine
	$(MOCKERY) --dir pkg/content --output pkg/content/mocks --case underscore --name Extractor
	$(MOCKERY) --dir pkg/content --output pkg/content/mocks --case underscore --name Retriever
	$(MOCKERY) --dir pkg/search --output pkg/search/mocks --case underscore --name Searcher
	$(PIGEON) -optimize-grammar -optimize-parser -o pkg/query/kql/dictionary_gen.go pkg/query/kql/dictionary.peg

.PHONY: ci-node-generate
ci-node-generate:

############ licenses ############
.PHONY: ci-node-check-licenses
ci-node-check-licenses:

.PHONY: ci-node-save-licenses
ci-node-save-licenses:
