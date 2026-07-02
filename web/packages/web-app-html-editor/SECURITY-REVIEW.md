# Security review - `web-app-html-editor`

This app renders attacker-influenceable HTML (a user may open an `.html` file that
was shared to them or synced from an untrusted source), so it had a dedicated
adversarial security review before merge. This file records the threat model, the
findings, and what was remediated.

## Threat model

Primary threat: a victim opens an `.html`/`.htm`/`.xhtml` file whose contents are
attacker-controlled. The previewed HTML is hostile. Attacker goals, by value:

1. steal the victim's OIDC access token or act as the victim against the oCIS API
2. execute script in the oCIS web origin (XSS of the parent app / iframe break-out)
3. read or exfiltrate cookies, storage, the token, or other files
4. phishing / open-redirect / UI redress via the preview
5. denial of service of the editor tab

Two facts the analysis leans on: the preview iframe has no `allow-same-origin` (so
the `srcdoc` document runs at an opaque origin), and oCIS web auth is an OIDC
**bearer token held in parent-origin JS**, not a cookie.

## Verdict

No high or critical findings. Goals 1-3 are structurally prevented by the opaque
origin plus the bearer-token-in-JS auth model, and verified at the framework
level: Vue binds `:srcdoc` as a property (not an HTML-parsing sink), the editor is
a plaintext CodeMirror (it never executes the HTML), the app adds no parent-side
`postMessage` listener and no custom network egress. Residual risk was goal 4
(phishing/beacon) and goal 5 (tab DoS); both have been reduced by the hardening
below.

## Findings and remediation

| Id | Severity | Finding | Status |
|----|----------|---------|--------|
| F1/F4 | low-med | Hostile preview script / very large doc can hang the victim's tab; size limit was advisory and split view auto-rendered on open | **Fixed**: live preview is paused above `PREVIEW_SIZE_LIMIT` until the user explicitly opts in (`helpers/preview.ts`, `App.vue`) |
| F2 | low | `allow-forms` + `allow-popups` enabled zero-click phishing / external form-POST beaconing from the preview | **Fixed**: sandbox reduced to `allow-scripts` only (`HtmlPreviewPane.vue`) |
| F5/F8 | info | App relied on the deployment proxy CSP; `data:`/`blob:` exfil channels survived the opaque origin | **Fixed**: a strict iframe-scoped CSP (`default-src 'none'; form-action 'none'; ...`) is **prepended** to the `srcdoc` so the first `<meta>` policy wins and governs the whole document — including any element placed before a real or decoy `<head>` (`helpers/preview.ts`) |
| F3 | low | The sandbox regression test pinned only one token | **Fixed**: the full sandbox contract and `srcdoc`-not-`src` are now asserted (`HtmlPreviewPane.spec.ts`) |
| F7 | info | web-runtime delegated-auth `postMessage` handler skips the origin check when unconfigured | **Out of scope** (web-runtime, not reachable in this app's deployment) |
| F9 | info | `AppWrapper` save triggers gate on `isDirty` not `isReadOnly` | **Out of scope** (web-pkg; not reachable via hostile file content here, since the editor enforces read-only) |
| F10 | info | Parent-origin clickjacking backstop (`X-Frame-Options`/COOP) | **Out of scope** (deployment/proxy concern) |
| F11 | low | The sandboxed preview can navigate **itself** (`<meta http-equiv="refresh">`, `location` assignment) to an external URL: CSP `default-src 'none'` does not restrict a frame navigating itself, and the sandbox blocks only **top-frame** navigation (`allow-scripts` is granted, `allow-top-navigation` is not) | **Accepted residual risk (v1)**: such a navigation can carry only data the attacker already controls (a phishing / open-confirmation beacon), never the victim's token, cookies, or storage — those live in the parent origin and the bearer-token-in-JS model keeps them out of the opaque-origin frame. Recorded so the "network-isolated" claim is precise: it isolates *exfiltration of victim secrets*, not *self-navigation of the hostile frame*. |

## Verified safe (no change required)

- Opaque-origin isolation defeats token/cookie/storage theft, authenticated API
  calls as the victim, parent-DOM XSS, top-frame navigation, and popup escape.
- No `v-html`/`innerHTML`/`eval`/`document.write` anywhere in the package; the
  editor is plaintext (CodeMirror syntax highlighting only).
- No custom `fetch`/XHR/WebSocket, no content/token logging, no URL construction
  or open-redirect from file content; registration is by extension.
- Editor logic: no reactive save loop, debounce timer cleared on unmount, autosave
  uses `.drop()`, and read-only is enforced inside the editor.
- Supply chain: CodeMirror closure has no install scripts and is sha512-pinned.

The trade-off introduced by the injected CSP: previews of documents that reference
**external** stylesheets/scripts/images will not load them. This is intentional
for v1 to keep the preview self-protecting and network-isolated. The one channel
the CSP cannot close is the frame navigating *itself* away (F11); that is an
accepted low residual risk because it can leak only attacker-controlled data,
never the victim's credentials.
