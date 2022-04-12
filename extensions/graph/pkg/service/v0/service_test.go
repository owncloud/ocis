package svc

import (
	"net/http"
	"testing"
)

func TestParsePurgeHeader(t *testing.T) {
	tests := map[string]bool{
		"":         false,
		"f":        false,
		"F":        false,
		"anything": false,
		"t":        true,
		"T":        true,
	}

	for input, expected := range tests {
		h := make(http.Header)
		h.Add(HeaderPurge, input)

		if expected != parsePurgeHeader(h) {
			t.Errorf("parsePurgeHeader with input %s got %t expected %t", input, !expected, expected)
		}
	}

	if h := make(http.Header); parsePurgeHeader(h) {
		t.Error("parsePurgeHeader without Purge header set got true expected false")
	}
}
