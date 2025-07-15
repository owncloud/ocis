// Copyright 2011 Andy Balholm. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Reading and parsing of ICAP requests.

// Package icap provides an extensible ICAP server.
package icap

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

type badStringError struct {
	what string
	str  string
}

func (e *badStringError) Error() string { return fmt.Sprintf("%s %q", e.what, e.str) }

// A Request represents a parsed ICAP request.
type Request struct {
	Method     string               // REQMOD, RESPMOD, OPTIONS, etc.
	RawURL     string               // The URL given in the request.
	URL        *url.URL             // Parsed URL.
	Proto      string               // The protocol version.
	Header     textproto.MIMEHeader // The ICAP header
	RemoteAddr string               // the address of the computer sending the request
	Preview    []byte               // the body data for an ICAP preview

	// The HTTP messages.
	Request  *http.Request
	Response *http.Response
}

// ReadRequest reads and parses a request from b.
func ReadRequest(b *bufio.ReadWriter) (req *Request, err error) {
	tp := textproto.NewReader(b.Reader)
	req = new(Request)

	// Read first line.
	var s string
	s, err = tp.ReadLine()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}

	f := strings.SplitN(s, " ", 3)
	if len(f) < 3 {
		return nil, &badStringError{"malformed ICAP request", s}
	}
	req.Method, req.RawURL, req.Proto = f[0], f[1], f[2]

	req.URL, err = url.ParseRequestURI(req.RawURL)
	if err != nil {
		return nil, err
	}

	req.Header, err = tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	s = req.Header.Get("Encapsulated")
	if s == "" {
		return req, nil // No HTTP headers or body.
	}
	eList := strings.Split(s, ", ")
	var initialOffset, reqHdrLen, respHdrLen int
	var hasBody bool
	var prevKey string
	var prevValue int
	for _, item := range eList {
		eq := strings.Index(item, "=")
		if eq == -1 {
			return nil, &badStringError{"malformed Encapsulated: header", s}
		}
		key := item[:eq]
		value, err := strconv.Atoi(item[eq+1:])
		if err != nil {
			return nil, &badStringError{"malformed Encapsulated: header", s}
		}

		// Calculate the length of the previous section.
		switch prevKey {
		case "":
			initialOffset = value
		case "req-hdr":
			reqHdrLen = value - prevValue
		case "res-hdr":
			respHdrLen = value - prevValue
		case "req-body", "opt-body", "res-body", "null-body":
			return nil, fmt.Errorf("%s must be the last section", prevKey)
		}

		switch key {
		case "req-hdr", "res-hdr", "null-body":
		case "req-body", "res-body", "opt-body":
			hasBody = true
		default:
			return nil, &badStringError{"invalid key for Encapsulated: header", key}
		}

		prevValue = value
		prevKey = key
	}

	// Read the HTTP headers.
	var rawReqHdr, rawRespHdr []byte
	if initialOffset > 0 {
		junk := make([]byte, initialOffset)
		_, err = io.ReadFull(b, junk)
		if err != nil {
			return nil, err
		}
	}
	if reqHdrLen > 0 {
		rawReqHdr = make([]byte, reqHdrLen)
		_, err = io.ReadFull(b, rawReqHdr)
		if err != nil {
			return nil, err
		}
	}
	if respHdrLen > 0 {
		rawRespHdr = make([]byte, respHdrLen)
		_, err = io.ReadFull(b, rawRespHdr)
		if err != nil {
			return nil, err
		}
	}

	var bodyReader io.ReadCloser = emptyReader(0)
	if hasBody {
		if p := req.Header.Get("Preview"); p != "" {
			moreBody := true
			req.Preview, err = ioutil.ReadAll(newChunkedReader(b))
			if err != nil {
				if strings.Contains(err.Error(), "ieof") {
					// The data ended with "0; ieof", which the HTTP chunked reader doesn't understand.
					moreBody = false
					err = nil
				} else {
					return nil, err
				}
			}
			var r io.Reader = bytes.NewBuffer(req.Preview)
			if moreBody {
				r = io.MultiReader(r, &continueReader{buf: b})
			}
			bodyReader = ioutil.NopCloser(r)
		} else {
			bodyReader = ioutil.NopCloser(newChunkedReader(b))
		}
	}

	// Construct the http.Request.
	if rawReqHdr != nil {
		invalidURLEscapeFixed := false
		req.Request, err = http.ReadRequest(bufio.NewReader(bytes.NewBuffer(rawReqHdr)))
		if err != nil && strings.Contains(err.Error(), "invalid URL escape") {
			//Fix the request url
			// Convert the rawReqHdr to string
			// find the url\path start and end(sould be in the status line
			// convert the percents into %25
			// Then reparse the whole request
			rawReqHdrStr := string(rawReqHdr)
			result := strings.Split(rawReqHdrStr, "\n")
			result[0] = strings.Replace(result[0], "%", "%25", -1)
			// The next is a compromise since when adding "\r\n" it causes the request parsing to fail
			newReq := strings.Join(result, "\n")
			req.Request, err = http.ReadRequest(bufio.NewReader(bytes.NewBuffer([]byte(newReq))))
			if err != nil {
				return req, fmt.Errorf("error while parsing HTTP request: %v", err)
			}
			invalidURLEscapeFixed = true
		}
		if err != nil && !invalidURLEscapeFixed {
			return req, fmt.Errorf("error while parsing HTTP request: %v", err)
		}

		if req.Method == "REQMOD" {
			req.Request.Body = bodyReader
		} else {
			req.Request.Body = emptyReader(0)
		}
	}

	// Construct the http.Response.
	if rawRespHdr != nil {
		request := req.Request
		if request == nil {
			request, _ = http.NewRequest("GET", "/", nil)
		}
		req.Response, err = http.ReadResponse(bufio.NewReader(bytes.NewBuffer(rawRespHdr)), request)
		if err != nil {
			return req, fmt.Errorf("error while parsing HTTP response: %v", err)
		}

		if req.Method == "RESPMOD" {
			req.Response.Body = bodyReader
		} else {
			req.Response.Body = emptyReader(0)
		}
	}

	return
}

// An emptyReader is an io.ReadCloser that always returns os.EOF.
type emptyReader byte

func (emptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (emptyReader) Close() error {
	return nil
}

// A continueReader sends a "100 Continue" message the first time Read
// is called, creates a ChunkedReader, and reads from that.
type continueReader struct {
	buf *bufio.ReadWriter // the underlying connection
	cr  io.Reader         // the ChunkedReader
}

func (c *continueReader) Read(p []byte) (n int, err error) {
	if c.cr == nil {
		_, err := c.buf.WriteString("ICAP/1.0 100 Continue\r\n\r\n")
		if err != nil {
			return 0, err
		}
		err = c.buf.Flush()
		if err != nil {
			return 0, err
		}
		c.cr = newChunkedReader(c.buf)
	}

	return c.cr.Read(p)
}
