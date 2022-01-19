package net

const (
	// "github.com/cs3org/reva/internal/http/services/datagateway" is internal so we redeclare it here
	// HeaderTokenTransport holds the header key for the reva transfer token
	HeaderTokenTransport = "X-Reva-Transfer"
	// HeaderIfModifiedSince is used to mimic/pass on caching headers when using grpc
	HeaderIfModifiedSince = "If-Modified-Since"
)
