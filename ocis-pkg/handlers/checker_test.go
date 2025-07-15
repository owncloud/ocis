package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/test-go/testify/require"
	"github.com/tidwall/gjson"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

func TestCheckHandlerConfiguration(t *testing.T) {
	nopCheckCounter := 0
	nopCheck := func(_ context.Context) error { nopCheckCounter++; return nil }
	handlerConfiguration := handlers.NewCheckHandlerConfiguration().WithCheck("check-1", nopCheck)

	t.Run("add check", func(t *testing.T) {
		localCounter := 0
		handlers.NewCheckHandler(handlerConfiguration).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
		require.Equal(t, 1, nopCheckCounter)

		handlers.NewCheckHandler(handlerConfiguration.WithCheck("check-2", func(_ context.Context) error { localCounter++; return nil })).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
		require.Equal(t, 2, nopCheckCounter)
		require.Equal(t, 1, localCounter)
	})

	t.Run("checks are unique", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("checks should be unique")
			}
		}()

		handlerConfiguration.WithCheck("check-1", nopCheck)
		require.Equal(t, 3, nopCheckCounter)
	})
}

func TestCheckHandler(t *testing.T) {
	checkFactory := func(err error) func(ctx context.Context) error {
		return func(ctx context.Context) error {
			if err != nil {
				return err
			}

			<-ctx.Done()
			return nil
		}
	}

	t.Run("passes with custom status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := handlers.NewCheckHandler(
			handlers.
				NewCheckHandlerConfiguration().
				WithStatusSuccess(http.StatusCreated),
		)

		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
		require.Equal(t, http.StatusCreated, rec.Code)
		require.Equal(t, http.StatusText(http.StatusCreated), rec.Body.String())
	})

	t.Run("is not ok if any check fails", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := handlers.NewCheckHandler(
			handlers.
				NewCheckHandlerConfiguration().
				WithCheck("check-1", checkFactory(errors.New("failed"))),
		)
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, http.StatusText(http.StatusInternalServerError), rec.Body.String())
	})

	t.Run("fails with custom status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := handlers.NewCheckHandler(
			handlers.
				NewCheckHandlerConfiguration().
				WithCheck("check-1", checkFactory(errors.New("failed"))).
				WithStatusFailed(http.StatusTeapot),
		)
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
		require.Equal(t, http.StatusTeapot, rec.Code)
		require.Equal(t, http.StatusText(http.StatusTeapot), rec.Body.String())
	})

	t.Run("exits all other running tests on failure", func(t *testing.T) {
		var errs []error
		rec := httptest.NewRecorder()
		buffer := &bytes.Buffer{}
		logger := log.Logger{Logger: log.NewLogger().Output(buffer)}
		handler := handlers.NewCheckHandler(
			handlers.
				NewCheckHandlerConfiguration().
				WithLogger(logger).
				WithCheck("check-1", func(ctx context.Context) error {
					err := checkFactory(nil)(ctx)
					errs = append(errs, err)
					return err
				}).
				WithCheck("check-2", func(ctx context.Context) error {
					err := checkFactory(errors.New("failed"))(ctx)
					errs = append(errs, err)
					return err
				}).
				WithCheck("check-3", func(ctx context.Context) error {
					err := checkFactory(nil)(ctx)
					errs = append(errs, err)
					return err
				}),
		)
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

		require.Equal(t, "'check-2': failed", gjson.Get(buffer.String(), "error").String())
		require.Equal(t, 1, len(slices.DeleteFunc(errs, func(err error) bool { return err == nil })))
		require.Equal(t, 2, len(slices.DeleteFunc(errs, func(err error) bool { return err != nil })))
	})
}
