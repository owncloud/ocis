---
description: Investigate top flaky tests from the owncloud/ocis weekly pipeline report
argument-hint: "[workflow-run-url-or-local-json-path]"
allowed-tools: [Bash, Read, Glob, Grep, WebFetch, Agent]
---

Investigate flaky tests from the owncloud/ocis weekly pipeline CI report.

The argument (if provided) is: $ARGUMENTS

## Phase 1 – Fetch CI failure data

The weekly pipeline report workflow (`weekly_pipeline_report.yml`) runs a Go tool that emits a
JSON summary directly to the job log. There are no artifacts — the data is in the log.

**Step 1a – resolve the run ID.**

If $ARGUMENTS is a GitHub run URL (e.g. `https://github.com/owncloud/ocis/actions/runs/12345`),
extract the numeric run ID from the URL.

If $ARGUMENTS is a local file path, read it directly with the Read tool and skip to Phase 2.

If $ARGUMENTS is empty, fetch the most recent completed run:

```bash
gh api "repos/owncloud/ocis/actions/workflows/weekly_pipeline_report.yml/runs?per_page=5&status=completed" \
  --jq '.workflow_runs[0] | {id, conclusion, created_at, html_url}' 2>/dev/null
```

**Step 1b – get the "Generate Pipeline Report" job ID for that run:**

```bash
gh api "repos/owncloud/ocis/actions/runs/RUN_ID/jobs" \
  --jq '.jobs[] | select(.name == "Generate Pipeline Report") | {id, conclusion}' 2>/dev/null
```

**Step 1c – download the job log and extract the JSON report:**

The log ends with a JSON block starting at `{` after the line `Result: N failed / M total`.
Save the raw log then use Python to extract the JSON reliably:

```bash
gh api "repos/owncloud/ocis/actions/jobs/JOB_ID/logs" 2>/dev/null \
  > /tmp/ocis_raw.log

python3 - <<'EOF'
import re, json, sys

raw = open('/tmp/ocis_raw.log').readlines()
# Find the "Result:" line
start = next((i for i, l in enumerate(raw) if 'Result:' in l), None)
if start is None:
    print("ERROR: Result line not found", file=sys.stderr); sys.exit(1)

# Strip timestamps, fix secret redactions
lines = []
for l in raw[start+1:]:
    l = re.sub(r'^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z ', '', l)
    lines.append(l)
text = ''.join(lines)
text = text.replace('***"', '[REDACTED]')

# raw_decode stops at end of valid JSON, ignoring trailing log lines
data, _ = json.JSONDecoder().raw_decode(text)
with open('/tmp/ocis_report.json', 'w') as f:
    json.dump(data, f, indent=2)
print(f"OK: {len(data.get('count', []))} count entries, {len(data.get('pipelines', []))} pipelines")
EOF
```

If the API call fails or the Python script exits with an error, ask the user:
> "I couldn't fetch the CI report automatically. Please provide either:
> - A local path to the report JSON file
> - The URL of a specific GitHub Actions run (e.g. `https://github.com/owncloud/ocis/actions/runs/12345`)
> - The raw JSON pasted directly"

Once you have the JSON (either from the log or from the user), proceed to Phase 2.

**Step 1d – extract failure logs for the top-ranked steps:**

After `/tmp/ocis_report.json` is written, extract the `logs` field for each failing step entry
from the `pipelines` array and save them to `/tmp/ocis_step_logs.json`:

```bash
python3 - <<'EOF'
import json, sys

report = json.load(open('/tmp/ocis_report.json'))
count  = report.get('count', [])

# Build a lookup: (stage_name, step_name) -> logs
logs_by_key = {}
for pipeline in report.get('pipelines', []):
    for stage in pipeline.get('pipeline_info', {}).get('pipeline_stages', []):
        for step in stage.get('steps', []):
            if step.get('status') == 'failure':
                key = (step.get('stage_name', ''), step.get('step_name', ''))
                logs_by_key[key] = step.get('logs', '')

# Attach logs to top-5 count entries (excluding CI plumbing — same filter applied in Phase 2)
EXCLUDE = {'scan-result-cache', 'codacy', 'upload', 'publish', 'tarball',
           'output', 'diff', 'clone', 'checkout', 'generate'}
def is_test_step(name):
    n = name.lower()
    if any(n.startswith(e) for e in EXCLUDE): return False
    return any(k in n for k in ('test', 'e2e', 'api', 'cli'))

test_entries = [e for e in count if is_test_step(e.get('step_name', ''))]
test_entries.sort(key=lambda e: e.get('count', 0), reverse=True)
top5 = test_entries[:5]

result = []
for e in top5:
    key = (e.get('stage_name', ''), e.get('step_name', ''))
    result.append({
        'stage_name': e['stage_name'],
        'step_name':  e['step_name'],
        'count':      e['count'],
        'logs':       logs_by_key.get(key, ''),
    })

with open('/tmp/ocis_step_logs.json', 'w') as f:
    json.dump(result, f, indent=2)

for i, r in enumerate(result, 1):
    preview = (r['logs'] or '(no logs)')[:120].replace('\n', ' ')
    print(f"  {i}. [{r['count']}x] {r['stage_name']} / {r['step_name']}")
    print(f"     {preview}")
EOF
```

