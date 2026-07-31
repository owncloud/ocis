# Archaeology — `web-app-html-editor`

Phase 1 findings. Every architectural decision in `DECISIONS.md` cites a line below.
All paths are relative to the repository root.

The primary reference is `packages/web-app-text-editor`. The crucial discovery is
that an oCIS editor app is **thin**: registration, WebDAV load/save, dirty
tracking, the save button, `Ctrl+S`, the unsaved-changes guard, and error toasts
are all provided by the framework's `AppWrapper`. The app itself only renders an
editor surface bound to `currentContent`.

---

## 1.1 The reference extension (`web-app-text-editor`)

Full file list (`find packages/web-app-text-editor -type f`):

```
packages/web-app-text-editor/l10n/.tx/config
packages/web-app-text-editor/l10n/translations.json
packages/web-app-text-editor/package.json
packages/web-app-text-editor/src/App.vue
packages/web-app-text-editor/src/index.ts
packages/web-app-text-editor/tests/unit/__snapshots__/app.spec.ts.snap
packages/web-app-text-editor/tests/unit/app.spec.ts
packages/web-app-text-editor/vitest.config.ts
```

Note what is **absent**: there is no `vite.config.ts` and no `tsconfig.json` in the
package. The spec's assumed scaffold (Phase 3) lists both; the real reference has
neither. Apps are built by the root build, not per-package (see 1.5).

### How registration works — by file EXTENSION, not MIME type

`packages/web-app-text-editor/src/index.ts:16` exports
`defineWebApplication({ setup() {...} })` from `@ownclouders/web-pkg`.

- The app declares the file types it opens as **extensions**, not MIME types:
  `ApplicationFileExtension` entries with an `extension` string
  (`src/index.ts:24-70`, e.g. `{ extension: 'txt', label: ... }`).
- These become `appInfo.extensions` (`src/index.ts:122-129`), and the route is
  `'/:driveAliasAndItem(.*)?'` wrapped by
  `AppWrapperRoute(TextEditor, { applicationId: appId })` (`src/index.ts:98-111`).
- `appInfo` carries `id`, `icon`, `color`, `defaultExtension`, and
  `meta.fileSizeLimit` (`src/index.ts:113-121`; the text editor uses `2000000`).
- An optional app-menu item is built with `useOpenEmptyEditor`
  (`src/index.ts:10, 132-149`) so the app can create a new empty file.

There is **no `mimeType` / `mimeTypes` field** in the registration. The grep the
spec suggested (`grep -r "mimeType..." packages/web-pkg packages/web-runtime`)
finds MIME handling only in unrelated preview/route resolution code, never in app
registration. **App↔file association in oCIS is by extension.**

### `App.vue` is a thin shell

`packages/web-app-text-editor/src/App.vue:16-32`:

- Props received from the wrapper: `applicationConfig: AppConfigObject`,
  `currentContent: string`, `isReadOnly?: boolean`, `resource: Resource`
  (`App.vue:20-25`).
- Emits exactly one event: `(e: 'update:currentContent', value)` (`App.vue:26-28`).
- It renders the web-pkg `TextEditor` component bound to `currentContent` and
  re-emits its change (`App.vue:6-12`). It contains **no load, no save, no dirty
  logic** — that all lives upstream in the wrapper.

---

## The framework wrapper — `AppWrapper` (the linchpin)

`packages/web-pkg/src/components/AppTemplates/AppWrapperRoute.ts` is the route
factory; `packages/web-pkg/src/components/AppTemplates/AppWrapper.vue` is the
component. Behaviours, all in `AppWrapper.vue`:

- **Two flags driven by the wrapped component's interface:**
  - `isEditor = wrappedComponent.emits?.includes('update:currentContent')`
    (`AppWrapper.vue:156-158`). Declaring that emit is what enables the Save
    button, `Ctrl+S`, autosave, and the nav guard.
  - `hasProp('currentContent')` (`AppWrapper.vue:160-162, 352`) gates the WebDAV
    GET. Declaring that prop is what makes the wrapper fetch file contents.
- **Load (WebDAV GET):** `useAppDefaults().getFileContents` (`AppWrapper.vue:185,
  353-357`); the body is assigned to `currentContent` and `serverContent`, and
  `OC-ETag` is captured (`AppWrapper.vue:356-357`). Auth is injected by the
  client service inside `useAppDefaults`.
