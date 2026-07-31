# HTML Editor (`html-editor`)

An ownCloud Web app that opens `.html` / `.htm` / `.xhtml` files for editing with a
live, sandboxed preview. No WOPI and no extra server: the file is loaded and saved
over WebDAV by the standard app wrapper, the source is edited in CodeMirror, and
the rendered result is shown in a sandboxed `<iframe>`.

## Features

- **CodeMirror 6** source editor with HTML syntax highlighting, line numbers, and
  theme that follows the active ownCloud theme (light/dark).
- **Live preview** in a strictly sandboxed iframe (`srcdoc`, `sandbox="allow-scripts"`,
  no `allow-same-origin`/`allow-forms`/`allow-popups`), with a self-contained CSP
  injected into the document. The preview is updated as you type, and is paused for
  very large files until you opt in. The preview is network-isolated, so documents
  that reference external CSS/JS/images will not load them.
- **Three view modes** from the toolbar: Split (default), Editor only, Preview
  only.
- **Save, dirty tracking, `Ctrl+S`, the unsaved-changes guard, and error
  notifications** are provided by the shared app wrapper (`AppWrapper`), exactly
  like the text editor; this app does not reimplement them.

## How it works

`src/index.ts` registers the app for the `html`, `htm`, and `xhtml` extensions via
`defineWebApplication` and routes through `AppWrapperRoute`. `AppWrapper` performs
the WebDAV `GET`, passes the file body to the app as the `currentContent` prop, and
re-saves with `PUT` whenever the user saves. Declaring the `currentContent` prop and
the `update:currentContent` emit is what enables the wrapper's save button, autosave,
`Ctrl+S` and unsaved-changes guard (see `ARCHAEOLOGY.md` and `DECISIONS.md`).

```
src/
  index.ts                     app + extension registration
  App.vue                      thin shell: props in, update:currentContent out, view-mode layout
  components/
    HtmlEditorPane.vue         CodeMirror 6 wrapper (HTML mode, themed)
    HtmlPreviewPane.vue        sandboxed iframe (srcdoc)
    HtmlToolbar.vue            view-mode toggle (Editor | Split | Preview)
```

## Enabling the app

The app is enabled by default through `config/config.json.dist` and
`config/config.json.sample-ocis` (the `apps` array contains `"html-editor"`). For a
custom deployment, add `"html-editor"` to the `apps` list of your `config.json`.

## Content-Security-Policy

The preview iframe uses `srcdoc` with a `sandbox` attribute and no
`allow-same-origin`, so it runs at an opaque origin and cannot read the shell's
cookies or storage. oCIS' default CSP (`script-src 'self' 'unsafe-inline'`,
`frame-src 'self' blob:`) permits this: inline scripts inside the preview run under
`'unsafe-inline'`, while external/`'self'` resource loads from the untrusted HTML
are blocked by the opaque origin. No CSP change is required.

## Tests

```
pnpm --filter html-editor test:unit
```
