package net

// Common HTTP headers.
const (
	HeaderAcceptRanges               = "Accept-Ranges"
	HeaderAccessControlAllowHeaders  = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	HeaderContentDisposistion        = "Content-Disposition"
	HeaderContentLength              = "Content-Length"
	HeaderContentRange               = "Content-Range"
	HeaderContentType                = "Content-Type"
	HeaderETag                       = "ETag"
	HeaderLastModified               = "Last-Modified"
	HeaderLocation                   = "Location"
	HeaderRange                      = "Range"
	HeaderIfMatch                    = "If-Match"
)

// webdav headers
const (
	HeaderDav         = "DAV"         // https://datatracker.ietf.org/doc/html/rfc4918#section-10.1
	HeaderDepth       = "Depth"       // https://datatracker.ietf.org/doc/html/rfc4918#section-10.2
	HeaderDestination = "Destination" // https://datatracker.ietf.org/doc/html/rfc4918#section-10.3
	HeaderIf          = "If"          // https://datatracker.ietf.org/doc/html/rfc4918#section-10.4
	HeaderLockToken   = "Lock-Token"  // https://datatracker.ietf.org/doc/html/rfc4918#section-10.5
	HeaderOverwrite   = "Overwrite"   // https://datatracker.ietf.org/doc/html/rfc4918#section-10.6
	HeaderTimeout     = "Timeout"     // https://datatracker.ietf.org/doc/html/rfc4918#section-10.7
)

// Non standard HTTP headers.
const (
	HeaderOCFileID             = "OC-FileId"
	HeaderOCETag               = "OC-ETag"
	HeaderOCChecksum           = "OC-Checksum"
	HeaderOCPermissions        = "OC-Perm"
	HeaderTusResumable         = "Tus-Resumable"
	HeaderTusVersion           = "Tus-Version"
	HeaderTusExtension         = "Tus-Extension"
	HeaderTusChecksumAlgorithm = "Tus-Checksum-Algorithm"
	HeaderTusUploadExpires     = "Upload-Expires"
	HeaderUploadChecksum       = "Upload-Checksum"
	HeaderUploadLength         = "Upload-Length"
	HeaderUploadMetadata       = "Upload-Metadata"
	HeaderUploadOffset         = "Upload-Offset"
	HeaderOCMtime              = "X-OC-Mtime"
	HeaderExpectedEntityLength = "X-Expected-Entity-Length"
	HeaderLitmus               = "X-Litmus"
)
