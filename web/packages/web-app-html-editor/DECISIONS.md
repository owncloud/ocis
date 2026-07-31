# Architecture decisions — `web-app-html-editor`

Each decision cites a finding in `ARCHAEOLOGY.md`. Where archaeology contradicted
the original prompt, the correction and its rationale are recorded — the prompt
required that "every architectural decision must cite a specific file found during
Phase 1," and several of its assumptions (MIME registration, a `useHtmlFile`
load/save composable, a custom dirty guard / `Ctrl+S` / error toasts, a brand-color
override) do not match how oCIS Web actually works.

## D1 — Editor library: **CodeMirror 6**, used directly

Monaco is absent; CodeMirror 6 is already fully resolved in the lockfile
(Arch 1.3). Per the prompt's own rule ("if only CodeMirror is present, use
CodeMirror"), we use CodeMirror 6.

There is **no** reusable raw-CodeMirror Vue component to reuse (Arch 1.3): the only
CodeMirror in the tree is wrapped by `md-editor-v3`, a Markdown editor unsuitable
for HTML *source* editing. We therefore add the individual CodeMirror packages as
direct dependencies at the versions already in the lockfile — `@codemirror/state`,
`@codemirror/view`, `@codemirror/commands`, `@codemirror/language`,
`@codemirror/lang-html` — and write a ~60-line wrapper (`HtmlEditorPane.vue`).
Bundle impact is effectively zero: every one of these is already resolved and
shipped transitively via `web-pkg` → `md-editor-v3`. No web worker is involved
(unlike Monaco), so no Vite worker configuration is required.

## D2 — File-type registration: by **extension** `html`, `htm`

Correction to the prompt: oCIS associates apps to files by **file extension, not
MIME type** (Arch 1.1, 1.2). `defineWebApplication` exposes
`appInfo.extensions[].extension`; there is no `mimeType` field. We register
`html` and `htm` (both are common and neither is claimed today — Arch 1.2).
`application/xhtml+xml` does not apply in an extension-based registry; we instead
also accept the `xhtml` extension. There is no priority field to set (Arch 1.2).
The app id is `html-editor` and the package `name` is `"html-editor"` (Arch 1.5).

## D3 — Preview sandboxing: `srcdoc` + minimal sandbox + self-contained CSP

We inject content via `srcdoc` (never a `src` URL), with `referrerpolicy="no-referrer"`
and **without** `allow-same-origin`, so the preview runs at an opaque origin and
cannot read the shell's cookies, storage or OIDC token (Arch 1.4). That opaque
origin is the primary, repo-local control and holds even with no proxy CSP.

Hardened after the security review (see `SECURITY-REVIEW.md`), which corrected the
prompt's `sandbox="allow-scripts allow-forms allow-popups"`:

- **Sandbox reduced to `allow-scripts` only.** `allow-forms` and `allow-popups`
  added no value for a preview but enabled zero-click phishing / external
  form-POST beaconing from the opaque-origin frame (those channels are independent
  of CSP). Dropping them removes the vector at the source. `allow-scripts` is kept
  so self-contained pages still render and run their inline JS.
- **A strict, iframe-scoped CSP is prepended to the `srcdoc`** (`helpers/preview.ts`,
  applied in `App.vue`). It is prepended as the document's first bytes — not inserted
  into the `<head>` — so the first `<meta>` policy wins and governs the *whole*
  document, including any element placed before a real or decoy `<head>` that
  head-insertion would leave uncovered. `default-src 'none'` denies all network
  egress (so a hostile script cannot beacon/exfiltrate), `form-action 'none'` and
  `base-uri 'none'` are belt-and-suspenders, while inline `script-src`/`style-src`
  and `data:`/`blob:` images, fonts and media keep self-contained previews working.
  This makes the preview self-protecting rather than reliant on the deployment proxy
  CSP. Trade-off: previews of documents that load **external** CSS/JS/images will not
  fetch them (and the preview can still navigate itself away — accepted low residual
  risk F11, see `SECURITY-REVIEW.md`).
- The exact sandbox string and `srcdoc`-not-`src` are pinned by a regression test
  (`HtmlPreviewPane.spec.ts`), so loosening the contract fails CI.

## D4 — Layout: three view modes via CSS **grid**

