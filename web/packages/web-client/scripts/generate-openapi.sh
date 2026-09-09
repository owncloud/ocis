#!/usr/bin/env sh
set -eu

# Regenerates the libre-graph client under src/graph/generated.
#
# The spec marks many fields `readOnly: true`, including every field of `Quota` and the
# `parentReference` identifiers. The typescript-fetch templates take that literally: they
# emit `readonly` modifiers and, more importantly, omit those fields from the generated
# `*ToJSON` serializers, so a PATCH body such as `{ quota: { total: 500 } }` would go out
# as `{ quota: {} }`. oCIS does accept these fields on write, so the annotation is stripped
# from the spec before generating.

SPEC_URL="https://raw.githubusercontent.com/owncloud/libre-graph-api/main/api/openapi-spec/v1.0.yaml"
GRAPH_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")/../src/graph" && pwd)"
SPEC_FILE="$GRAPH_DIR/openapi-spec.yaml"

cleanup() {
  rm -f "$SPEC_FILE"
}
trap cleanup EXIT

rm -rf "$GRAPH_DIR/generated"
curl -sSfL "$SPEC_URL" | sed '/^ *readOnly: true$/d' >"$SPEC_FILE"

docker run --rm -v "$GRAPH_DIR:/local" openapitools/openapi-generator-cli generate \
  -i /local/openapi-spec.yaml \
  -g typescript-fetch \
  --type-mappings=DateTime=string \
  -o /local/generated
