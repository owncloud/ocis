package connector

// ConnectorError defines an error in the connector. It contains an error code
// and a message.
// For convenience, the error code can be used as HTTP error code, although
// the connector shouldn't know anything about HTTP.
type ConnectorError struct {
	HttpCodeOut int
	Msg         string
}

func (e *ConnectorError) Error() string {
	return e.Msg
}

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

func NewConnector(fc FileConnectorService, cc ContentConnectorService) *Connector {
	return &Connector{
		fileConnector:    fc,
		contentConnector: cc,
	}
}

func (c *Connector) GetFileConnector() FileConnectorService {
	return c.fileConnector
}

func (c *Connector) GetContentConnector() ContentConnectorService {
	return c.contentConnector
}
