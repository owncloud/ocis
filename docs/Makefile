SHELL := bash

include ../.bingo/Variables.mk

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docs-generate
docs-generate: ## run docs-generate for all oCIS extensions
	@$(MAKE) --no-print-directory -C ../ docs-generate

.PHONY: docs-init
docs-init:
	@mkdir -p hugo/content/
	@cd hugo && git init
	@cd hugo && git config advice.detachedHead false
	@cd hugo && git remote rm origin || true
	@cd hugo && git remote add origin https://github.com/owncloud/owncloud.github.io
	@cd hugo && git fetch --depth=1
	@cd hugo && git checkout origin/source -f
	@$(MAKE) -C hugo theme

.PHONY: docs-serve
docs-serve: docs-init docs-generate docs-copy ## serve docs with hugo
	@bash -c "trap 'trap - SIGINT SIGTERM ERR; $(MAKE) --no-print-directory docs-sync-stop; exit 0' SIGINT SIGTERM ERR; $(MAKE) --no-print-directory docs-sync-start && $(MAKE) --no-print-directory hugo-serve"

.PHONY: test
test: $(HUGO)
	@cd hugo && $(HUGO)

.PHONY: hugo-serve
hugo-serve: $(HUGO)
	@cd hugo && $(HUGO) server

.PHONY: docs-copy
docs-copy: docs-init docs-sync-start docs-sync-stop

.PHONY: docs-sync-start
docs-sync-start: $(MUTAGEN)
	@$(MUTAGEN) project terminate || true
	@$(MUTAGEN) daemon stop || true
	@$(MUTAGEN) project start
	@$(MUTAGEN) project flush

.PHONY: docs-sync-stop
docs-sync-stop: $(MUTAGEN)
	@$(MUTAGEN) project terminate
	@$(MUTAGEN) daemon stop

.PHONY: clean
clean: ## clean up docs build artifacts
	@rm -rf hugo
