// Copyright 2011 Andy Balholm. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Responding to ICAP requests.

package icap

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

// ResponseWriter ---
type ResponseWriter interface {
	// Header returns the header map that will be sent by WriteHeader.
	// Changing the header after a call to WriteHeader (or Write) has
	// no effect.
	Header() http.Header

	// Write writes the data to the connection as part of an ICAP reply.
	// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK, nil)
	// before writing the data.
	Write([]byte) (int, error)

	// Write raw data to the connection.
	WriteRaw(string)

	// WriteHeader sends an ICAP response header with status code.
	// Then it sends an HTTP header if httpMessage is not nil.
	// httpMessage may be an *http.Request or an *http.Response.
	// hasBody should be true if there will be calls to Write(), generating a message body.
	WriteHeader(code int, httpMessage interface{}, hasBody bool)
}

type respWriter struct {
	conn        *conn          // information on the connection
	req         *Request       // the request that is being responded to
	header      http.Header    // the ICAP header to write for the response
	wroteHeader bool           // true if the headers have already been written
	wroteRaw    bool           // true if raw data was written to the connection
	cw          io.WriteCloser // the chunked writer used to write the body
}

func (w *respWriter) Header() http.Header {
	return w.header
}

func (w *respWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK, nil, true)
	}

	if w.cw == nil {
		return 0, errors.New("called Write() on an icap.ResponseWriter that should not have a body")
	}
	return w.cw.Write(p)
}

func (w *respWriter) WriteRaw(p string) {
	bw := w.conn.buf.Writer
	io.WriteString(bw, p)
	w.wroteRaw = true
}

func (w *respWriter) WriteHeader(code int, httpMessage interface{}, hasBody bool) {
	if w.wroteHeader {
		log.Println("Called WriteHeader twice on the same connection")
		return
	}

	// Make the HTTP header and the Encapsulated: header.
	var header []byte
	var encap string
	var err error

	switch msg := httpMessage.(type) {
	case *http.Request:
		header, err = httpRequestHeader(msg)
		if err != nil {
			break
		}
		if hasBody {
			encap = fmt.Sprintf("req-hdr=0, req-body=%d", len(header))
		} else {
			encap = fmt.Sprintf("req-hdr=0, null-body=%d", len(header))
		}

	case *http.Response:
		header, err = httpResponseHeader(msg)
		if err != nil {
			break
		}
		if hasBody {
			encap = fmt.Sprintf("res-hdr=0, res-body=%d", len(header))
		} else {
			encap = fmt.Sprintf("res-hdr=0, null-body=%d", len(header))
		}
	}

	if encap == "" {
		if hasBody {
			method := w.req.Method
			if len(method) > 3 {
				method = method[0:3]
			}
			method = strings.ToLower(method)
			encap = fmt.Sprintf("%s-body=0", method)
		} else {
			encap = "null-body=0"
		}
	}

	w.header.Set("Encapsulated", encap)
	if _, ok := w.header["Date"]; !ok {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}

	w.header.Set("Connection", "close")

	bw := w.conn.buf.Writer
	status := StatusText(code)
	if status == "" {
		status = fmt.Sprintf("status code %d", code)
	}
	fmt.Fprintf(bw, "ICAP/1.0 %d %s\r\n", code, status)
	w.header.Write(bw)
	io.WriteString(bw, "\r\n")

	if header != nil {
		bw.Write(header)
	}

	w.wroteHeader = true

	if hasBody {
		w.cw = httputil.NewChunkedWriter(w.conn.buf.Writer)
	}
}

func (w *respWriter) finishRequest() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK, nil, false)
	}

	if w.cw != nil && !w.wroteRaw {
		w.cw.Close()
		w.cw = nil
		io.WriteString(w.conn.buf, "\r\n")
	}

	w.conn.buf.Flush()
}

// httpRequestHeader returns the headers for an HTTP request
// as a slice of bytes in a form suitable for including in an ICAP message.
func httpRequestHeader(req *http.Request) (hdr []byte, err error) {
	buf := new(bytes.Buffer)

	if req.URL == nil {
		if err != nil {
			return nil, errors.New("icap: httpRequestHeader called on Request with no URL")
		}
	}

	host := req.URL.Host
	if host == "" {
		host = req.Host
	}
	req.Header.Set("Host", host)

	uri := req.URL.String()

	fmt.Fprintf(buf, "%s %s %s\r\n", valueOrDefault(req.Method, "GET"), uri, valueOrDefault(req.Proto, "HTTP/1.1"))
	req.Header.WriteSubset(buf, map[string]bool{
		"Transfer-Encoding": true,
		"Content-Length":    true,
	})
	io.WriteString(buf, "\r\n")

	return buf.Bytes(), nil
}

// httpResponseHeader returns the headers for an HTTP response
// as a slice of bytes.
func httpResponseHeader(resp *http.Response) (hdr []byte, err error) {
	buf := new(bytes.Buffer)

	// Status line
	text := resp.Status
	if text == "" {
		text = http.StatusText(resp.StatusCode)
		if text == "" {
			text = "status code " + strconv.Itoa(resp.StatusCode)
		}
	}
	proto := resp.Proto
	if proto == "" {
		proto = "HTTP/1.1"
	}
	fmt.Fprintf(buf, "%s %d %s\r\n", proto, resp.StatusCode, text)
	if _, xIcap206Exists := resp.Header["X-Icap-206"]; xIcap206Exists {
		resp.Header.Write(buf)
	} else {
		resp.Header.WriteSubset(buf, map[string]bool{
			"Transfer-Encoding": true,
			"Content-Length":    false,
		})
	}
	io.WriteString(buf, "\r\n")

	return buf.Bytes(), nil
}

// Return value if nonempty, def otherwise.
func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
