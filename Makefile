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
	onlyoffice \
	proxy \
	settings \
	storage \
	store \
	thumbnails \
	web \
	webdav

ifneq (, $(shell which go 2> /dev/null)) # supress `command not found warnings` for non go targets in CI
include .bingo/Variables.mk
endif

.PHONY: help
help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo
	@echo -e "${GREEN}Testing with test suite natively installed:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.github.io/ocis/development/testing/#testing-with-test-suite-natively-installed${RESET}\n"
	@echo -e "\tmake test-acceptance-api\t${BLUE}run API acceptance tests${RESET}"
	@echo -e "\tmake clean-tests\t\t${BLUE}delete API tests framework dependencies${RESET}"
	@echo
	@echo -e "${BLACK}---------------------------------------------------------${RESET}"
	@echo
	@echo -e "${RED}You also should have a look at other available Makefiles:${RESET}"
	@echo
	@echo -e "${GREEN}oCIS:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.github.io/ocis/development/building/${RESET}\n"
	@echo -e "\tsee ./ocis/Makefile"
	@echo -e "\tor run ${YELLOW}make -C ocis help${RESET}"
	@echo
	@echo -e "${GREEN}Documentation:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.github.io/ocis/development/building-docs/${RESET}\n"
	@echo -e "\tsee ./docs/Makefile"
	@echo -e "\tor run ${YELLOW}make -C docs help${RESET}"
	@echo
	@echo -e "${GREEN}Testing with test suite in docker:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.github.io/ocis/development/testing/#testing-with-test-suite-in-docker${RESET}\n"
	@echo -e "\tsee ./tests/acceptance/docker/Makefile"
	@echo -e "\tor run ${YELLOW}make -C tests/acceptance/docker help${RESET}"
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
	$(BINGO) get -u

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
