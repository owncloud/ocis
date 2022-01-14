package headers

const (
	// "github.com/cs3org/reva/internal/http/services/datagateway" is internal so we redeclare it here
	// TokenTransportHeader holds the header key for the reva transfer token
	TokenTransportHeader = "X-Reva-Transfer"
	// IfModifiedSince is used to mimic/pass on caching headers when using grpc
	IfModifiedSince = "If-Modified-Since"
)
