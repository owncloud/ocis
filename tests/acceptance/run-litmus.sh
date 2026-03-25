#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OCIS_BIN="$REPO_ROOT/ocis/bin/ocis"
WRAPPER_BIN="$REPO_ROOT/tests/ociswrapper/bin/ociswrapper"
OCIS_URL="https://localhost:9200"
OCIS_CONFIG_DIR="$HOME/.ocis/config"

# build
make -C "$REPO_ROOT/ocis" build
GOWORK=off make -C "$REPO_ROOT/tests/ociswrapper" build

# init + start ocis
"$OCIS_BIN" init --insecure true
cp "$REPO_ROOT/tests/config/drone/app-registry.yaml" "$OCIS_CONFIG_DIR/app-registry.yaml"

OCIS_URL=$OCIS_URL \
OCIS_CONFIG_DIR=$OCIS_CONFIG_DIR \
STORAGE_USERS_DRIVER=ocis \
PROXY_ENABLE_BASIC_AUTH=true \
OCIS_EXCLUDE_RUN_SERVICES=idp \
OCIS_LOG_LEVEL=error \
IDM_CREATE_DEMO_USERS=true \
IDM_ADMIN_PASSWORD=admin \
OCIS_ASYNC_UPLOADS=true \
OCIS_EVENTS_ENABLE_TLS=false \
NATS_NATS_HOST=0.0.0.0 \
NATS_NATS_PORT=9233 \
OCIS_JWT_SECRET=some-ocis-jwt-secret \
WEB_UI_CONFIG_FILE="$REPO_ROOT/tests/config/drone/ocis-config.json" \
  "$WRAPPER_BIN" serve \
    --bin "$OCIS_BIN" \
    --url "$OCIS_URL" \
    --admin-username admin \
    --admin-password admin &
WRAPPER_PID=$!
trap "kill $WRAPPER_PID 2>/dev/null || true" EXIT

echo "Waiting for ocis..."
timeout 300 bash -c \
  "while [ \$(curl -sk -uadmin:admin $OCIS_URL/graph/v1.0/users/admin \
    -w %{http_code} -o /dev/null) != 200 ]; do sleep 1; done"
echo "ocis ready."

# setup: creates test folder, share, and exports SPACE_ID + PUBLIC_TOKEN to .env
TEST_SERVER_URL=$OCIS_URL bash "$REPO_ROOT/tests/config/drone/setup-for-litmus.sh"
source .env

# run litmus against each WebDAV endpoint
ENDPOINTS=(
  "$OCIS_URL/remote.php/webdav"
  "$OCIS_URL/remote.php/dav/files/admin"
  "$OCIS_URL/remote.php/dav/files/admin/Shares/new_folder/"
  "$OCIS_URL/remote.php/webdav/Shares/new_folder/"
  "$OCIS_URL/remote.php/dav/spaces/$SPACE_ID"
)

for ENDPOINT in "${ENDPOINTS[@]}"; do
  echo "Testing endpoint: $ENDPOINT"
  docker run --rm --network host \
    -e LITMUS_URL="$ENDPOINT" \
    -e LITMUS_USERNAME=admin \
    -e LITMUS_PASSWORD=admin \
    -e TESTS="basic copymove props http" \
    owncloudci/litmus:latest \
    /usr/local/bin/litmus-wrapper
done
