package icapclient

import "time"

// the icap request methods
const (
	MethodOPTIONS = "OPTIONS"
	MethodRESPMOD = "RESPMOD"
	MethodREQMOD  = "REQMOD"
)

// the error messages
const (
	ErrInvalidScheme       = "the url scheme must be icap://"
	ErrMethodNotRegistered = "the requested method is not registered"
	ErrInvalidHost         = "the requested host is invalid"
	ErrConnectionNotOpen   = "no open connection to close"
	ErrInvalidTCPMsg       = "invalid tcp message"
	ErrREQMODWithNoReq     = "http request cannot be nil for method REQMOD"
	ErrREQMODWithResp      = "http response must be nil for method REQMOD"
	ErrRESPMODWithNoResp   = "http response cannot be nil for method RESPMOD"
)

// general constants required for the package
const (
	SchemeICAP                      = "icap"
	ICAPVersion                     = "ICAP/1.0"
	HTTPVersion                     = "HTTP/1.1"
	SchemeHTTPReq                   = "http_request"
	SchemeHTTPResp                  = "http_response"
	CRLF                            = "\r\n"
	DoubleCRLF                      = "\r\n\r\n"
	LF                              = "\n"
	bodyEndIndicator                = CRLF + "0" + CRLF
	fullBodyEndIndicatorPreviewMode = "; ieof" + DoubleCRLF
	icap100ContinueMsg              = "ICAP/1.0 100 Continue" + DoubleCRLF
	icap204NoModsMsg                = "ICAP/1.0 204 No modifications"
	defaultChunkLength              = 512
	defaultTimeout                  = 15 * time.Second
)

// Common ICAP headers
const (
	PreviewHeader          = "Preview"
	MethodsHeader          = "Methods"
	AllowHeader            = "Allow"
	EncapsulatedHeader     = "Encapsulated"
	TransferPreviewHeader  = "Transfer-Preview"
	ServiceHeader          = "Service"
	ISTagHeader            = "ISTag"
	OptBodyTypeHeader      = "Opt-body-type"
	MaxConnectionsHeader   = "Max-Connections"
	OptionsTTLHeader       = "Options-TTL"
	ServiceIDHeader        = "Service-ID"
	TransferIgnoreHeader   = "Transfer-Ignore"
	TransferCompleteHeader = "Transfer-Complete"
)
