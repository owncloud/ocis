# Tools

YARN ?= yarn
MSGCAT ?= msgcat
MSGMERGE ?= msgmerge
MSGFMT ?= msgfmt

# Variables

POT = translation.pot
POS ?= $(wildcard *.po)

# Translations

.PHONY: build
build: json

.PYTHON: json
json:
	@for po in $(POS); do \
		lang=$$(echo $$po | sed "s/\.po//"); \
		$(YARN) i18next-conv -K --skipUntranslated -l $$lang -s ../i18n/$$po -t ../src/locales/$$lang/translation.json; \
	done
	$(YARN) node ../i18n/build-json.js ../src/locales/locales.json $(POS)

dev-translation.json: FORCE
	$(YARN) i18next --fail-on-warnings

.PHONY: extract
extract: pot

.PHONY: pot
pot: dev-translation.json
	@tmpfile1=$(shell mktemp).po; \
	tmpfile2=$(shell mktemp).po; \
	trap 'rm -f "$$tmpfile1" "$$tmpfile2"' EXIT; \
	$(YARN) i18next-conv --project "LibreGraph Connect Identifier" -K -l en -s i18n/dev-translation.json -t $$tmpfile1; \
	$(YARN) node i18n/build-pot.js $$tmpfile1 $$tmpfile2; \
	$(MSGCAT) --no-wrap -o $(POT) $$tmpfile2

.PHONY: merge
merge: $(POS)

$(POS): FORCE $(POT)
	@echo -n "$@ " && \
	$(MSGMERGE) -U \
		--backup=none \
		--no-wrap \
		--sort-output \
		$@ $(POT)

.PHONY: stats
stats:
	$(foreach po, $(POS), $(shell $(MSGFMT) -v --statistics $(po)))
	@- true

FORCE:
