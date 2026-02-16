package icapclient

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"strings"
)

// Client represents the icap client who makes the icap server calls
type Client struct {
	conn Conn
}

// NewClient creates a new icap client
func NewClient(options ...ConfigOption) (Client, error) {
	config := DefaultConfig()
	for _, option := range options {
		option(&config)
	}

	conn, err := NewICAPConn(config.ICAPConn)
	if err != nil {
		return Client{}, err
	}

	return Client{
		conn: conn,
	}, nil
}

// Do is the main function of the client that makes the ICAP request
func (c *Client) Do(req Request) (res Response, err error) {
	// establish connection to the icap server
	err = c.conn.Connect(req.ctx, req.URL.Host)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		err = errors.Join(err, c.conn.Close())
	}()

	req.setDefaultRequestHeaders()

	// convert the request to icap message
	message, err := toICAPRequest(req)
	if err != nil {
		return Response{}, err
	}

	// send the icap message to the server
	dataRes, err := c.conn.Send(message)
	if err != nil {
		return Response{}, err
	}

	res, err = toClientResponse(bufio.NewReader(strings.NewReader(string(dataRes))))
	if err != nil {
		return Response{}, err
	}

	// check if the message is fully done scanning or if it needs to be sent another chunk
	done := !(res.StatusCode == http.StatusContinue && !req.bodyFittedInPreview && req.previewSet)
	if done {
		return res, nil
	}

	// get the remaining body bytes
	data := req.remainingPreviewBytes
	if !bodyIsChunked(string(data)) {
		data = []byte(addHexBodyByteNotations(string(data)))
	}

	// hydrate the icap message with closing doubleCRLF suffix
	if !bytes.HasSuffix(data, []byte(doubleCRLF)) {
		data = append(data, []byte(crlf)...)
	}

	// send the remaining body bytes to the server
	dataRes, err = c.conn.Send(data)
	if err != nil {
		return Response{}, err
	}

	return toClientResponse(bufio.NewReader(strings.NewReader(string(dataRes))))
}
