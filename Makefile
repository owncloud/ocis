SHELL := bash

.PHONY: help
help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo
	@echo -e "Testing:\n"
	@echo -e "make test-acceptance-api\trun API acceptance tests"
	@echo -e "make clean-tests\t\tdelete API tests framework dependencies"
	@echo
	@echo -e "See the Makefile in the ocis folder for other build and test targets"

.PHONY: clean-tests
clean-tests:
	rm -Rf vendor-bin/**/vendor vendor-bin/**/composer.lock

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
