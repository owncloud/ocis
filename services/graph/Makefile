SHELL := bash
NAME := graph

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
ci-go-generate: $(MOCKERY) # CI runs ci-node-generate automatically before this target
	$(MOCKERY) --dir pkg/service/v0 --case underscore --name HTTPClient
	$(MOCKERY) --dir pkg/service/v0 --case underscore --name Publisher
	$(MOCKERY) --dir pkg/service/v0 --case underscore --name Permissions
	$(MOCKERY) --dir pkg/service/v0 --case underscore --name RoleService
	$(MOCKERY) --dir pkg/identity --output pkg/identity/mocks --case underscore --name Backend
	$(MOCKERY) --dir pkg/identity --output pkg/identity/mocks --case underscore --name EducationBackend
	$(MOCKERY) --srcpkg github.com/go-ldap/ldap/v3 --case underscore --filename ldapclient.go --name Client


.PHONY: ci-node-generate
ci-node-generate:

############ licenses ############
.PHONY: ci-node-check-licenses
ci-node-check-licenses:

.PHONY: ci-node-save-licenses
ci-node-save-licenses:
