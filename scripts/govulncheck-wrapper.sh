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
trap 'rm -f "$TMPFILE"' EXIT

echo "Running govulncheck..."
"$GOVULNCHECK" -format json ./... > "$TMPFILE" 2>&1 || true

python3 - "$TMPFILE" <<'PYEOF'
import json
import sys

def parse_json_stream(path):
    with open(path) as f:
        content = f.read()
    decoder = json.JSONDecoder()
    idx = 0
    objects = []
    while idx < len(content):
        while idx < len(content) and content[idx] in ' \t\n\r':
            idx += 1
        if idx >= len(content):
            break
        obj, end_idx = decoder.raw_decode(content, idx)
        idx = end_idx
        objects.append(obj)
    return objects

objects = parse_json_stream(sys.argv[1])

# Collect OSV details
osvs = {}
for obj in objects:
    if 'osv' in obj:
        osv = obj['osv']
        osvs[osv['id']] = osv

# Collect findings
findings = [obj['finding'] for obj in objects if 'finding' in obj]

# Group by vuln ID
from collections import defaultdict
by_vuln = defaultdict(list)
for f in findings:
    by_vuln[f['osv']].append(f)

fail_vulns = []
warn_vulns = []

for vid, entries in sorted(by_vuln.items()):
    # Check if any trace reaches symbol level (has 'function' in trace frames)
    is_called = any(
        any('function' in frame for frame in entry.get('trace', []))
        for entry in entries
    )
    fixed_version = entries[0].get('fixed_version', '')
    trace = entries[0].get('trace', [])
    module = trace[0].get('module', '') if trace else ''

    # Determine category
    if not is_called:
        category = "IMPORTED"
    elif not fixed_version:
        category = "NO_FIX"
    elif module == "stdlib":
        category = "STDLIB"
    else:
        category = "FIXABLE"

    osv = osvs.get(vid, {})
    summary = osv.get('summary', vid)

    info = {
        'id': vid,
        'category': category,
        'module': module,
        'fixed_version': fixed_version,
        'summary': summary,
    }

    if category == "FIXABLE":
        fail_vulns.append(info)
    else:
        warn_vulns.append(info)

# Print warnings
if warn_vulns:
    print("\n⚠ Vulnerabilities acknowledged (not blocking):")
    for v in warn_vulns:
        reason = {
            'NO_FIX': 'no upstream fix available',
            'STDLIB': f'needs Go toolchain upgrade to {v["fixed_version"]}',
            'IMPORTED': 'code does not call vulnerable function',
        }.get(v['category'], v['category'])
        print(f"  {v['id']}: {v['summary']}")
        print(f"    module={v['module']} ({reason})")

# Print failures
if fail_vulns:
    print(f"\n✗ {len(fail_vulns)} fixable vulnerability(ies) found:")
    for v in fail_vulns:
        print(f"  {v['id']}: {v['summary']}")
        print(f"    module={v['module']}, fix: bump to {v['fixed_version']}")
    sys.exit(1)
else:
    print(f"\n✓ No fixable vulnerabilities found ({len(warn_vulns)} acknowledged warnings)")
    sys.exit(0)
PYEOF
