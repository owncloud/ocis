#!/usr/bin/env bash
# Start a stand-alone NATS server in operator (JWT) mode with JetStream.
#
# The script creates new keys at each run and writes server.conf, user.creds, and
# other.creds to $NATS_JWT_DIR. It does not print the content of a .creds file.
# See lib.sh for the variables.
set -euo pipefail
cd "$(dirname "$0")"
# shellcheck source=lib.sh
. ./lib.sh

: "${NATS_JWT_DIR:?set NATS_JWT_DIR}"
REPO_ROOT="$(cd ../.. && pwd)"

mkdir -p "$NATS_JWT_DIR"
if [ -n "${NATS_CREDS_BIN:-}" ]; then
  "$NATS_CREDS_BIN" -out "$NATS_JWT_DIR"
else
  (cd "$REPO_ROOT" && GOWORK=off go run ./tests/nats-jwt/natscreds -out "$NATS_JWT_DIR")
fi

docker rm -f "$NATS_CONTAINER" >/dev/null 2>&1 || true
# The new container has an empty log. Thus verify.sh must not subtract the errors of its old checks.
rm -f "$NATS_JWT_DIR/probe-errors"

network=(--network host)
if [ "$NATS_DOCKER_NETWORK" = publish ]; then
  network=(-p 127.0.0.1:9233:9233 -p 127.0.0.1:8222:8222)
fi

log "starting $NATS_IMAGE as $NATS_CONTAINER"
docker run -d --name "$NATS_CONTAINER" "${network[@]}" \
  -v "$NATS_JWT_DIR/server.conf:/etc/nats/server.conf:ro" \
  "$NATS_IMAGE" -c /etc/nats/server.conf -m 8222 >/dev/null

# Wait until JetStream is ready.
for _ in $(seq 1 60); do
  if curl -sf "$NATS_MONITOR_URL/healthz?js-enabled-only=true" >/dev/null 2>&1; then
    log "NATS is ready. user creds: $NATS_JWT_DIR/user.creds"
    exit 0
  fi
  sleep 1
done
log "NATS did not become ready within 60 s"
docker logs "$NATS_CONTAINER" || true
exit 1
