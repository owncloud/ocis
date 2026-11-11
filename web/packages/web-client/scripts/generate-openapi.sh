#!/usr/bin/env sh
set -eu

# Regenerates the libre-graph client under src/graph/generated.
#
# `readOnly: true` in the spec is load-bearing here: the typescript-fetch templates omit
# those fields from the generated `*ToJSON` serializers, so a field that oCIS does accept
# on write but that the spec marks read-only is silently dropped from the request body.
# Note that a `readOnly` next to a `$ref` has no effect (OpenAPI 3.0 ignores siblings of
# `$ref`) — what propagates is the annotation on the referenced schema itself.
#
# oCIS returns a field the spec does not declare: `attributes` on users, added by
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
