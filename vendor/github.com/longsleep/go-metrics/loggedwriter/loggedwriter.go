package loggedwriter

import (
	"net/http"
)

// LoggedResponseWriter define http.ResponseWriter with Status
type LoggedResponseWriter interface {
	http.ResponseWriter
	Status() int
}

// loggedResponseWriter is a http.ResponseWriter with Status.
type loggedResponseWriter struct {
	http.ResponseWriter
	status int
}

// loggedResponseWriterHijacker is a http.ResponseWriter and http.Hijacker
// with Status.
type loggedResponseWriterHijacker struct {
	LoggedResponseWriter
	http.Hijacker
}

// loggedResponseWriterHijacker is a http.ResponseWriter and http.Flusher
// with Status.
type loggedResponseWriterFlusher struct {
	LoggedResponseWriter
	http.Flusher
}

// loggedResponseWriterHijacker is a http.ResponseWriter and http.Pusher
// with Status.
type loggedResponseWriterPusher struct {
	LoggedResponseWriter
	http.Pusher
}

// loggedResponseWriterHijackerFlusherPusher is a http.ResponseWriter, http.Hijacker,
// http.Flusher and http.Pusher with Status.
type loggedResponseWriterHijackerFlusherPusher struct {
	LoggedResponseWriter
	http.Hijacker
	http.Flusher
	http.Pusher
}

// loggedResponseWriterHijackerFlusher is a http.ResponseWriter, http.Hijacker,
// and http.Flusher with Status.
type loggedResponseWriterHijackerFlusher struct {
	LoggedResponseWriter
	http.Hijacker
	http.Flusher
}

// loggedResponseWriterHijackerPusher is a http.ResponseWriter, http.Hijacker,
// and http.Pusher with Status.
type loggedResponseWriterHijackerPusher struct {
	LoggedResponseWriter
	http.Hijacker
	http.Pusher
}

// loggedResponseWriterFlusherPusher is a http.ResponseWriter, http.Flusher and
// http.Pusher with Status.
type loggedResponseWriterFlusherPusher struct {
	LoggedResponseWriter
	http.Flusher
	http.Pusher
}

// NewLoggedResponseWriter wraps the provided http.ResponseWriter with Status
// preserving the support to hijack the connection if supported by the provided
// http.ResponseWriter.
func NewLoggedResponseWriter(w http.ResponseWriter) LoggedResponseWriter {
	lw := &loggedResponseWriter{ResponseWriter: w}

	hj, hj_ok := w.(http.Hijacker)
	fl, fl_ok := w.(http.Flusher)
	ps, ps_ok := w.(http.Pusher)

	if hj_ok {
		if fl_ok && ps_ok {
			return &loggedResponseWriterHijackerFlusherPusher{
				LoggedResponseWriter: lw,
				Hijacker:             hj,
				Flusher:              fl,
				Pusher:               ps,
			}
		}
		if fl_ok {
			return &loggedResponseWriterHijackerFlusher{
				LoggedResponseWriter: lw,
				Hijacker:             hj,
				Flusher:              fl,
			}
		}
		if ps_ok {
			return &loggedResponseWriterHijackerPusher{
				LoggedResponseWriter: lw,
				Hijacker:             hj,
				Pusher:               ps,
			}
		}
		return &loggedResponseWriterHijacker{LoggedResponseWriter: lw, Hijacker: hj}
	} else if fl_ok {
		if ps_ok {
			return &loggedResponseWriterFlusherPusher{
				LoggedResponseWriter: lw,
				Flusher:              fl,
				Pusher:               ps,
			}
		}
		return &loggedResponseWriterFlusher{LoggedResponseWriter: lw, Flusher: fl}
	} else if ps_ok {
		return &loggedResponseWriterPusher{LoggedResponseWriter: lw, Pusher: ps}
	}

	return lw
}

// WriteHeader sends an HTTP response header with status code.
func (w *loggedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Status returns the written HTTP response header status code or http.StatusOK
// if none was written.
func (w *loggedResponseWriter) Status() int {
	status := w.status
	if status == 0 {
		status = http.StatusOK
	}
	return status
}
