SHELL := bash
NAME := invitations

include ../../.make/recursion.mk

############ tooling ############
ifneq (, $(shell command -v go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../../.bingo/Variables.mk
GOARCH := $(shell go env GOARCH)
endif

############ go tooling ############
include ../../.make/go.mk

############ release ############
include ../../.make/release.mk

############ docs generate ############
include ../../.make/docs.mk

.PHONY: docs-generate
docs-generate: config-docs-generate

############ dev-docker ############
.PHONY: dev-docker
dev-docker:
	$(MAKE) --no-print-directory release-linux-docker-$(GOARCH)
	docker build -f docker/Dockerfile.linux.$(GOARCH) -t owncloud/ocis-services:$(NAME)-dev .

############ generate ############
include ../../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: $(MOCKERY) # CI runs ci-node-generate automatically before this target
	$(MOCKERY) --output pkg/mocks --case underscore --name Client  --srcpkg "github.com/owncloud/ocis/v2/ocis-pkg/keycloak"


.PHONY: ci-node-generate
ci-node-generate:

############ licenses ############
.PHONY: ci-node-check-licenses
ci-node-check-licenses:

.PHONY: ci-node-save-licenses
ci-node-save-licenses:
