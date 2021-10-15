
.PHONY: protoc-deps
protoc-deps: $(BINGO)
	@# #TODO: bingo creates symlinks from the -l option in the current directory
	@# if no GOPATH and GOBIN is set, but they should reside inside .bingo
	@# for now we move them manually
	@cd .. && GOPATH="" GOBIN="" $(BINGO) get -l google.golang.org/protobuf/cmd/protoc-gen-go
	@cd .. && mv protoc-gen-go .bingo/
	@cd .. && GOPATH="" GOBIN="" $(BINGO) get -l github.com/asim/go-micro/cmd/protoc-gen-micro/v3
	@cd .. && mv protoc-gen-micro .bingo/
	@cd .. && GOPATH="" GOBIN="" $(BINGO) get -l github.com/owncloud/protoc-gen-microweb
	@cd .. && mv protoc-gen-microweb .bingo/
	@cd .. && GOPATH="" GOBIN="" $(BINGO) get -l github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	@cd .. && mv protoc-gen-openapiv2 .bingo/
	@cd .. && GOPATH="" GOBIN="" $(BINGO) get -l github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	@cd .. && mv protoc-gen-doc .bingo/

.PHONY: buf-generate
buf-generate: $(BUF) protoc-deps
	$(BUF) generate