`SPLIT` (default, 50/50) · `EDITOR` (100%) · `PREVIEW` (100%), toggled from the
in-app toolbar. The container is a CSS grid whose `grid-template-columns` switches
between `1fr 1fr`, `1fr 0`, and `0 1fr` so column widths stay stable. No drag
handle in v1 (v2). The toolbar lives inside the app content; the page chrome
(filename, Save, close, action menu) is the framework's `AppTopBar` (Arch — the
linchpin section, `AppWrapper.vue:7-16`).

## D5 — Dirty state & navigation guard: **provided by the framework, not reimplemented**

Correction to the prompt: we do **not** add an `isDirty` ref, a navigation guard,
or an `UnsavedChangesModal`. `AppWrapper` already computes `isDirty` from
`currentContent !== serverContent` (`AppWrapper.vue:164-166`), guards route-leave
and reload (`AppWrapper.vue:646-669, 168-179`), and shows the dialog. Our only
obligation is to **declare a `currentContent` prop and emit
`update:currentContent`** on every edit (Arch linchpin: `AppWrapper.vue:156-158,
352`). Doing so turns on dirty tracking, the Save button, autosave, and the guard
automatically. Reimplementing any of it would double the modal and the
`beforeunload` handler.

## D6 — Save shortcut: **provided by the framework, not bound in the editor**

Correction to the prompt: `Ctrl+S` / `Cmd+S` is already bound once at the document
level by `AppWrapper` via `useKeyboardActions`, firing `save()` only when dirty
(`AppWrapper.vue:535-541`). CodeMirror's default keymap does **not** bind `Ctrl+S`,
so the keystroke propagates to the wrapper unhindered. We therefore add **no**
save keybinding in `HtmlEditorPane` — adding one would either double-save or
swallow the wrapper's handler. (The editor does keep CodeMirror's standard editing
keymap: history/undo, indent, etc.)

## D7 — Errors: **surfaced by the framework**

Correction to the prompt: load/save/conflict/quota errors are already surfaced by
`AppWrapper` via `useMessages().showErrorMessage`, with explicit `401/403`,
`409/412`, and `507` handling (`AppWrapper.vue:420-479`). The file-too-large case
is a pre-load warning modal driven by `meta.fileSizeLimit`
(`AppWrapper.vue:200-202, 392-414`); we set `meta.fileSizeLimit` in registration
(D-extra) so it applies. The app itself contains no `alert()` and no custom error
modal. The preview pane fails safe (renders empty) on empty content.

## D-extra — Toolbar scope, file-size limit, theme, brand CSS

- **Toolbar = view-mode toggle only.** Filename, the Save button (with its
  dirty/read-only disabling), and the action menu are rendered by `AppTopBar`
  (`AppWrapper.vue:7-16, 543-556`). A second Save button / filename / dirty dot in
  our toolbar would duplicate the framework, so `HtmlToolbar` carries only the
  three-segment Editor | Split | Preview control, built as a `.oc-button-group` of
  `OcButton`s toggling `appearance` (Arch 1.6, the `ViewOptions.vue` pattern).
- **File-size limit:** set `appInfo.meta.fileSizeLimit` (mirroring the text
  editor's `2000000`) so the wrapper's large-file warning applies (D7).
- **Theme:** the editor reads `useThemeStore().currentTheme.isDark` (Arch 1.3, 1.6)
  and reconfigures a CodeMirror theme compartment; the editor chrome is expressed
  in ODS tokens (`--oc-color-background-default`, `--oc-color-text-default`,
  `--oc-color-border`) so it follows the active theme.
- **No brand-color override.** Correction to the prompt's Phase 4 CSS block: an
  upstream app must inherit the deployment's theme, not hardcode
  `--oc-color-swatch-primary-default: #041E42` or a Source Sans font. Those tokens
  are owned by the theme (Arch 1.6); overriding them in the app would break themed
  deployments. We therefore do not ship that `:root` override.

## Scaffold correction

The real reference has **no `vite.config.ts` and no `tsconfig.json`** (Arch 1.1,
1.5); apps are auto-discovered and built by the root config. We mirror the real
structure (`package.json`, `vitest.config.ts`, `l10n/`, `src/`, `tests/`) and do
**not** create those two files, contrary to the prompt's Phase 3 tree.
