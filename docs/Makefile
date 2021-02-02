SHELL := bash

OCIS_MODULES = \
	accounts \
	glauth \
	idp \
	ocis \
	ocs \
	ocis-pkg \
	proxy \
	settings \
	storage \
	store \
	thumbnails \
	web \
	webdav

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: config-docs-generate
config-docs-generate: ## run config-docs-generate for all oCIS extensions
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C ../$$mod config-docs-generate; \
    done

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
docs-serve: config-docs-generate docs-copy ## serve docs with hugo
	@docker run --rm --network host -v $(shell pwd)/hugo:/src owncloudci/hugo:0 server

.PHONY: clean
clean: ## clean up docs build artifacts
	@rm -rf hugo
