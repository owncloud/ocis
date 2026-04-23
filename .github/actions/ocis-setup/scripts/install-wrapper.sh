#!/usr/bin/env bash
set -euo pipefail

# ociswrapper is not published as a release artifact.
# Sparse-clone only tests/ociswrapper/ from the matching tag and build it.
VERSION="${OCIS_VERSION:-latest}"

if [[ "$VERSION" == "latest" ]]; then
  if [[ -f /tmp/ocis-resolved-version ]]; then
    VERSION=$(cat /tmp/ocis-resolved-version)
  else
    VERSION=$(curl -s https://api.github.com/repos/owncloud/ocis/releases/latest \
      | python3 -c "import sys,json; print(json.load(sys.stdin)['tag_name'].lstrip('v'))")
  fi
fi

CLONE_DIR="/tmp/ocis-src-wrapper"
rm -rf "$CLONE_DIR"

echo "Sparse-cloning ociswrapper from tag v${VERSION}..."
git clone --depth=1 --filter=blob:none --sparse \
  --branch "v${VERSION}" \
  https://github.com/owncloud/ocis.git "$CLONE_DIR"

cd "$CLONE_DIR"
git sparse-checkout set tests/ociswrapper

cd tests/ociswrapper
echo "Building ociswrapper..."
GOWORK=off go build -o /tmp/ociswrapper .
sudo mv /tmp/ociswrapper /usr/local/bin/ociswrapper

echo "ociswrapper installed at $(which ociswrapper)"
rm -rf "$CLONE_DIR"
