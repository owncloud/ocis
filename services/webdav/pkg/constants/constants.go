package constants

type contextKey int

const (
	ContextKeyID contextKey = iota
	ContextKeyPath

	// RFC1123 time that mimics oc10. time.RFC1123 would end in "UTC", see https://github.com/golang/go/issues/13781
	// Format copied from internal package https://github.com/cs3org/reva/blob/edge/internal/http/services/owncloud/ocdav/net/net.go.
	// It's needed to match the times shown in PROPFIND and REPORT requests.
	RFC1123 = "Mon, 02 Jan 2006 15:04:05 GMT"
)
