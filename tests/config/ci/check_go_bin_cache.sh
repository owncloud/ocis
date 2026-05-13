#!/usr/bin/env bash

#
# $1 - root path where .bingo resides
# $2 - name of the cache item
#

ROOT_PATH="$1"
if [ -z "$1" ]; then
  ROOT_PATH="."
fi
BINGO_DIR="$ROOT_PATH/.bingo"

# generate hash of a .bingo folder
BINGO_HASH=$(cat "$BINGO_DIR"/* | sha256sum | cut -d ' ' -f 1)

go_cache=$(mc find s3/$CACHE_BUCKET/ocis/go-bin/$BINGO_HASH/$2 2>&1 | grep 'Object does not exist')

if [[ -z "$go_cache" ]]
then
  echo "[INFO] Go bin cache with has '$BINGO_HASH' exists."
  exit 78
else
  # stored hash of a .bingo folder to '.bingo_hash' file
  echo "$BINGO_HASH" >"$ROOT_PATH/.bingo_hash"
  echo "[INFO] Go bin cache with has '$BINGO_HASH' does not exist."
fi
