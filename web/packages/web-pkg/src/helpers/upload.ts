/**
 * The shape of an Uppy file's `data`: a `File` (has `lastModified`), a plain `Blob` (doesn't), the
 * `{ size }` stub Uppy uses for remote/companion files, or `undefined` for a ghost file.
 */
type UploadFileData = Blob | File | { size: number | null } | undefined

/**
 * A `lastModified` value is only usable as an upload mtime when it is a finite number strictly
 * greater than zero. `undefined`, `NaN`, `Infinity`, `0` and negative values all mean "no mtime" -
 * the server then falls back to stamping the upload with its own current time rather than
 * receiving an unparseable value (which it rejects with a 400).
 */
const isValidLastModified = (lastModified: number | undefined): lastModified is number =>
  Number.isFinite(lastModified) && lastModified > 0

/**
 * Reads the `lastModified` (milliseconds) off an Uppy file's `data`. A plain `Blob` (or the
 * remote-file `{ size }` stub) has no `lastModified`, only a `File` does, so this is the one place
 * that needs to know about that.
 */
const getLastModified = (data: UploadFileData): number | undefined => (data as File)?.lastModified

/**
 * Derives the tus-meta `mtime` in seconds from an Uppy file's `data.lastModified` (milliseconds),
 * or `undefined` when there is no valid value to send.
 */
export const getUploadMtimeSeconds = (data: UploadFileData): number | undefined => {
  const lastModified = getLastModified(data)
  if (!isValidLastModified(lastModified)) {
    return undefined
  }
  return lastModified / 1000
}

/**
 * Renders the `x-oc-mtime` header value - the rounded integer-seconds string - from an Uppy file's
 * `data.lastModified` (milliseconds), or `undefined` when there is no valid value to send.
 */
export const getUploadMtimeHeaderValue = (data: UploadFileData): string | undefined => {
  const seconds = getUploadMtimeSeconds(data)
  if (seconds === undefined) {
    return undefined
  }
  return seconds.toFixed(0)
}

/**
 * Builds the `mtime` tus-meta field for an Uppy file: `{ mtime: <seconds> }` when the file's
 * `data.lastModified` is valid, or `{}` (no `mtime` key at all) otherwise - ready to spread into
 * `file.meta`.
 */
export const getUploadMtimeMeta = (data: UploadFileData): { mtime?: number } => {
  const mtime = getUploadMtimeSeconds(data)
  return mtime === undefined ? {} : { mtime }
}

/**
 * Builds the `x-oc-mtime` header field for an Uppy file: `{ 'x-oc-mtime': <seconds string> }` when
 * the file's `data.lastModified` is valid, or `{}` (no header key at all) otherwise - ready to
 * spread into a headers object.
 */
export const getUploadMtimeHeader = (data: UploadFileData): { 'x-oc-mtime'?: string } => {
  const mtimeHeaderValue = getUploadMtimeHeaderValue(data)
  return mtimeHeaderValue === undefined ? {} : { 'x-oc-mtime': mtimeHeaderValue }
}
