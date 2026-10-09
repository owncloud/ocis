package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// A request with a reva token but no user, like reva's internal download requests to
// the data gateway, is passed on unchanged and does not log an error.
func TestCreateHomeWithoutUserPassesThroughQuietly(t *testing.T) {
	// Other tests in this package raise the global level; let error lines through.
	defer zerolog.SetGlobalLevel(zerolog.GlobalLevel())
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	var buf bytes.Buffer
	logger := log.Logger{Logger: zerolog.New(&buf).Level(zerolog.InfoLevel)}

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})
	sut := CreateHome(Logger(logger))(next)

	req := httptest.NewRequest(http.MethodGet, "/data/some-transfer-token", nil)
	req.Header.Set("x-access-token", "a-token")
	rw := httptest.NewRecorder()
	sut.ServeHTTP(rw, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Empty(t, buf.String())
}
