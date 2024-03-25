package connector

import (
	"encoding/json"
	"errors"
	"net/http"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/rs/zerolog"
)

const (
	HeaderWopiLock    string = "X-WOPI-Lock"
	HeaderWopiOldLock string = "X-WOPI-OldLock"
)

// HttpAdapter will adapt the responses from the connector to HTTP.
//
// The adapter will use the request's context for the connector operations,
// this means that the request MUST have a valid WOPI context and a
// pre-configured logger. This should have been prepared in the routing.
//
// All operations are expected to follow the definitions found in
// https://learn.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/endpoints
type HttpAdapter struct {
	con *Connector
}

func NewHttpAdapter(gwc gatewayv1beta1.GatewayAPIClient, cfg *config.Config) *HttpAdapter {
	return &HttpAdapter{
		con: NewConnector(
			NewFileConnector(gwc, cfg),
			NewContentConnector(gwc, cfg),
		),
	}
}

func (h *HttpAdapter) GetLock(w http.ResponseWriter, r *http.Request) {
	fileCon := h.con.GetFileConnector()

	lockID, err := fileCon.GetLock(r.Context())
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set(HeaderWopiLock, lockID)
}

func (h *HttpAdapter) Lock(w http.ResponseWriter, r *http.Request) {
	oldLockID := r.Header.Get(HeaderWopiOldLock)
	lockID := r.Header.Get(HeaderWopiLock)

	fileCon := h.con.GetFileConnector()
	newLockID, err := fileCon.Lock(r.Context(), lockID, oldLockID)
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			if conError.HttpCodeOut == 409 {
				w.Header().Set(HeaderWopiLock, newLockID)
			}
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	// If no error, a HTTP 200 should be sent automatically.
	// X-WOPI-Lock header isn't needed on HTTP 200
}

func (h *HttpAdapter) RefreshLock(w http.ResponseWriter, r *http.Request) {
	lockID := r.Header.Get(HeaderWopiLock)

	fileCon := h.con.GetFileConnector()
	newLockID, err := fileCon.RefreshLock(r.Context(), lockID)
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			if conError.HttpCodeOut == 409 {
				w.Header().Set(HeaderWopiLock, newLockID)
			}
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	// If no error, a HTTP 200 should be sent automatically.
	// X-WOPI-Lock header isn't needed on HTTP 200
}

func (h *HttpAdapter) UnLock(w http.ResponseWriter, r *http.Request) {
	lockID := r.Header.Get(HeaderWopiLock)

	fileCon := h.con.GetFileConnector()
	newLockID, err := fileCon.UnLock(r.Context(), lockID)
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			if conError.HttpCodeOut == 409 {
				w.Header().Set(HeaderWopiLock, newLockID)
			}
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	// If no error, a HTTP 200 should be sent automatically.
	// X-WOPI-Lock header isn't needed on HTTP 200
}

func (h *HttpAdapter) CheckFileInfo(w http.ResponseWriter, r *http.Request) {
	fileCon := h.con.GetFileConnector()

	fileInfo, err := fileCon.CheckFileInfo(r.Context())
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	logger := zerolog.Ctx(r.Context())

	jsonFileInfo, err := json.Marshal(fileInfo)
	if err != nil {
		logger.Error().Err(err).Msg("CheckFileInfo: failed to marshal fileinfo")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, err := w.Write(jsonFileInfo)

	if err != nil {
		logger.Error().
			Err(err).
			Int("TotalBytes", len(jsonFileInfo)).
			Int("WrittenBytes", bytes).
			Msg("CheckFileInfo: failed to write contents in the HTTP response")
	}
}

func (h *HttpAdapter) GetFile(w http.ResponseWriter, r *http.Request) {
	contentCon := h.con.GetContentConnector()
	err := contentCon.GetFile(r.Context(), w)
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func (h *HttpAdapter) PutFile(w http.ResponseWriter, r *http.Request) {
	lockID := r.Header.Get(HeaderWopiLock)

	contentCon := h.con.GetContentConnector()
	newLockID, err := contentCon.PutFile(r.Context(), r.Body, r.ContentLength, lockID)
	if err != nil {
		var conError *ConnectorError
		if errors.As(err, &conError) {
			if conError.HttpCodeOut == 409 {
				w.Header().Set(HeaderWopiLock, newLockID)
			}
			http.Error(w, http.StatusText(conError.HttpCodeOut), conError.HttpCodeOut)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	// If no error, a HTTP 200 should be sent automatically.
	// X-WOPI-Lock header isn't needed on HTTP 200
}
