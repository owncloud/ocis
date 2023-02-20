#!/bin/bash

# generate hash for .bingo folder
BINGO_HASH=$(find .bingo -type f -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d ' ' -f 1)

echo "[INOF] BINGO_HASH: $BINGO_HASH"

URL="$CACHE_ENDPOINT/$CACHE_BUCKET/ocis/go-bin/$BINGO_HASH.tar.gz"

echo "[INFO] Checking for the go bin cache at '$URL'."

if curl --output /dev/null --silent --head --fail "$URL"; then
    echo "[INFO] Go bin cache with has '$BINGO_HASH' exists."
    # https://discourse.drone.io/t/how-to-exit-a-pipeline-early-without-failing/3951
    # exit a Pipeline early without failing
    exit 78
else
    echo "[INFO] Go bin cache with has '$BINGO_HASH' does not exist."
fi
