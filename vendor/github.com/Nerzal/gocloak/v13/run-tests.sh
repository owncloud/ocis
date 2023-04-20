#!/bin/sh

docker-compose down
docker-compose up -d

sleep 10

ARGS=()
if [ $# -gt 0 ]; then
    ARGS+=("-run")
    ARGS+=("^($@)$")
fi

go test -failfast -race -cover -coverprofile=coverage.out -covermode=atomic -p 10 -cpu 1,2 -bench . -benchmem ${ARGS[@]}

docker-compose down