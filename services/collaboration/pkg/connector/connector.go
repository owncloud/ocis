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
	fileConnector    *FileConnector
	contentConnector *ContentConnector
}

func NewConnector(fc *FileConnector, cc *ContentConnector) *Connector {
	return &Connector{
		fileConnector:    fc,
		contentConnector: cc,
	}
}

func (c *Connector) GetFileConnector() *FileConnector {
	return c.fileConnector
}

func (c *Connector) GetContentConnector() *ContentConnector {
	return c.contentConnector
}
