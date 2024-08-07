package connector

// ConnectorResponse represent a response from the FileConnectorService.
// The ConnectorResponse is oriented to HTTP, so it has the Status, Headers
// and Body that the actual HTTP response should have. This includes HTTP
// errors with status 4xx and 5xx, which will also represent some error
// conditions for the FileConnectorService.
// Note that the Body is expected to be JSON-encoded outside before sending.
type ConnectorResponse struct {
	Status  int
	Headers map[string]string
	Body    interface{}
}

// ConnectorError defines an error in the connector. It contains an error code
// and a message.
// For convenience, the error code can be used as HTTP error code, although
// the connector shouldn't know anything about HTTP.
type ConnectorError struct {
	HttpCodeOut int
	Msg         string
}

// Error gets the error message
func (e *ConnectorError) Error() string {
	return e.Msg
}

// NewConnectorError creates a new connector error using the provided parameters
func NewConnectorError(code int, msg string) *ConnectorError {
	return &ConnectorError{
		HttpCodeOut: code,
		Msg:         msg,
	}
}

// ConnectorService is the interface to implement the WOPI operations. They're
// divided into multiple endpoints.
// The IFileConnector will implement the "File" endpoint
// The IContentConnector will implement the "File content" endpoint
type ConnectorService interface {
	GetFileConnector() FileConnectorService
	GetContentConnector() ContentConnectorService
}

// Connector will implement the WOPI operations.
// For convenience, the connector splits the operations based on the
// WOPI endpoints, so you'll need to get the specific connector first.
//
// Available endpoints:
// * "Files" -> GetFileConnector()
// * "File contents" -> GetContentConnector()
//
// Other endpoints aren't available for now.
type Connector struct {
	fileConnector    FileConnectorService
	contentConnector ContentConnectorService
}

// NewConnector creates a new connector
func NewConnector(fc FileConnectorService, cc ContentConnectorService) *Connector {
	return &Connector{
		fileConnector:    fc,
		contentConnector: cc,
	}
}

// GetFileConnector gets the file connector service associated to this connector
func (c *Connector) GetFileConnector() FileConnectorService {
	return c.fileConnector
}

// GetContentConnector gets the content connector service associated to this connector
func (c *Connector) GetContentConnector() ContentConnectorService {
	return c.contentConnector
}
