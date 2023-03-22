package handlers

import (
	"io"
	"net/http"
)

// Health can be used for a health endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// TODO: check if services are up and running

	_, err := io.WriteString(w, http.StatusText(http.StatusOK))
	// io.WriteString should not fail but if it does we want to know.
	if err != nil {
		panic(err)
	}
}

// Ready can be used as a ready endpoint
func Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// TODO: check if services are up and running

	_, err := io.WriteString(w, http.StatusText(http.StatusOK))
	// io.WriteString should not fail but if it does we want to know.
	if err != nil {
		panic(err)
	}
}
