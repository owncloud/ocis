package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type dummyHandler struct{}

func (h dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestToken(t *testing.T) {
	dh := dummyHandler{}
	handler := Token("test-api-key")(dh)

	req, err := http.NewRequest("GET", "/token-protected", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	req, err = http.NewRequest("GET", "/token-protected", nil)
	req.Header.Set("Authorization", "Bearer test-api-key")
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
