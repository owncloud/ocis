#!/usr/bin/env bash
set -euo pipefail

# OCIS_VERSION: "latest" or a specific version like "8.0.1"
VERSION="${OCIS_VERSION:-latest}"

if [[ "$VERSION" == "latest" ]]; then
  VERSION=$(curl -s https://api.github.com/repos/owncloud/ocis/releases/latest \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['tag_name'].lstrip('v'))")
fi

# Write resolved version so install-wrapper.sh uses the same tag without a second API call.
echo "$VERSION" > /tmp/ocis-resolved-version

echo "Installing oCIS $VERSION..."
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH_SUFFIX="amd64" ;;
  aarch64) ARCH_SUFFIX="arm64" ;;
  *)       echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

BASE_URL="https://github.com/owncloud/ocis/releases/download/v${VERSION}"
BINARY="ocis-${VERSION}-linux-${ARCH_SUFFIX}"

curl -sLo /tmp/ocis "${BASE_URL}/${BINARY}"
curl -sLo /tmp/ocis.sha256 "${BASE_URL}/${BINARY}.sha256"

# sha256 file contains "HASH  filename" — rewrite to match our local path
HASH=$(awk '{print $1}' /tmp/ocis.sha256)
echo "${HASH}  /tmp/ocis" | sha256sum -c -

chmod +x /tmp/ocis
sudo mv /tmp/ocis /usr/local/bin/ocis

echo "oCIS $(ocis --version 2>&1 | head -1) installed."
