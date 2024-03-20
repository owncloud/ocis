package connector

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
