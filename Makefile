SHELL := bash
NAME := ocis
IMPORT := github.com/owncloud/$(NAME)
HUGO := hugo

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean:
	rm -rf $(HUGO)

.PHONY: docs-copy
docs-copy:
	mkdir -p $(HUGO); \
	mkdir -p $(HUGO)/content/; \
	cd $(HUGO); \
	git init; \
	git remote rm origin; \
	git remote add origin https://github.com/owncloud/owncloud.github.io; \
	git fetch --depth=1; \
	git checkout origin/source -f; \
	rsync -ax --delete ../docs/ content/; \

.PHONY: config-docs-generate
config-docs-generate:
	go run github.com/owncloud/flaex >| docs/configuration.md

.PHONY: docs-build
docs-build:
	cd $(HUGO); hugo

.PHONY: docs
docs: docs-copy docs-build

BEHAT_BIN=vendor-bin/behat/vendor/bin/behat

.PHONY: test-acceptance-api
test-acceptance-api: vendor-bin/behat/vendor
	BEHAT_BIN=$(BEHAT_BIN) $(PATH_TO_CORE)/tests/acceptance/run.sh --remote --type api

vendor/bamarni/composer-bin-plugin: composer.lock
	composer install

vendor-bin/behat/vendor: vendor/bamarni/composer-bin-plugin vendor-bin/behat/composer.lock
	composer bin behat install --no-progress

vendor-bin/behat/composer.lock: vendor-bin/behat/composer.json
	@echo behat composer.lock is not up to date.

composer.lock: composer.json
	@echo composer.lock is not up to date.
