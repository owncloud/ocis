.PHONY: $(OPENAPI_SRC)/${NAME}.types.go
$(OPENAPI_SRC)/${NAME}.types.go: $(OAPI_CODEGEN)
	@echo "$(NAME): generating $(OPENAPI_SRC)/${NAME}.types.go"
	@$(OAPI_CODEGEN) \
		-generate types \
		-o $(OPENAPI_SRC)/${NAME}.types.go \
		$(OPENAPI_SRC)/${NAME}-${OPENGRAPH_VERSION}.yml

.PHONY: $(OPENAPI_SRC)/${NAME}.server.go
$(OPENAPI_SRC)/${NAME}.server.go: $(OAPI_CODEGEN)
	@echo "$(NAME): generating $(OPENAPI_SRC)/${NAME}.types.go"
	@$(OAPI_CODEGEN) \
		-generate chi-server \
		-o $(OPENAPI_SRC)/${NAME}.server.go \
		$(OPENAPI_SRC)/${NAME}-${OPENGRAPH_VERSION}.yml
