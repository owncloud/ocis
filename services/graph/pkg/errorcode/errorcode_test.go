package errorcode_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

type customErr struct{}

func (customErr) Error() string {
	return "some error"
}

func TestRenderError(t *testing.T) {
	t.Parallel()

	t.Run("errorcode.Error value error", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		err := errorcode.New(errorcode.ItemNotFound, "test error")
		errorcode.RenderError(w, r, err)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("errorcode.Error zero value error", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		var err errorcode.Error
		errorcode.RenderError(w, r, err)
		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("custom error", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		var err customErr
		errorcode.RenderError(w, r, err)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
