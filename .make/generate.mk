.PHONY: generate
generate: ci-node-generate ci-go-generate

.PHONY: embed.yml
embed.yml: $(FILEB0X)
	@cd pkg/assets/ && echo -n "$(NAME): embed.go - " && $(FILEB0X) embed.yml