The output of this script provides the per-step log previews shown in the Phase 2 triage table,
and the full `logs` content is passed verbatim to each Phase 3 deep-dive agent.

## Phase 2 – Triage: rank top 5

The JSON has a `count` array of objects: `{stage_name, step_name, count}`.

**Exclude** these step names — they are CI plumbing, not test failures:
- `scan-result-cache`, `coverage-cache-*`, `restore_*`, `sync-from-cache`
- `codacy`, `upload`, `publish`, `tarball`, `output`, `diff`
- `clone`, `checkout`, `generate`, `build*`, `lint*`, `docs-*`
- `wait-for-*`, `health-check-*`, `restore-*`, `unzip-*`
- Any step whose name does not involve running actual tests

**Keep** entries whose step name contains: `run-api-tests`, `localApiTests-*`, `e2e-tests`,
`cli-tests`, `unit-tests`, `test`.

Sort the kept entries by `count` descending, take top 5.

Print a ranked table:

```
Rank | Failures | Stage / Step
-----|----------|--------------------------------------
  1  |    N     | stage-name / step-name
  2  |    N     | ...
  ...
```

Also print the overall summary line:
> `N failed / M total commits (P%) in the past 14 days`

Then ask the user:
> "Which of the top 5 should I investigate? Enter one or more numbers (e.g. `1` or `1,3,5`), or `all` for all five."

Wait for the user's selection before proceeding.

## Phase 3 – Deep-dive per selected offender

For each selected offender, launch a parallel Explore agent (one Agent call per offender).
Run all selected offenders in parallel.

Each agent receives:
- `stage_name` and `step_name`
- The full `logs` field extracted in Step 1d (Behat output: scenario title, feature file path,
  failing step text, and assertion error)

Because the failure log is already known, the agent must **not** search for it from scratch.
Instead:

1. **Parse the provided failure log** to extract:
   - The exact `Scenario:` title
   - The feature file path and line number (already in the log)
   - The failing Gherkin step text
   - The assertion error / exception message

2. **Find the step definition** implementing the failing Gherkin step:
   - Search `tests/acceptance/features/bootstrap/` for the PHP file and method matching the step.
   - Note `file:line` for the matching step method.

3. **Find the Go server handler** invoked by that step:
   - Trace the HTTP call (WebDAV, Graph API, CS3 gRPC) from the step definition to the Go backend.
   - Search `services/` for the relevant service handler.
   - Search `vendor/github.com/cs3org/reva/` for lower-level handlers.
   - Note `file:line` and function name.

4. **Identify the flakiness mechanism**. Classify as one of:
   - `timing` — test doesn't wait long enough for an async operation
   - `race` — concurrent access without synchronization
   - `ordering` — assumes deterministic result ordering
   - `stale-cache` — reads from an index/cache before it's updated
   - `resource-leak` — leftover state from a previous test
   - `env-dep` — depends on external service availability or network
   - `other` — describe specifically

5. **Return a structured result:**
   ```
   STAGE: <stage-name>
   STEP: <step-name>
   SCENARIO: <full scenario title or "unknown — suite-level failure">
   FEATURE_FILE: <path>:<line>
   STEP_DEFINITION: <path>:<line> — <step text>
   GO_HANDLER: <path>:<line> — <function name>
   MECHANISM: <classification>
   FAILURE_ASSERTION: <exact assert/error from logs if available>
   EVIDENCE: <1-3 sentences describing what you found in the code>
   ```

## Phase 4 – Synthesize

For each investigated offender, produce a report section:

### [Rank]. `<Stage / Step>`

**Root cause** (one sentence): ...

**Evidence**

| What | Where | Detail |
|------|-------|--------|
| Failing assertion | `file:line` | exact error text |
| Step definition | `file:line` | step text |
| Server handler | `file:line` | function name |

**5-Why analysis**

1. Why did the test fail? — ...
2. Why did that condition occur? — ...
3. Why isn't that condition prevented? — ...
4. Why wasn't this caught earlier? — ...
5. Why does this recur? — ...

**Fix options** (least to most invasive)

1. *Test-side*: ...
2. *Test-side (alternative)*: ...
3. *Server-side*: ...

**Recommended fix**: Option N — rationale in one sentence.

---

## Phase 5 – Output

Print the full synthesized report for all investigated offenders.

Then ask the user:
> "Should I write this report to a markdown file? If yes, provide a filename or I'll use `flaky-test-report-YYYY-MM-DD.md`."

If the user confirms, write the report to the specified (or default) filename in the current directory.
