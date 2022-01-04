# bingo creates symlinks from the -l option in GOBIN, from where
# we can easily use it with buf. To have the symlinks inside this
# repo and on a known location, we set GOBIN to .bingo in the root
# of the repository (therefore we need to cd ..)
.PHONY: protoc-deps
protoc-deps: $(BINGO)
	@cd .. && GOPATH="" GOBIN=".bingo" $(BINGO) get -l google.golang.org/protobuf/cmd/protoc-gen-go
	@cd .. && GOPATH="" GOBIN=".bingo" $(BINGO) get -l github.com/asim/go-micro/cmd/protoc-gen-micro/v4
	@cd .. && GOPATH="" GOBIN=".bingo" $(BINGO) get -l github.com/owncloud/protoc-gen-microweb
	@cd .. && GOPATH="" GOBIN=".bingo" $(BINGO) get -l github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	@cd .. && GOPATH="" GOBIN=".bingo" $(BINGO) get -l github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

.PHONY: buf-generate
buf-generate: $(BUF) protoc-deps
	$(BUF) generate