- **Save (WebDAV PUT):** `putFileContents(currentFileContext, { content,
  previousEntityTag })` in `saveFileTask` (`AppWrapper.vue:433-441`). Exposed as
  `save()` (`AppWrapper.vue:483-485`).
- **Dirty state:** `isDirty = currentContent !== serverContent`
  (`AppWrapper.vue:164-166`).
- **Save button + filename:** rendered by `AppTopBar`
  (`AppWrapper.vue:7-16`), fed by `fileActionsSave` whose `isVisible = isEditor`
  and `isDisabled = isReadOnly || !isDirty` (`AppWrapper.vue:543-556`). The page
  title is `"<App name> for <fileName>"` (`AppWrapper.vue:204-211`).
- **`Ctrl+S`:** bound once, at the document level, via `useKeyboardActions`
  (`AppWrapper.vue:535-541`); it calls `save()` only when dirty. **Not** a global
  `keydown` listener and not per-editor.
- **Unsaved-changes guard:** `onBeforeRouteLeave` dispatches `UnsavedChangesModal`
  when dirty (`AppWrapper.vue:646-669`), plus a `beforeunload` listener
  (`AppWrapper.vue:168-179`).
- **Error toasts:** `useMessages().showErrorMessage` (`AppWrapper.vue:115, 420-427`)
  with explicit handling for `401/403`, `409/412` (conflict), and `507` (quota)
  (`AppWrapper.vue:443-479`).
- **File-size guard:** if `resource.size > meta.fileSizeLimit` a warning modal is
  shown before load with Continue/Cancel (`AppWrapper.vue:200-202, 392-414`).
- **Props passed to the wrapped component** (`slotAttrs`, `AppWrapper.vue:671-697`):
  `url, space, resource, activeFiles, isDirty, isReadOnly, applicationConfig,
  currentFileContext, currentContent, isFolderLoading`, plus handlers
  `onUpdate:resource`, `onUpdate:currentContent`, `onSave`, `onClose`.

**Conclusion:** an HTML editor app must (a) declare a `currentContent` prop and an
`update:currentContent` emit, and (b) render an editing surface. Everything else
(load, save, dirty, `Ctrl+S`, guard, toasts, size limit) is inherited for free.

---

## 1.2 App registry / association pattern

- Association is by extension via `appInfo.extensions[].extension` (1.1).
- `text/html` is **not** claimed by any existing app. The text editor handles
  `txt, md, markdown, js, json, xml, py, php, yaml, log, conf`
  (`web-app-text-editor/src/index.ts:24-70`) — **not** `html`/`htm`. No other
  `packages/web-app-*` registers `html`. So `html`/`htm` are free.
- There is no numeric "priority" field on extension registration; ordering when
  multiple apps claim an extension is handled by the apps/config layer, not by the
  editor. The text editor sets none, so the HTML editor sets none.

## 1.3 Monaco / CodeMirror availability

- **Monaco is absent** from the workspace (no `monaco-editor` in `package.json`
  files or `pnpm-lock.yaml`).
- **CodeMirror 6 is present and fully resolved** in `pnpm-lock.yaml` (transitively
  via `md-editor-v3`, a dependency of `web-pkg`): `@codemirror/state@6.5.2`,
  `@codemirror/view@6.38.2`, `@codemirror/commands@6.8.1`,
  `@codemirror/language@6.11.3`, `@codemirror/lang-html@6.4.9`,
  `@codemirror/autocomplete@6.18.7`, `codemirror@6.0.2`.
- There is **no reusable raw-CodeMirror Vue component** exported. `web-pkg`'s
  `TextEditor` (`packages/web-pkg/src/components/TextEditor/TextEditor.vue`) wraps
  `md-editor-v3`, a Markdown editor (preview/toolbars/prettier/katex/mermaid). It
  is not a plain HTML source editor.
- Theme is read from the shell via `useThemeStore().currentTheme.isDark`
  (`packages/web-pkg/src/components/TextEditor/TextEditor.vue:110, 126`).

## 1.4 CSP / iframe constraints

- This repo sets **no** CSP in source or `index.html` (`index.html:7` is only an
  `x-ua-compatible` meta). CSP is applied by the oCIS proxy from
  `dev/docker/ocis/csp.yaml`, `tests/drone/csp.yaml`,
  `deployments/examples/ocis_web/config/ocis/csp.yaml`.
