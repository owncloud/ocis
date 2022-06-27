/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kcc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eternnoir/gncp"
)

const (
	soapUserAgent = "kcc-go-fakesoap"
	soapHeader    = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xop="http://www.w3.org/2004/08/xop/include" xmlns:xmlmime="http://www.w3.org/2004/11/xmlmime" xmlns:ns="urn:zarafa"><SOAP-ENV:Body SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">`
	soapFooter = `</SOAP-ENV:Body></SOAP-ENV:Envelope>`
)

func soapEnvelope(payload *string) *bytes.Buffer {
	var b bytes.Buffer
	b.WriteString(soapHeader)
	b.WriteString(*payload)
	b.WriteString(soapFooter)

	if debug {
		raw, _ := ioutil.ReadAll(&b)
		b = *bytes.NewBuffer(raw)
		fmt.Printf("SOAP --- request start ---\n%s\nSOAP --- request end  ---\n", string(raw))
	}
	return &b
}

func newSOAPRequest(ctx context.Context, url string, payload *string) (*http.Request, error) {
	body := soapEnvelope(payload)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("User-Agent", soapUserAgent+"/"+Version)

	return req, nil
}

func debugRawResponse(code int, data io.Reader) (io.Reader, error) {
	raw, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SOAP --- response %d start ---\n%s\nSOAP --- response end  ---\n", code, string(raw))

	return bytes.NewBuffer(raw), nil
}

func parseSOAPResponse(code int, data io.Reader, v interface{}) error {
	if debug {
		var err error
		data, err = debugRawResponse(code, data)
		if err != nil {
			return err
		}
	}

	decoder := xml.NewDecoder(data)

	match := false
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if match {
				return decoder.DecodeElement(v, &se)
			}

			if se.Name.Local == "Body" {
				match = true
			}
		}
	}

	return fmt.Errorf("failed to unmarshal SOAP response body")
}

// A SOAPClient is a network client which sends SOAP requests.
type SOAPClient interface {
	DoRequest(ctx context.Context, payload *string, v interface{}) error
}

// A SOAPClientConfig is a collection of configuration settings used when
// constructing SOAP clients.
type SOAPClientConfig struct {
	HTTPClient   *http.Client
	SocketDialer *net.Dialer
}

// DefaultSOAPClientConfig is the default SOAP client config which is used when
// constructing SOAP clients with default settings.
var DefaultSOAPClientConfig = &SOAPClientConfig{}

// A SOAPHTTPClient implements a SOAP client using the HTTP protocol.
type SOAPHTTPClient struct {
	Client *http.Client
	URI    string
}

// A SOAPSocketClient implements a SOAP client connecting to a unix socket.
type SOAPSocketClient struct {
	Dialer *net.Dialer
	Pool   gncp.ConnPool
	Path   string
}

// NewSOAPClient creates a new SOAP client for the protocol matching the
// provided URL using default connection settings. If the protocol is
// unsupported, an error is returned.
func NewSOAPClient(uri *url.URL) (SOAPClient, error) {
	if uri == nil {
		uri, _ = url.Parse(DefaultURI)
	}
	switch uri.Scheme {
	case "https":
		fallthrough
	case "http":
		return NewSOAPHTTPClient(uri, nil)

	case "file":
		return NewSOAPSocketClient(uri, nil)

	default:
		return nil, fmt.Errorf("invalid scheme '%v' for SOAP client", uri.Scheme)
	}
}

// NewSOAPClientWithConfig create new SOAP client for the protocol matching
// the provided URL using defaulft uri and config if nil is providedl. If the
// protocol is unsupported, an error is returned.
func NewSOAPClientWithConfig(uri *url.URL, config *SOAPClientConfig) (SOAPClient, error) {
	if uri == nil {
		uri, _ = url.Parse(DefaultURI)
	}
	if config == nil {
		config = DefaultSOAPClientConfig
	}
	switch uri.Scheme {
	case "https":
		fallthrough
	case "http":
		return NewSOAPHTTPClient(uri, config.HTTPClient)

	case "file":
		return NewSOAPSocketClient(uri, config.SocketDialer)

	default:
		return nil, fmt.Errorf("invalid scheme '%v' for SOAP client", uri.Scheme)
	}
}

// NewSOAPHTTPClient creates a new SOAP HTTP client for the protocol matching the
// provided URL. A http.Client can be provided to further customize the behavior
// of the client instead of using the defaults. If the protocol is unsupported,
// an error is returned.
func NewSOAPHTTPClient(uri *url.URL, client *http.Client) (*SOAPHTTPClient, error) {
	var err error

	if uri == nil {
		uri, err = uri.Parse(DefaultURI)
		if err != nil {
			return nil, err
		}
	}

	if client == nil {
		client = DefaultHTTPClient
	}

	switch uri.Scheme {
	case "https":
		fallthrough
	case "http":
		c := &SOAPHTTPClient{
			Client: client,
			URI:    uri.String(),
		}
		return c, nil
	default:
		return nil, fmt.Errorf("invalid scheme '%v' for SOAP HTTP client", uri.Scheme)
	}
}

