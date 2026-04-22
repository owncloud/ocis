#!/usr/bin/env bash
set -euo pipefail

# Initialize oCIS config directory and copy required config files.
# OCIS_ACTION_PATH: path to the action directory (contains config/)
# The action is used from within a repo checkout so GITHUB_WORKSPACE is set.

CONFIG_DIR="${HOME}/.ocis/config"
mkdir -p "$CONFIG_DIR"

ocis init --insecure true

# app-registry.yaml: bundled with this action so it works without a repo checkout
cp "${OCIS_ACTION_PATH}/config/app-registry.yaml" "${CONFIG_DIR}/app-registry.yaml"

echo "oCIS config initialized at ${CONFIG_DIR}"
