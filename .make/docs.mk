
SKIP_CONFIG_DOCS_GENERATE ?= 0
CONFIG_DOCS_BASE_PATH ?= ../docs/extensions

#.PHONY: config-docs-generate
config-docs-generate: $(FLAEX)
#	@if [ $(SKIP_CONFIG_DOCS_GENERATE) -ne 1 ]; then \
#		$(FLAEX) >| $(CONFIG_DOCS_BASE_PATH)/$(NAME)/configuration.md \
#	; fi;

.PHONY: grpc-docs-generate
grpc-docs-generate: buf-generate
