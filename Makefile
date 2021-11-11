SHELL := bash

# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)

L10N_MODULES := $(shell find . -path '*.tx*' -name 'config' | sed 's|/[^/]*$$||' | sed 's|/[^/]*$$||' | sed 's|/[^/]*$$||')

# if you add a module here please also add it to the .drone.star file
OCIS_MODULES = \
	accounts \
	glauth \
	graph \
	graph-explorer \
	idp \
	ocis \
	ocis-pkg \
	ocs \
	proxy \
	settings \
	storage \
	store \
	thumbnails \
	web \
	webdav

# bin file definitions
PHP_CS_FIXER=php -d zend.enable_gc=0 vendor-bin/owncloud-codestyle/vendor/bin/php-cs-fixer
PHP_CODESNIFFER=vendor-bin/php_codesniffer/vendor/bin/phpcs
PHP_CODEBEAUTIFIER=vendor-bin/php_codesniffer/vendor/bin/phpcbf
PHAN=php -d zend.enable_gc=0 vendor-bin/phan/vendor/bin/phan
PHPSTAN=php -d zend.enable_gc=0 vendor-bin/phpstan/vendor/bin/phpstan

ifneq (, $(shell which go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include .bingo/Variables.mk
endif

.PHONY: help
help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo
	@echo -e "${GREEN}Testing with test suite natively installed:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.dev/ocis/development/testing/#testing-with-test-suite-natively-installed${RESET}\n"
	@echo -e "\tmake test-acceptance-api\t${BLUE}run API acceptance tests${RESET}"
	@echo -e "\tmake clean-tests\t\t${BLUE}delete API tests framework dependencies${RESET}"
	@echo
	@echo -e "${BLACK}---------------------------------------------------------${RESET}"
	@echo
	@echo -e "${RED}You also should have a look at other available Makefiles:${RESET}"
	@echo
	@echo -e "${GREEN}oCIS:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.dev/ocis/development/build/${RESET}\n"
	@echo -e "\tsee ./ocis/Makefile"
	@echo -e "\tor run ${YELLOW}make -C ocis help${RESET}"
	@echo
	@echo -e "${GREEN}Documentation:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.dev/ocis/development/build-docs/${RESET}\n"
	@echo -e "\tsee ./docs/Makefile"
	@echo -e "\tor run ${YELLOW}make -C docs help${RESET}"
	@echo
	@echo -e "${GREEN}Testing with test suite in docker:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.dev/ocis/development/testing/#testing-with-test-suite-in-docker${RESET}\n"
	@echo -e "\tsee ./tests/acceptance/docker/Makefile"
	@echo -e "\tor run ${YELLOW}make -C tests/acceptance/docker help${RESET}"
	@echo
	@echo -e "${GREEN}Tools for developing tests:\n${RESET}"
	@echo -e "\tmake test-php-style\t\t${BLUE}run PHP code style checks${RESET}"
	@echo -e "\tmake test-php-style-fix\t\t${BLUE}run PHP code style checks and fix any issues found${RESET}"
	@echo

.PHONY: clean-tests
clean-tests:
	@rm -Rf vendor-bin/**/vendor vendor-bin/**/composer.lock tests/acceptance/output

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

.PHONY: generate
generate:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod generate; \
    done

.PHONY: vet
vet:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod vet; \
    done

.PHONY: clean
clean:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod clean; \
    done

.PHONY: docs-generate
docs-generate:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod docs-generate; \
    done

.PHONY: ci-go-generate
ci-go-generate:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod ci-go-generate; \
    done

.PHONY: ci-node-generate
ci-node-generate:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod ci-node-generate; \
    done

.PHONY: go-mod-tidy
go-mod-tidy:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod go-mod-tidy; \
    done

.PHONY: test
test:
	@for mod in $(OCIS_MODULES); do \
        $(MAKE) --no-print-directory -C $$mod test; \
    done

.PHONY: go-coverage
go-coverage:
	@if [ ! -f coverage.out ]; then $(MAKE) test  &>/dev/null; fi;
	@for mod in $(OCIS_MODULES); do \
        echo -n "% coverage $$mod: "; $(MAKE) --no-print-directory -C $$mod go-coverage; \
    done

.PHONY: protobuf
protobuf:
	@for mod in $(OCIS_MODULES); do \
        echo -n "% protobuf $$mod: "; $(MAKE) --no-print-directory -C $$mod protobuf; \
    done

.PHONY: bingo-update
bingo-update: $(BINGO)
	$(BINGO) get -l -u

CHANGELOG_VERSION =

.PHONY: changelog
changelog: $(CALENS)
ifndef CHANGELOG_VERSION
	$(error CHANGELOG_VERSION is undefined)
endif
	$(CALENS) --version $(CHANGELOG_VERSION) -o ocis/dist/CHANGELOG.md

.PHONY: l10n-push
l10n-push:
	@for extension in $(L10N_MODULES); do \
    	make -C $$extension l10n-push; \
	done

.PHONY: l10n-pull
l10n-pull:
	@for extension in $(L10N_MODULES); do \
    	make -C $$extension l10n-pull; \
	done

.PHONY: l10n-clean
l10n-clean:
	@for extension in $(L10N_MODULES); do \
		make -C $$extension l10n-clean; \
	done

.PHONY: l10n-read
l10n-read:
	@for extension in $(L10N_MODULES); do \
    	make -C $$extension l10n-read; \
    done

.PHONY: l10n-write
l10n-write:
	@for extension in $(L10N_MODULES); do \
    	make -C $$extension l10n-write; \
    done

.PHONY: ci-format
ci-format: $(BUILDIFIER)
	$(BUILDIFIER) --mode=fix .drone.star

.PHONY: test-php-style
test-php-style: vendor-bin/owncloud-codestyle/vendor vendor-bin/php_codesniffer/vendor
	$(PHP_CS_FIXER) fix -v --diff --allow-risky yes --dry-run
	$(PHP_CODESNIFFER) --cache --runtime-set ignore_warnings_on_exit --standard=phpcs.xml tests/acceptance

.PHONY: test-php-style-fix
test-php-style-fix: vendor-bin/owncloud-codestyle/vendor
	$(PHP_CS_FIXER) fix -v --diff --allow-risky yes
	$(PHP_CODEBEAUTIFIER) --cache --runtime-set ignore_warnings_on_exit --standard=phpcs.xml tests/acceptance

vendor-bin/owncloud-codestyle/vendor: vendor/bamarni/composer-bin-plugin vendor-bin/owncloud-codestyle/composer.lock
	composer bin owncloud-codestyle install --no-progress

vendor-bin/owncloud-codestyle/composer.lock: vendor-bin/owncloud-codestyle/composer.json
	@echo owncloud-codestyle composer.lock is not up to date.

vendor-bin/php_codesniffer/vendor: vendor/bamarni/composer-bin-plugin vendor-bin/php_codesniffer/composer.lock
	composer bin php_codesniffer install --no-progress

vendor-bin/php_codesniffer/composer.lock: vendor-bin/php_codesniffer/composer.json
	@echo php_codesniffer composer.lock is not up to date.
