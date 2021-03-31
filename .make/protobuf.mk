.PHONY: protoc-gen-openapiv2
protoc-gen-openapiv2:
	GO111MODULE=off go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2


.PHONY: $(PROTO_SRC)/${NAME}.pb.go
$(PROTO_SRC)/${NAME}.pb.go: $(BUF) protoc-gen-openapiv2 $(PROTOC_GEN_GO)
	@echo "$(NAME): generating $(PROTO_SRC)/${NAME}.pb.go"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=../third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-go=$(PROTOC_GEN_GO) \
		--go_out=$(PROTO_SRC) --go_opt=paths=source_relative \
		$(PROTO_SRC)/${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.pb.micro.go
$(PROTO_SRC)/${NAME}.pb.micro.go: $(BUF) protoc-gen-openapiv2 $(PROTOC_GEN_MICRO)
	@echo "$(NAME): generating $(PROTO_SRC)/${NAME}.pb.micro.go"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=../third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-micro=$(PROTOC_GEN_MICRO) \
		--micro_out=$(PROTO_SRC) --micro_opt=paths=source_relative \
		$(PROTO_SRC)/${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.pb.web.go
$(PROTO_SRC)/${NAME}.pb.web.go: $(BUF) protoc-gen-openapiv2 $(PROTOC_GEN_MICROWEB)
	@echo "$(NAME): generating $(PROTO_SRC)/${NAME}.pb.web.go"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=../third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-microweb=$(PROTOC_GEN_MICROWEB) \
		--microweb_out=$(PROTO_SRC) --microweb_opt=paths=source_relative \
		$(PROTO_SRC)/${NAME}.proto

.PHONY: $(PROTO_SRC)/${NAME}.swagger.json
$(PROTO_SRC)/${NAME}.swagger.json: $(BUF) protoc-gen-openapiv2 $(PROTOC_GEN_OPENAPIV2)
	@echo "$(NAME): generating $(PROTO_SRC)/${NAME}.swagger.json"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=../third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-openapiv2=$(PROTOC_GEN_OPENAPIV2) \
		--openapiv2_out=$(PROTO_SRC)/ \
		$(PROTO_SRC)/${NAME}.proto

.PHONY: ../docs/extensions/${NAME}/grpc.md
../docs/extensions/${NAME}/grpc.md: $(BUF) protoc-gen-openapiv2 $(PROTOC_GEN_DOC)
	@echo "$(NAME): generating ../docs/extensions/${NAME}/grpc.md"
	@$(BUF) protoc \
		-I=$(PROTO_SRC)/ \
		-I=../third_party/ \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		--plugin protoc-gen-doc=$(PROTOC_GEN_DOC) \
		--doc_opt=./templates/GRPC.tmpl,grpc.md \
		--doc_out=../docs/extensions/${NAME} \
		$(PROTO_SRC)/${NAME}.proto
