#!/usr/bin/env bash
# Shared settings and helpers for the NATS JWT test scripts. Source this file.
#
# Variables you can set:
#   NATS_JWT_DIR           folder for server.conf and the .creds files (required)
#   NATS_CONTAINER         name of the NATS container (default: nats-jwt)
#   NATS_DOCKER_NETWORK    "host" (default, for Linux) or "publish" (for Docker Desktop)
#   NATS_URL               client URL of the NATS server (default: nats://127.0.0.1:9233)
#   NATS_CREDS_BIN         path of a built natscreds binary. Without it, start-nats.sh runs "go run".
#   NATS_CLI               path of a local "nats" binary. Without it, the scripts use the nats-box image.

# The image pins have a tag and a digest, like the pins of the actions in the workflows.
NATS_IMAGE="${NATS_IMAGE:-nats:2.12-alpine@sha256:b270f5e2428354c0335612694d7dd2fb588148e567a5757fdff325ef9c9332e6}"
NATS_BOX_IMAGE="${NATS_BOX_IMAGE:-natsio/nats-box:0.19.0@sha256:00f10dd4156f28ece84bb3a52cf356ef220ebaa0494310caa90dec184d9e2b89}"
NATS_CONTAINER="${NATS_CONTAINER:-nats-jwt}"
NATS_DOCKER_NETWORK="${NATS_DOCKER_NETWORK:-host}"
NATS_URL="${NATS_URL:-nats://127.0.0.1:9233}"
NATS_MONITOR_URL="${NATS_MONITOR_URL:-http://127.0.0.1:8222}"

log() { echo "[nats-jwt] $*"; }

# nats_as <user|other|none> <nats arguments...>
# Runs the nats CLI with the creds of the given user. "none" sends no creds.
nats_as() {
  local who="$1"
  shift
  local args=(-s "$NATS_URL")
  if [ -n "${NATS_CLI:-}" ]; then
    [ "$who" != none ] && args+=(--creds "$NATS_JWT_DIR/$who.creds")
    "$NATS_CLI" "${args[@]}" "$@"
  else
    [ "$who" != none ] && args+=(--creds "/creds/$who.creds")
    docker run --rm --network host -v "$NATS_JWT_DIR:/creds:ro" "$NATS_BOX_IMAGE" nats "${args[@]}" "$@"
  fi
}
