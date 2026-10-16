SHELL := bash
DIST := ${CURDIR}/dist
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
	@echo -e "${GREEN}Any changes for the web developer documentation have to be made at:${RESET}\n"
	@echo -e "${PURPLE}\thttps://github.com/owncloud/docs.owncloud.com/tree/main/content/ocis/${RESET}\n"
	@echo

.PHONY: list
list:
	@echo -e 'Available .PHONY targets: \n'
	@grep -P -o '(?<=^\.PHONY: )(.*)' Makefile | sort -u
	@echo -e ''

.PHONY: clean
clean:
	rm -rf ${DIST} ${RELEASE} ${NODE_MODULES}

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
