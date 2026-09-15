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
#
# oCIS also returns a field the spec does not declare: `attributes` on users, added by
# `UserWithAttributes` in services/graph and filled from
# OCIS_USER_SEARCH_DISPLAYED_ATTRIBUTES. The typescript-fetch templates rebuild every
# response from the declared fields only, so an undeclared field is dropped during decoding.
# It is declared on the read model below to keep it.

SPEC_URL="https://raw.githubusercontent.com/owncloud/libre-graph-api/main/api/openapi-spec/v1.0.yaml"
GRAPH_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")/../src/graph" && pwd)"
SPEC_FILE="$GRAPH_DIR/openapi-spec.yaml"

# The single line the `user` schema composes `userUpdate` into. The two other references to
# that schema are request bodies and sit at a different indentation.
USER_ALL_OF_LINE="        - \$ref: '#/components/schemas/userUpdate'"

cleanup() {
  rm -f "$SPEC_FILE"
}
trap cleanup EXIT

rm -rf "$GRAPH_DIR/generated"
curl -sSfL "$SPEC_URL" \
  | sed '/^ *readOnly: true$/d' \
  | awk -v anchor="$USER_ALL_OF_LINE" '
      { print }
      $0 == anchor {
        print "        - type: object"
        print "          properties:"
        print "            attributes:"
        print "              type: array"
        print "              items:"
        print "                type: string"
        print "              description: Attributes of the user as configured via OCIS_USER_SEARCH_DISPLAYED_ATTRIBUTES. Not part of the upstream spec, added by oCIS. Read-only."
        patched = 1
      }
      END { if (!patched) exit 1 }
    ' >"$SPEC_FILE" || {
  echo "failed to declare the oCIS-only user 'attributes' field: the spec no longer contains" >&2
  echo "the expected line, adapt USER_ALL_OF_LINE to how the 'user' schema now composes" >&2
  echo "'userUpdate'. Regenerating without it silently drops the field from user responses." >&2
  exit 1
}

docker run --rm -v "$GRAPH_DIR:/local" openapitools/openapi-generator-cli generate \
  -i /local/openapi-spec.yaml \
  -g typescript-fetch \
  --type-mappings=DateTime=string \
  -o /local/generated
