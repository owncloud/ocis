SHELL := bash
DIST := ${CURDIR}/dist
HUGO := docs/hugo
RELEASE := ${CURDIR}/release
NODE_MODULES := ${CURDIR}/node_modules

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

node_modules: package.json pnpm-lock.yaml
	[ -n "${NO_INSTALL}" ] || pnpm install
	touch ${NODE_MODULES}

.PHONY: help
DEFAULT_GOAL := help
help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo
	@echo -e "${GREEN}List all available .PHONY targets:${RESET}\n"
	@echo -e "\tmake list\t\t${BLUE}sorted alphabetically${RESET}"
	@echo -e "${BLACK}---------------------------------------------------------${RESET}"
	@echo -e "${GREEN}Install/update required packages:${RESET}\n"
	@echo -e "\tmake\t\t\t${BLUE}run make without target${RESET}"
	@echo -e "${BLACK}---------------------------------------------------------${RESET}"
	@echo
	@echo -e "${GREEN}Documentation:${RESET}\n"
	@echo -e "${PURPLE}\tdocs: https://owncloud.dev/ocis/build-docs/${RESET}\n"
	@echo -e "\trun ${YELLOW}make list | grep docs-\t\t${BLUE}note: run all docs command via this makefile${RESET}"
	@echo

.PHONY: list
list:
	@echo -e 'Available .PHONY targets: \n'
	@grep -P -o '(?<=^\.PHONY: )(.*)' Makefile | sort -u
	@echo -e ''

.PHONY: clean
clean:
	rm -rf ${DIST} ${HUGO} ${RELEASE} ${NODE_MODULES}

.PHONY: release
release: clean
	make -f Makefile.release

#
# Release
# make this app compatible with the ownCloud
# default build tools
#
.PHONY: dist
dist:
	make -f Makefile.release

# note that everything docs related is located in the docs/ folder
# we keep this original calls for the sake of history and ease of use
# for drone only, prepare docs, do not run manually
.PHONY: docs-generate          # 1. prepare docs
docs-generate:
# initialize the docs build environment
	@$(MAKE) --no-print-directory -C docs docs-init

# remind one that web needs prerequisites installed
	@$(MAKE) --no-print-directory -C docs docs-first-time-message

# copy required resources into hugo/content
.PHONY: docs-copy              # 2. copy required doc resources
docs-copy:
	@$(MAKE) --no-print-directory -C docs docs-copy

# the docs-build|serve commands requires that docs-init was run first for the required data to exists
# create a docs build
.PHONY: docs-build             # 3. build prepared docs
docs-build:
	@$(MAKE) --no-print-directory -C docs docs-build

# serve built docs with hugo
.PHONY: docs-serve             # serve the docs build
docs-serve:
	@$(MAKE) --no-print-directory -C docs docs-serve

# clean up doc build artifacts 
.PHONY: docs-clean             # clean all docs artifacts, must be run as sudo
docs-clean:
	@$(MAKE) --no-print-directory -C docs docs-clean

# imitate a full drone run locally to build docs without pushing to the web.
# this can help identify uncaught issues when running `make docs-serve` only.
.PHONY: docs-local             # run all steps as drone would do it (1, 2, 3)
docs-local:
	@$(MAKE) --no-print-directory docs-generate
	@$(MAKE) --no-print-directory docs-copy
	@$(MAKE) --no-print-directory docs-build 

# prepare a link from the root to the hugo folder because the image requires it
# note that on local building, the referenced container of inside the hugo/makefile is used
.PHONY: docs-hugo-drone-prep   # only used for drone !
docs-hugo-drone-prep:
	@$(MAKE) --no-print-directory -C docs docs-hugo-drone-prep

# translation relevant
.PHONY: l10n-push
l10n-push:
	@$(MAKE) --no-print-directory -C packages/web-runtime/l10n push

.PHONY: l10n-pull
l10n-pull:
	@$(MAKE) --no-print-directory -C packages/web-runtime/l10n pull

.PHONY: l10n-clean
l10n-clean:
	@$(MAKE) --no-print-directory -C packages/web-runtime/l10n clean

.PHONY: l10n-read
l10n-read: node_modules
	@$(MAKE) --no-print-directory -C packages/web-runtime/l10n extract

.PHONY: l10n-write
l10n-write: node_modules
	@$(MAKE) --no-print-directory -C packages/web-runtime/l10n translations

.PHONY: generate-qa-activity-report
generate-qa-activity-report: node_modules
	@if [ -z "${MONTH}" ] || [ -z "${YEAR}" ]; then \
		echo "Please set the MONTH and YEAR environment variables. Usage: make generate-qa-activity-report MONTH=<month> YEAR=<year>"; \
		exit 1; \
	fi
	pnpm exec node generate-qa-activity-report.js --month ${MONTH} --year ${YEAR}
