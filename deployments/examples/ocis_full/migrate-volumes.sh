#!/usr/bin/env bash
#
# One-time helper for upgrading an existing ocis_full deployment from the old
# docker-managed named volumes to the new ./data/ bind mounts.
#
# Docker named volumes are not directly accessible from the host, so for each
# volume that still exists we spin up a throwaway container that mounts both
# the old volume and the new host directory, and copy the data across.
#
# Safe to re-run: volumes that no longer exist (fresh installs) are skipped,
# and targets that already contain data are left untouched.
#
# Usage: run this from the ocis_full example directory, before starting the
# stack with `docker compose up` for the first time after upgrading:
#   ./migrate-volumes.sh

set -euo pipefail

cd "$(dirname "$0")"

# shellcheck disable=SC1091
[ -f .env ] && source .env

MIGRATIONS=(
  "ocis-config:${OCIS_CONFIG_DIR:-./data/ocis-config}"
  "ocis-data:${OCIS_DATA_DIR:-./data/ocis-data}"
  "ocis-apps:./data/ocis-apps"
  "keycloak_postgres_data:./data/keycloak-postgres"
  "minio-data:./data/minio"
  "clamav-db:./data/clamav-db"
  "companion-data:./data/companion"
)

for entry in "${MIGRATIONS[@]}"; do
  vol="${entry%%:*}"
  target="${entry#*:}"

  if ! docker volume inspect "$vol" >/dev/null 2>&1; then
    echo "skip ${vol}: no such docker volume (fresh install?)"
    continue
  fi

  mkdir -p "$target"
  if [ -n "$(ls -A "$target" 2>/dev/null)" ]; then
    echo "skip ${vol}: ${target} already has data, not overwriting"
    continue
  fi

  docker run --rm \
    -v "${vol}:/from" \
    -v "$(cd "$target" && pwd):/to" \
    alpine sh -c "cp -a /from/. /to/"
  echo "migrated ${vol} -> ${target}"
done

vol_names=()
for entry in "${MIGRATIONS[@]}"; do
  vol_names+=("${entry%%:*}")
done

echo
echo "Done. Verify the deployment works with the new ./data/ paths, then you"
echo "can remove the old named volumes with, e.g.:"
echo "  docker volume rm ${vol_names[*]}"
