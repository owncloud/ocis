// Defense-in-depth Content-Security-Policy injected into the sandboxed preview
// iframe so the preview is self-protecting and does not depend on the
// deployment's proxy CSP. The iframe sandbox (no `allow-same-origin`) is the
// primary control that isolates the origin; this CSP additionally denies all
// network egress (`default-src 'none'`) and form submission so a hostile
// document cannot beacon or exfiltrate. Inline scripts/styles are allowed so the
// preview still renders self-contained pages and runs their inline JS.
export const PREVIEW_CSP = [
  "default-src 'none'",
  "script-src 'unsafe-inline'",
  "style-src 'unsafe-inline'",
  // `data:`/`blob:` are local, non-network sources (a `blob:` URL is minted by the
  // sandboxed document's own script and never leaves the opaque origin), so they
  // cannot exfiltrate. They are allowed symmetrically for every media kind a
  // self-contained document may inline.
  'img-src data: blob:',
  'font-src data: blob:',
  'media-src data: blob:',
  "form-action 'none'",
  "base-uri 'none'"
].join('; ')

// Above this size the live preview is paused: rendering (and re-parsing the
// whole document into the iframe on every settled edit) of a very large or
// hostile document can hang the browser tab. The user can still opt in to render
// it once via an explicit action.
export const PREVIEW_SIZE_LIMIT = 500_000

const META = `<meta http-equiv="Content-Security-Policy" content="${PREVIEW_CSP}">`

/**
 * Prepend the preview CSP `<meta>` as the very first bytes of the document.
 *
 * A `Content-Security-Policy` delivered via `<meta>` governs every resource the
 * parser encounters *after* it, and the first such meta wins. Prepending
 * unconditionally is therefore the only placement that covers the *whole*
 * document: any element positioned before a `<head>` — or before a crafted,
 * commented-out one — would escape a policy injected "into the real head" and be
 * free to fetch/beacon. The only cost is that a document without a doctype renders
 * in quirks mode, which is an acceptable trade for a preview whose entire purpose
 * is to be sandboxed and network-isolated; correctness of the isolation outranks
 * fidelity of the rendering mode. Keeping this a pure string prepend (no parsing,
 * no positioning) also means there is no head-detection logic for hostile markup
 * to defeat.
 */
export function wrapWithPreviewCsp(html: string): string {
  return `${META}${html ?? ''}`
}

/** Whether a document is large enough that the live preview should be paused. */
export function isPreviewTooLarge(html: string): boolean {
  return (html?.length ?? 0) > PREVIEW_SIZE_LIMIT
}
