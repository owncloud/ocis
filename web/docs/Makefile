SHELL := bash
NAME := web
# we need an override so that this variable is cleanly set in this makefile
# else it might come from somewhere else
override HUGO := ${CURDIR}/hugo
override HDIR := ${CURDIR}/hugo/.git

# print a hint when doing a local build because it may fail if you have not run pnpm commands at least once
docs-first-time-message:
	@echo
	@echo 'Note, on a fresh web clone, you need to run `pnpm install` and `pnpm build` once to create required directories.'

# if the directory exists, print that we do not need to prepare
hugo-exists:
	@if [ -d $(HDIR) ]; then \
		echo; \
		echo 'Build environment for hugo exists, nothing to setup, continuing.'; \
	fi

# if the directory does not extist, initialize the framework and create the theme
$(HDIR):
	@echo 'Initialize the hugo framework and build the theme, no docs get built.'
	@mkdir -p $(HUGO)/content/
	@mkdir -p $(HUGO)/extensions/
	@mkdir -p $(HUGO)/public/
	@cd $(HUGO) && git init
	@cd $(HUGO) && git config advice.detachedHead false
	@cd $(HUGO) && git remote rm origin || true
	@cd $(HUGO) && git remote add origin https://github.com/owncloud/owncloud.github.io
	@cd $(HUGO) && git fetch --depth=1
	@cd $(HUGO) && git checkout origin/main -f
	@$(MAKE) -C $(HUGO) --no-print-directory theme

# print available make targets
.PHONY: help
.DEFAULT_GOAL := help
help:
	@echo -e 'Available .PHONY targets: \n'
	@grep -P -o '(?<=^\.PHONY: )(.*)' Makefile | sort -u
	@echo

# alternative target because we also use it in the root's makefile
.PHONY:list
list: help

# for drone only. drone uses an image for hugo that needs the hugo folder in the root
# CURDIR will work for make only but not when running manually in the shell
.PHONY: docs-hugo-drone-prep
docs-hugo-drone-prep:
	@echo
	@echo 'Linking directory hugo for drone.'
	@rm -f ../hugo
	@ln -s $(HUGO) ../hugo

# create the docs build environment 
.PHONY: docs-init
docs-init: hugo-exists $(HDIR)
# if there is nothing to do, stop printing that
# another way would be https://stackoverflow.com/questions/58039810/makefiles-what-is-an-order-only-prerequisite
	@:

# the docs-build|serve commands require that docs-init was run first for the required data to exists
# create a docs build
.PHONY: docs-build
docs-build:
	@echo
	@echo 'Building web for owncloud.dev'
	@$(MAKE) -C $(HUGO) --no-print-directory hugo-build
	@echo 'To view the rendered docs in the browser, run: make docs-serve'

# serve built docs with hugo
.PHONY: docs-serve
docs-serve:
	@$(MAKE) -C $(HUGO) --no-print-directory hugo-server

# copy required resources into hugo
.PHONY: docs-copy
docs-copy:
	@echo
	@echo 'Syncing required doc files and directories for the build process.'
# brace expansion not with sh
	@rsync --delete -ax --exclude={hugo,Makefile,.gitignore,README.md} ./ $(HUGO)/content/$(NAME)

# clean up docs build artifacts (removing the hugo folder)
.PHONY: docs-clean
docs-clean:
# NOTE: the whitespace at the beginning of the ifeq/else/endif statements is intentional.
	@echo
 	ifneq ($(shell id -u), 0)
		@echo "You must run the clean command using sudo or as root because some doc files are created by root"
 	else
		@echo 'Clean up docs build artifacts'
		@rm -rf $(HUGO)
		@echo "Removed folder: $(HUGO)"
 	endif
