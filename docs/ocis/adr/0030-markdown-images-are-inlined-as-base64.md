---
title: "30. Markdown Images Are Inlined As Base64, Not Stored As Files"
date: 2026-09-15T00:00:00+02:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/adr
geekdocFilePath: 0030-markdown-images-are-inlined-as-base64.md
---

* Status: draft
* Date: 2026-09-15
* Ticket: OCISDEV-137 (formerly owncloud/web#12407)

## Context and Problem Statement

Image upload in the Markdown editor does nothing — neither drag-and-drop nor the "Upload Images"
toolbar button produces a result, a spinner, or an error. This is not a regression. The feature is
deliberately suppressed by a single attribute, `no-upload-img`, on the `MdEditor` in
`web/packages/web-pkg/src/components/TextEditor/TextEditor.vue`. No upload handler was ever written.

It was suppressed because three questions had no answer (owncloud/web#12407, April 2025):

1. **No upload target.** Nothing decided *where* an uploaded image should be written.
2. **Share recipients cannot read it.** For a public link to a *single* `.md` file, the recipient has
   no grant on the parent folder, so an image stored beside the document renders broken for exactly
   the audience most likely to open the link.
3. **Rendering could not await.** `markdown-it` has no async parse rules. oCIS file URLs require an
   `Authorization` header, so a bare `<img src="/dav/…">` cannot load in the preview pane, and the
   blob URL that *would* work could not be resolved during a synchronous parse.

The upstream issue was closed on 2026-07-14 only because the `owncloud/web` repository was archived.
Nothing was resolved and no code landed, so the decision falls to this fork.

Two of those blockers have since expired. `md-editor-v3` 6.5.4 exposes
`transformImgUrl: (t: string) => string | Promise<string>`, which resolves image URLs *outside*
`markdown-it` and accepts a Promise — blocker 3 is gone for either design. Its CodeMirror
`linkShortener` extension folds long URLs in the editor pane and its match regex carries a dedicated
data-URI branch (`/data:[a-z]+\/[a-z0-9.+-]+(?:;base64)?,[a-z0-9+\/=%]+/i`), so an inlined image
collapses to `...` rather than flooding the source.

Blocker 2 has not expired, and cannot be solved on the client. A single-file public link grants
access to one file; any *separate* image file needs a grant nobody has issued.

## Decision

Uploaded images are **inlined into the Markdown document as base64 data URIs**. No image file is
written to storage, so there is no upload target and no second resource to grant access to. The
bytes travel inside the document, which means a public single-file link carries its own images with
no additional permission — blocker 2 does not arise rather than being worked around.

Because the bytes now live *in* the document, size is bounded at two levels:

* **Per image** — validated against the raw `File.size` handed to `onUploadImg`
  (`(files: Array<File>, callBack: UploadImgCallBack) => void`), because that is the number the user
  recognises from their file manager. Default 2 MB.
* **Per document** — a budget over the *encoded* size of all inlined images, counting those already
  present in the document and not merely the current batch. Default 10 MB. Base64 inflates by 4/3
  plus a short `data:…;base64,` prefix, so ≈1.34×.

Both are configurable through the editor's `applicationConfig`, following `showPreviewOnlyMd`.

An upload over either limit is **rejected with a validation message**, not silently altered. The
message is rendered in the editor footer via the existing `#defFooters` slot, and states the actual
size and the limit, after GitLab's attachment pattern. Downscaling the image to fit was rejected:
silently re-encoding a user's asset trades one surprise for another.

The two limits carry different messages, because they are different failures. Exceeding the
per-image cap is about the file; exhausting the document budget is about the document, and must
report what remains rather than what is allowed.

### Consequences

* **Documents grow, and every version keeps a full copy.** `OCIS_DISABLE_VERSIONING` defaults to
  `false` (`services/storage-users/pkg/config/config.go:188`, and `:325` for S3NG) and oCIS stores
  complete copies, not deltas. There is no retention policy — no max count, no max age, no size cap
  — so an inlined image is re-written and retained on every save of its document. This amplification
  is the price paid for the permission model, and it is the reason a document budget exists at all.
* **The document budget shadows a server-side limit the web app cannot read.** The real constraint
  is `SEARCH_CONTENT_EXTRACTION_SIZE_LIMIT`, default 20 MB, which truncates with `io.LimitReader`
  before extraction (`services/search/pkg/config/config.go:30`,
  `services/search/pkg/content/tika.go:77-82`). Base64 payloads consume the budget that the
  document's prose needs, so an over-large document has its trailing text silently dropped from the
  search index — the user is given no error and simply cannot find their own writing. The 10 MB
  default is half the extraction limit, chosen for headroom. An operator who lowers
  `SEARCH_CONTENT_EXTRACTION_SIZE_LIMIT` re-opens this hole with no signal, because the two values
  live in different services and nothing links them.
* **Neither limit is enforced server-side.** No body-size limit exists in the proxy or webdav
  configuration, and `httputil.ReverseProxy` imposes none. These caps are editor affordances that
  keep search working for documents this editor writes; they are not a control. A client-side cap is
  not a size limit, in the same sense that a hidden file action is not a permission check.
* **`![alt](data:image/png;base64,…)` in a document is deliberate.** A reader who assumes it is a bug
  and "fixes" it by extracting the images to files reintroduces blocker 2.
* **The existing height-cap CSS for long URLs is dead.** `web-app-text-editor/src/App.vue:44-49`
  targets `.toastui-editor-md-link-url` with a comment about base64 images, and
  `TextEditor.vue:240` targets `.toastui-editor-tabs`. Both are ToastUI-era leftovers matching
  nothing since the move to `md-editor-v3`; the `linkShortener` extension supersedes them.

## Considered Options

* **Write the image as a file in storage and reference it as `![alt](./image.png)`.** The
  conventional design, and it produces a far smaller document with no version amplification. There
  is even a precedent for the upload target: `useFileActionsSetImage.ts:31-34` parks a Space's cover
  image in the `.space` meta folder via `getDefaultMetaFolder` / `createDefaultMetaFolder`. Rejected
  because blocker 2 survives: a recipient of a single-file public link has no grant on any separate
  file, so shared documents show broken images. Solving that needs a server-side permission model
  for document-referenced resources, which is a much larger piece of work than this ticket.
* **Inline base64 with no limit.** Simplest to ship. Rejected: it silently degrades search on
  exactly the documents the feature is meant to enrich, and the degradation is invisible.
* **Downscale over-large images to fit instead of rejecting.** Keeps drag-and-drop of a phone photo
  working, which a 2 MB cap otherwise fails routinely. Rejected as the primary behaviour because it
  alters the user's asset without asking. Worth revisiting as an explicit, offered action.
* **Cap per image only, as GitLab does.** Rejected: GitLab stores attachments *outside* the
  document, so bytes never accumulate and per-file is the only limit it needs. Here they accumulate,
  and `onUploadImg` receives an `Array<File>` — a single drag-and-drop of five individually legal
  images can blow the document budget with every file passing validation.
* **Leave the feature hidden.** Rejected: it is the ticket, the blockers that justified hiding it
  have either expired or been designed around, and the current behaviour — nothing happens, with no
  error — is the worst available.
