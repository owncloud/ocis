SHELL := bash
NAME := graph

############ tooling ############
ifneq (, $(shell which go 2> /dev/null)) # supress `command not found warnings` for non go targets in CI
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
	@go generate $(GENERATE)

.PHONY: ci-node-generate
ci-node-generate:

############ openapi ############
OPENAPI_VERSION := v0
OPENAPI_SRC := pkg/openapi/$(OPENAPI_VERSION)
OPENGRAPH_VERSION := v0.0

include ../.make/openapi.mk

$(OPENAPI_SRC):
	@mkdir -p $(OPENAPI_SRC)

.PHONY: openapi
openapi: $(OPENAPI_SRC) \
		 $(OPENAPI_SRC)/${NAME}-${OPENGRAPH_VERSION}.yml \
		 $(OPENAPI_SRC)/${NAME}.types.go \
		 $(OPENAPI_SRC)/${NAME}.server.go

		  
SPEC_URL = https://raw.githubusercontent.com/owncloud/open-graph-api/fix-code-generation/api/openapi-spec/${OPENGRAPH_VERSION}.yml

#.PHONY: $(OPENAPI_SRC)/${NAME}-${OPENGRAPH_VERSION}.yml # force overwrite
$(OPENAPI_SRC)/${NAME}-${OPENGRAPH_VERSION}.yml:
	curl -o $(OPENAPI_SRC)/${NAME}-${OPENGRAPH_VERSION}.yml ${SPEC_URL}