// NewSOAPSocketClient creates a new SOAP socket client for the protocol
// matching the provided URL. A net.Dialer can be provided to further customize
// the behavior of the client instead of using the defaults. If the protocol is
//  unsupported, an error is returned.
func NewSOAPSocketClient(uri *url.URL, dialer *net.Dialer) (*SOAPSocketClient, error) {
	var err error

	if uri == nil {
		uri, err = uri.Parse(DefaultURI)
		if err != nil {
			return nil, err
		}
	}

	if dialer == nil {
		dialer = DefaultUnixDialer
	}

	if uri.Scheme != "file" {
		return nil, fmt.Errorf("invalid scheme '%v' for SOAP socket client", uri.Scheme)
	}

	c := &SOAPSocketClient{
		Dialer: dialer,
		Path:   uri.Path,
	}

	pool, err := gncp.NewPool(0, DefaultUnixMaxConnections, c.connect)
	if err != nil {
		return nil, err
	}
	c.Pool = pool

	return c, nil
}

// DoRequest sends the provided payload data as SOAP through the means of the
// accociated client. Connections are automatically reused according to keep-alive
// configuration provided by the http.Client attached to the SOAPHTTPClient.
func (sc *SOAPHTTPClient) DoRequest(ctx context.Context, payload *string, v interface{}) error {
	body := soapEnvelope(payload)

	req, err := http.NewRequest(http.MethodPost, sc.URI, body)
	if err != nil {
		return err
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("User-Agent", soapUserAgent+"/"+Version)

	resp, err := sc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		debugRawResponse(resp.StatusCode, resp.Body)
		return fmt.Errorf("unexpected http response status: %v", resp.StatusCode)
	}

	return parseSOAPResponse(resp.StatusCode, resp.Body, v)
}

func (sc *SOAPHTTPClient) String() string {
	return fmt.Sprintf("<http:%s>", sc.URI)
}

// DoRequest sends the provided payload data as SOAP through the means of the
// accociated client.
func (sc *SOAPSocketClient) DoRequest(ctx context.Context, payload *string, v interface{}) error {
	for {
		// TODO(longsleep): Use a pool which allows to add additional connections
		// in burst situations. With this current implementation based on Go
		// channel select, requests can timeout on burst situations where
		// constantly more requests than pooled connections are available come
		// in as Go's select is non-deterministic.
		c, err := sc.Pool.GetWithTimeout(sc.Dialer.Timeout)
		if err != nil {
			return fmt.Errorf("failed to open unix socket: %v", err)
		}

		body := soapEnvelope(payload)

		r := bufio.NewReader(c)

		c.SetWriteDeadline(time.Now().Add(sc.Dialer.Timeout))
		_, err = body.WriteTo(c)
		if err != nil {
			// Remove from pool and retry on any write error. This will retry
			// until the pool is not able to return a socket connection fast
			// enough anymore.
			sc.Pool.Remove(c)
			continue
		}

		// NOTE(longsleep): Kopano SOAP socket return HTTP protocol data.
		c.SetReadDeadline(time.Now().Add(sc.Dialer.Timeout))
		resp, err := http.ReadResponse(r, nil)
		if err != nil {
			sc.Pool.Remove(c)
			return fmt.Errorf("failed to read from unix socket: %v", err)
		}

		canReuseConnection := resp.Header.Get("Connection") == "keep-alive"
		defer func() {
			resp.Body.Close()
			if canReuseConnection {
				// Close makes the connection available to the pool again.
				c.Close()
			} else {
				sc.Pool.Remove(c)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			debugRawResponse(resp.StatusCode, resp.Body)
			return fmt.Errorf("unexpected http response status: %v", resp.StatusCode)
		}

		return parseSOAPResponse(resp.StatusCode, resp.Body, v)
	}
}

func (sc *SOAPSocketClient) connect() (net.Conn, error) {
	return sc.Dialer.Dial("unix", sc.Path)
}

func (sc *SOAPSocketClient) String() string {
	return fmt.Sprintf("<socket:%s>", sc.Path)
}

type xmlCharData []byte

func (s xmlCharData) String() string {
	return string(s)
}

func (s xmlCharData) Escape() string {
	var b strings.Builder

	xml.EscapeText(&b, s)
	return b.String()
}

func (s xmlCharData) WriteTo(w io.Writer) error {
	return xml.EscapeText(w, s)
}
