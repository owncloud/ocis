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

# init + start ocis with gRPC gateway exposed for cs3api-validator
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
GATEWAY_GRPC_ADDR=0.0.0.0:9142 \
OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD=false \
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

docker run --rm --network host \
  owncloud/cs3api-validator:0.2.1 \
  /usr/bin/cs3api-validator /var/lib/cs3api-validator --endpoint=localhost:9142