- Those policies set `script-src 'self' 'unsafe-inline'` and
  `frame-src 'self' blob: <office domains>`. Inline scripts inside an iframe
  `srcdoc` run under `'unsafe-inline'`. A `srcdoc` document inherits the parent's
  CSP; with no `allow-same-origin` the iframe has an opaque origin, so the user
  HTML's *external* `'self'` script/resource loads will not match — acceptable and
  desirable for an untrusted preview.
- Iframes are already used widely **without** any sandbox attribute, e.g.
  `packages/web-app-external/src/App.vue:2-9`,
  `packages/web-pkg/src/components/Modals/SaveAsModal.vue:4-12`,
  `ExportAsPdfModal.vue`, `FilePickerModal.vue`,
  `packages/web-app-password-protected-folders/src/components/FolderViewModal.vue`.
  Our sandboxed `srcdoc` iframe is **stricter** than existing practice.
- **Nothing in this repo blocks an `<iframe srcdoc>` with
  `sandbox="allow-scripts allow-forms allow-popups"`.**

## 1.5 Build / package toolchain

- The root `vite.config.ts` **auto-discovers** apps: it globs `packages/` and
  takes any directory starting with `web-app` that has `src/index.{ts,js}` as a
  build input (`vite.config.ts:49-66`), and exposes them at runtime through
  `window.WEB_APPS_MAP` keyed by the `web-app-*` directory name
  (`vite.config.ts:329-345`). No per-app `vite.config.ts` is needed or used.
- `pnpm-workspace.yaml` globs `packages/*`, so a new `packages/web-app-html-editor`
  is automatically a workspace package.
- Which apps are **enabled** is listed by app id in `config/config.json.dist:8-15`
  and `config/config.json.sample-ocis` (`"files", "pdf-viewer", "text-editor", ...`).
  A new core editor must be added there to load by default.
- The package `name` is the bare app id, e.g. `"name": "text-editor"`
  (`web-app-text-editor/package.json:2`) — **not** `@ownclouders/...`. Runtime libs
  the shell provides are `peerDependencies` (`web-client`, `web-pkg`,
  `vue-concurrency`, `vue3-gettext`); the only `devDependency` is
  `@ownclouders/web-test-helpers` (`web-app-text-editor/package.json:10-18`).
- Tests use **vitest**: `package.json` script
  `"test:unit": "NODE_OPTIONS=--unhandled-rejections=throw vitest"` and a
  `vitest.config.ts` that merges `tests/unit/config/vitest.config` and sets a
  project `name` (`web-app-text-editor/vitest.config.ts`). Run from root via the
  workspace vitest projects, or `pnpm --filter <name> test:unit`.

## 1.6 ODS components in use

- **Buttons:** `OcButton` (kebab `oc-button`), auto-registered globally by the
  design-system plugin; props `variation` (`primary|passive|danger|...`),
  `appearance` (`filled|outline|raw|raw-inverse`), `disabled`, `@click`. Real use:
  `packages/web-app-files/src/components/AppBar/CreateAndUpload.vue:5-16`.
- **Segmented / view-mode toggle:** the idiomatic pattern is a `.oc-button-group`
  wrapper containing several `OcButton`s that toggle
  `appearance="filled"` (active) vs `"outline"` (inactive), each with an
  `OcIcon`. Real use:
  `packages/web-pkg/src/components/ViewOptions.vue:4-25` (the files view-mode
  switcher). There is no dedicated segmented-control component.
- **Notifications:** `useMessages()` from `@ownclouders/web-pkg`
  (`packages/web-pkg/src/composables/piniaStores/messages.ts:18-62`) →
  `showMessage`, `showErrorMessage`, `removeMessage`. (Used here only indirectly:
  the wrapper already surfaces load/save errors.)
- **Theme:** `useThemeStore()` from `@ownclouders/web-pkg`
  (`packages/web-pkg/src/composables/piniaStores/theme.ts`); read
  `currentTheme.isDark`.
- **CSS tokens (real names):** `--oc-color-background-default`,
  `--oc-color-background-muted`, `--oc-color-text-default`, `--oc-color-border`,
  `--oc-color-swatch-primary-default`, `--oc-font-family`, `--oc-space-*`
  (`packages/design-system/src/tokens/ods/{color,font,space}.yaml`). The default
  `--oc-font-family` upstream is `Inter, sans-serif`; brand fonts/colors come from
  the deployment theme, not hardcoded in an app.
