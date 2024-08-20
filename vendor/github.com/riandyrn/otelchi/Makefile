.PHONY: *

GO_VERSIONS="1.21 1.22"

# This is the command that will be used to run the tests
go-test:
	go build .
	go test ./...

# This is the command that will be used to run the tests in a Docker container, useful when executing the test locally
test:
	docker build \
		-t go-test \
		--build-arg GO_VERSIONS=${GO_VERSIONS} \
		-f ./test/infras/Dockerfile . && \
		docker run --rm go-test
	
	make test-build-examples

test-build-examples:
	make test-build-basic-example
	make test-build-multi-services-example

test-build-basic-example:
	docker build -f ./examples/basic/Dockerfile .

test-build-multi-services-example:
	docker build -f ./examples/multi-services/back-svc/Dockerfile .
	docker build -f ./examples/multi-services/front-svc/Dockerfile .