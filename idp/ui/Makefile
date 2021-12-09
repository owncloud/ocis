# Tools

YARN ?= yarn

# Variables

VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null | sed 's/^v//' || \
			cat $(CURDIR)/../.version 2> /dev/null || echo 0.0.0-unreleased)

# Build

.PHONY: all
all: build

.PHONY: build
build: vendor | src i18n ; $(info building identifier Webapp ...)	@
	@rm -rf build

	REACT_APP_KOPANO_BUILD="${VERSION}" $(YARN) run build

.PHONY: src
src:
	@$(MAKE) --no-print-directory -C src

.PHONY: i18n
i18n: vendor
	@$(MAKE) --no-print-directory -C i18n

.PHONY: lint
lint: vendor ; $(info running eslint ...)	@
	@$(YARN) eslint . --cache && echo "eslint: no lint errors"

.PHONY: lint-checkstyle
lint-checkstyle: vendor ; $(info running eslint checkstyle ...)	@
	@mkdir -p ../test
	$(YARN) eslint -f checkstyle -o ../test/tests.eslint.xml . || true

# Yarn

.PHONY: vendor
vendor: .yarninstall

.yarninstall: package.json ; $(info getting dependencies with yarn ...)   @
	@$(YARN) install --silent
	@touch $@

# Stuff

.PHONY: licenses
licenses:
	echo "## Konnect identifier webapp\n"
	@$(YARN) run  -s licenses

.PHONY: clean ; $(info cleaning identifier Webapp ...)	@
clean:
	$(YARN) cache clean
	@rm -rf build
	@rm -rf node_modules
	@rm -f .yarninstall

	@$(MAKE) --no-print-directory -C src clean

.PHONY: version
version:
	@echo $(VERSION)
