SHELL := bash

include ../.bingo/Variables.mk

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docs-generate
docs-generate: ## run docs-generate for all oCIS extensions
	@$(MAKE) --no-print-directory -C ../ docs-generate

.PHONY: config-docs-generate
config-docs-generate: ## run config-docs-generate for all oCIS extensions
	@$(MAKE) --no-print-directory -C ../ config-docs-generate

.PHONY: docs-copy
docs-copy: ## copy docs to hugo
	@mkdir -p hugo/content/
	@cd hugo; \
	git init; \
	git config advice.detachedHead false; \
	git remote rm origin; \
	git remote add origin https://github.com/owncloud/owncloud.github.io; \
	git fetch --depth=1; \
	git checkout origin/source -f; \
	rsync -ax --delete --exclude hugo/ --exclude Makefile --exclude .gitignore --exclude README.md ../. content/; \

.PHONY: docs-serve
docs-serve: $(HUGO) docs-generate docs-copy ## serve docs with hugo
	cd hugo && $(HUGO) server

.PHONY: docs-serve-simple
docs-serve-simple: $(HUGO) config-docs-generate docs-copy ## serve docs with hugo without protobuf generation
	cd hugo && $(HUGO) server

.PHONY: clean
clean: ## clean up docs build artifacts
	@rm -rf hugo
