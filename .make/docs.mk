.PHONY: config-docs-generate
config-docs-generate: $(FLAEX)
	@echo "$(NAME): generating config docs"
	@$(FLAEX) >| ../docs/extensions/$(NAME)/configuration.md

.PHONY: grpc-docs-generate
grpc-docs-generate: ../docs/extensions/${NAME}/grpc.md
