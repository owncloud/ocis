// Copyright 2011 Andy Balholm. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A bridge between ICAP and HTTP.
// It allows answering a REQMOD request with an HTTP response generated locally.

package icap

import (
	"log"
	"net/http"
	"time"
)

type bridgedRespWriter struct {
	irw         ResponseWriter // the underlying icap.ResponseWriter
	header      http.Header    // the headers for the HTTP response
	wroteHeader bool           // Have the headers been written yet?
}

func (w *bridgedRespWriter) Header() http.Header {
	return w.header
}

func (w *bridgedRespWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.irw.Write(p)
}

func (w *bridgedRespWriter) WriteHeader(code int) {
	if w.wroteHeader {
		log.Print("http: multiple response.WriteHeader calls")
		return
	}

	w.wroteHeader = true

	// Default output is HTML encoded in UTF-8.
	if w.header.Get("Content-Type") == "" {
		w.header.Set("Content-Type", "text/html; charset=utf-8")
	}

	if _, ok := w.header["Date"]; !ok {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}

	resp := new(http.Response)
	resp.StatusCode = code
	resp.Header = w.header

	w.irw.WriteHeader(200, resp, true)
}

// NewBridgedResponseWriter Create an http.ResponseWriter that encapsulates its response in an ICAP response.
func NewBridgedResponseWriter(w ResponseWriter) http.ResponseWriter {
	rw := new(bridgedRespWriter)
	rw.header = make(http.Header)
	rw.irw = w

	return rw
}

// ServeLocally Pass use the local HTTP server to generate a response for an ICAP request.
func ServeLocally(w ResponseWriter, req *Request) {
	brw := NewBridgedResponseWriter(w)
	http.DefaultServeMux.ServeHTTP(brw, req.Request)
}

// ServeLocallyFromHandler ---
func ServeLocallyFromHandler(w ResponseWriter, req *Request, mux http.Handler) {
	brw := NewBridgedResponseWriter(w)
	mux.ServeHTTP(brw, req.Request)
}
