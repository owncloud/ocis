#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: sbom-scan.sh [--ocis-path <dir>] [--web-path <dir>] [--out <dir>]
  e.g. sbom-scan.sh --ocis-path /path/to/ocis --web-path /path/to/web --out sbom

Produces CycloneDX SBOMs for the oCIS and web source directories (Trivy fs).

Requires: trivy, jq, cyclonedx-gomod
USAGE
}

main() {
  # Default to ocis repo root (this script lives in ocis/scripts)
  OCIS_PATH="$(cd "$(dirname "$0")/.." && pwd)"
  # Default to sibling web repo root
  WEB_PATH="$(cd "$(dirname "$0")/.." && pwd)/../web"
  OUT_DIR="sbom"

  while [[ $# -gt 0 ]]; do
    case "$1" in
      --ocis-path) OCIS_PATH="$2"; shift 2;;
      --web-path) WEB_PATH="$2"; shift 2;;
      --out) OUT_DIR="$2"; shift 2;;
      -h|--help) usage; exit 0;;
      *) echo "Unknown arg: $1" >&2; usage; exit 2;;
    esac
  done

  # Absolute paths for mounts and outputs  
  mkdir -p "$OUT_DIR"
  OUT="$(cd "$OUT_DIR" && pwd)"  
  # Allow override via --ocis-path and --web-path
  OCIS_PATH_ABS="$(cd "$OCIS_PATH" && pwd)"
  WEB_PATH_ABS="$(cd "$WEB_PATH" && pwd)"

  echo "[oCIS] generating SBOM (fs), cyclonedx-gomod"
  go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest
  GOFLAGS=-mod=mod GO111MODULE=on GOWORK=off \
    cyclonedx-gomod mod -licenses -json -output-version 1.6 \
      -output "$OUT/ocis_cyclonedx.cdx.json" "$OCIS_PATH_ABS"
  OCIS_CNT=$(jq -r '.components|length' "$OUT/ocis_cyclonedx.cdx.json")
  echo "[oCIS] components: $OCIS_CNT -> $OUT/ocis_cyclonedx.cdx.json"
  
  # Install Trivy locally (macOS/Linux) into $HOME/.local/bin and update PATH
  mkdir -p "$HOME/.local/bin"
  curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b "$HOME/.local/bin"
  export PATH="$HOME/.local/bin:$PATH"

  # --file-patterns 'gomod:^go\.mod$' \
  echo "[oCIS] generating SBOM (fs), trivy"
  trivy fs --quiet --timeout 10m \
    --skip-dirs node_modules \
    --skip-dirs dist \
    --skip-dirs tests \
    --skip-dirs docs \
    --skip-dirs vendor \
    --skip-dirs vendor-bin \
    --format cyclonedx --output "$OUT/ocis.cdx.json" "$OCIS_PATH_ABS"
  OCIS_CNT=$(jq -r '.components|length' "$OUT/ocis.cdx.json")
  echo "[oCIS] components: $OCIS_CNT -> $OUT/ocis.cdx.json"

  echo "[web ] generating SBOM (fs), trivy"
  trivy fs --quiet --timeout 10m \
    --skip-dirs node_modules \
    --skip-dirs dist \
    --skip-dirs tests \
    --skip-dirs docs \
    --format cyclonedx --output "$OUT/web.cdx.json" "$WEB_PATH_ABS"
  WEB_CNT=$(jq -r '.components|length' "$OUT/web.cdx.json")
  echo "[web ] components: $WEB_CNT -> $OUT/web.cdx.json"
  
  # diff ocis cyclonedx and trivy: cyclonedx go only vs trivy gathering all projects supported: go/nodejs/php/...
  brew install cyclonedx/cyclonedx/cyclonedx-cli
  echo "[diff] CycloneDX vs trivy"
  cyclonedx diff \
    "$OUT/ocis_cyclonedx.cdx.json" \
    "$OUT/ocis.cdx.json" \
    --output-format text --component-versions \
    > "$OUT/ocis_cyclonedx_vs_trivy.diff" 2>&1 || true

  echo "done  -> $OUT"
}

main "$@"
