install:
	go install -v

build:
	go build -v ./...

lint:
	golint ./...
	go vet ./...

test:
	go test -v ./... --cover

deps:
	go get -u gopkg.in/cenkalti/backoff.v1
	go get -u github.com/golang/lint/golint
	go get -u github.com/stretchr/testify

clean:
	go clean
