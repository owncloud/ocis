// Above this size the live preview is paused: mermaid's layout pass re-parses and
// re-lays-out the whole diagram on every settled edit, and a very large or
// pathological source can hang the browser tab. The user can still opt in to
// render it once via an explicit action.
export const PREVIEW_SIZE_LIMIT = 500_000

/** Whether a diagram source is large enough that the live preview should be paused. */
export function isPreviewTooLarge(content: string): boolean {
  return (content?.length ?? 0) > PREVIEW_SIZE_LIMIT
}
