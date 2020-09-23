SHELL := bash
NAME := ocis
IMPORT := github.com/owncloud/$(NAME)
HUGO := hugo
EXTENSIONS := accounts glauth graph konnectd ocis-phoenix ocis-reva ocs proxy settings store thumbnails webdav

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean:
	rm -rf $(HUGO)

.PHONY: generate-docs $(EXTENSIONS)
generate-docs: $(EXTENSIONS)
$(EXTENSIONS):
	$(MAKE) -C $@ docs; \\
	mkdir -p docs/extensions/$@; \\
	cp -R $@/docs/ docs/extensions/$@

.PHONY: clean-docs
clean-docs:
	rm -rf docs

.PHONY: ocis-docs
ocis-docs:
	mkdir -p docs/ocis; \\
	$(MAKE) -C ocis docs; \\
	cp -R ocis/docs/ docs/ocis

.PHONY: docs
docs: clean-docs generate-docs ocis-docs

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
