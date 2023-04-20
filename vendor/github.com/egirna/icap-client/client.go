package icapclient

import (
	"net/http"
	"strconv"
	"time"
)

// Client represents the icap client who makes the icap server calls
type Client struct {
	scktDriver *Driver
	Timeout    time.Duration
}

// Do makes  does everything required to make a call to the ICAP server
func (c *Client) Do(req *Request) (*Response, error) {

	if c.scktDriver == nil { // create a new socket driver if one wasn't explicitly created
		port, err := strconv.Atoi(req.URL.Port())

		if err != nil {
			return nil, err
		}
		c.scktDriver = NewDriver(req.URL.Hostname(), port)
	}

	c.setDefaultTimeouts() // assinging default timeouts if not set already

	if req.ctx != nil { // connect with the given context if context is set
		if err := c.scktDriver.ConnectWithContext(*req.ctx); err != nil {
			return nil, err
		}
	} else {
		if err := c.scktDriver.Connect(); err != nil {
			return nil, err
		}
	}

	defer c.scktDriver.Close() // closing the socket connection

	req.SetDefaultRequestHeaders() // assigning default headers if not set already

	logDebug("The request headers: ")
	dumpDebug(req.Header)

	d, err := DumpRequest(req) // getting the byte representation of the ICAP request

	if err != nil {
		return nil, err
	}

	if err := c.scktDriver.Send(d); err != nil { // sending the entire TCP message of the ICAP client to the server connected
		return nil, err
	}

	resp, err := c.scktDriver.Receive() // taking the response

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusContinue && !req.bodyFittedInPreview && req.previewSet { // this block suggests that the ICAP request contained preview body bytes and whole body did not fit in the preview, so the serber responded with 100 Continue and the client is to send the remaining body bytes only
		logDebug("Making request for the rest of the remaining body bytes after preview, as received 100 Continue from the server...")
		return c.DoRemaining(req)
	}

	return resp, nil
}

// DoRemaining requests an ICAP server with the remaining body bytes which did not fit in the preview in the original request
func (c *Client) DoRemaining(req *Request) (*Response, error) {

	data := req.remainingPreviewBytes

	if !bodyAlreadyChunked(string(data)) { // if the body is not already chunke, then add the basic hexa body bytes notation
		ds := string(data)
		addHexaBodyByteNotations(&ds)
		data = []byte(ds)
	}

	if err := c.scktDriver.Send(data); err != nil {
		return nil, err
	}

	resp, err := c.scktDriver.Receive()

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetDriver sets a new socket driver with the client
func (c *Client) SetDriver(d *Driver) {
	c.scktDriver = d
}

func (c *Client) setDefaultTimeouts() {
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}

	if c.scktDriver.DialerTimeout == 0 {
		c.scktDriver.DialerTimeout = c.Timeout
	}

	if c.scktDriver.ReadTimeout == 0 {
		c.scktDriver.ReadTimeout = c.Timeout
	}

	if c.scktDriver.WriteTimeout == 0 {
		c.scktDriver.WriteTimeout = c.Timeout
	}
}
