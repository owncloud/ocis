#!/usr/bin/env bash
#
# Wrapper around govulncheck that only fails on fixable, called vulnerabilities.
#
# - CALLED + fixable (non-stdlib) → FAIL (you should bump the dep)
# - CALLED + stdlib fix only      → WARN (needs Go toolchain upgrade)
# - CALLED + no fix available     → WARN (nothing to do yet)
# - imported/required only        → WARN (code doesn't call it)
#
# Usage: scripts/govulncheck-wrapper.sh [govulncheck-binary]
#   If no binary is provided, uses 'govulncheck' from PATH.

set -euo pipefail

GOVULNCHECK="${1:-govulncheck}"
TMPFILE=$(mktemp)
STDERRFILE=$(mktemp)
trap 'rm -f "$TMPFILE" "$STDERRFILE"' EXIT

echo "Running govulncheck..."
"$GOVULNCHECK" -format json ./... > "$TMPFILE" 2>"$STDERRFILE" || true
cat "$STDERRFILE" >&2

# govulncheck -format json outputs one JSON object per line (NDJSON).
# Classify each finding and display results.
jq -rn '
  [inputs] |

  # OSV summary lookup: { id: summary }
  (map(select(has("osv")))
   | map({ (.osv.id): (.osv.summary // .osv.id) })
   | add // {}) as $osvs |

  # Classify every vuln group
  (map(select(has("finding"))) | map(.finding)
   | group_by(.osv)
   | map(
       . as $e |
       ($e[0].osv) as $vid |
       ($e | any(.[].trace[]?; has("function"))) as $called |
       ($e[0].fixed_version // "") as $fixed |
       ($e[0].trace[0]?.module // "") as $mod |
       {
         id:            $vid,
         summary:       ($osvs[$vid] // $vid),
         module:        $mod,
         fixed_version: $fixed,
         category: (
           if   ($called | not) then "IMPORTED"
           elif ($fixed == "")  then "NO_FIX"
           elif ($mod == "stdlib") then "STDLIB"
           else "FIXABLE"
           end
         )
       }
     )
  ) as $vulns |

  ($vulns | map(select(.category != "FIXABLE"))) as $warns |
  ($vulns | map(select(.category == "FIXABLE"))) as $fails |

  # Warnings
  if ($warns | length) > 0 then "\n⚠ Vulnerabilities acknowledged (not blocking):" else empty end,
  ($warns[] |
    "  \(.id): \(.summary)",
    (if   .category == "NO_FIX"   then "    module=\(.module) (no upstream fix available)"
     elif .category == "STDLIB"   then "    module=\(.module) (needs Go toolchain upgrade to \(.fixed_version))"
     else                              "    module=\(.module) (code does not call vulnerable function)"
     end)
  ),

  # Failures or success
  if ($fails | length) > 0 then
    "\n✗ \($fails | length) fixable vulnerability(ies) found:",
    ($fails[] |
      "  \(.id): \(.summary)",
      "    module=\(.module), fix: bump to \(.fixed_version)"
    )
  else
    "\n✓ No fixable vulnerabilities found (\($warns | length) acknowledged warnings)"
  end
' "$TMPFILE"

# Exit 1 if any fixable vulnerabilities
FAIL_COUNT=$(jq -n '
  [inputs]
  | map(select(has("finding"))) | map(.finding)
  | group_by(.osv)
  | map(
      . as $e |
      ($e | any(.[].trace[]?; has("function"))) as $called |
      ($e[0].fixed_version // "") as $fixed |
      ($e[0].trace[0]?.module // "") as $mod |
      select($called and ($fixed != "") and ($mod != "stdlib"))
    )
  | length
' "$TMPFILE")
[ "$FAIL_COUNT" -eq 0 ]
