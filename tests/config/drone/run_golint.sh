#!/bin/bash

make -C "$1" ci-golangci-lint

SUCCESS=$?
if [ $SUCCESS -ne 0 ]; then
    echo "[WARN] golangci-lint failed."
    rm -rf /go/*
    make -C "$1" ci-golangci-lint
fi

# make bingo-update :x:
# rm -rf /go/* :x:
# go clean -modcache && go mod tidy :x:
