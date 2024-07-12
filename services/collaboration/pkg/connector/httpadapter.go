package connector

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector/utf7"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/locks"
	"github.com/rs/zerolog"
)

const (
	HeaderWopiLock          string = "X-WOPI-Lock"
	HeaderWopiOldLock       string = "X-WOPI-OldLock"
	HeaderWopiST            string = "X-WOPI-SuggestedTarget"
	HeaderWopiRT            string = "X-WOPI-RelativeTarget"
	HeaderWopiOverwriteRT   string = "X-WOPI-OverwriteRelativeTarget"
	HeaderWopiSize          string = "X-WOPI-Size"
	HeaderWopiValidRT       string = "X-WOPI-ValidRelativeTarget"
	HeaderWopiRequestedName string = "X-WOPI-RequestedName"
	HeaderContentLength     string = "Content-Length"
	HeaderContentType       string = "Content-Type"
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
	con   ConnectorService
	locks locks.LockParser
}

// NewHttpAdapter will create a new HTTP adapter. A new connector using the
// provided gateway API client and configuration will be used in the adapter
func NewHttpAdapter(gwc gatewayv1beta1.GatewayAPIClient, cfg *config.Config) *HttpAdapter {
	httpAdapter := &HttpAdapter{
		con: NewConnector(
			NewFileConnector(gwc, cfg),
			NewContentConnector(gwc, cfg),
		),
	}

	httpAdapter.locks = &locks.NoopLockParser{}
	if strings.ToLower(cfg.App.Name) == "microsoftofficeonline" {
		httpAdapter.locks = &locks.LegacyLockParser{}
	}
	return httpAdapter
}

// NewHttpAdapterWithConnector will create a new HTTP adapter that will use
// the provided connector service
func NewHttpAdapterWithConnector(con ConnectorService, l locks.LockParser) *HttpAdapter {
	return &HttpAdapter{
		con:   con,
		locks: l,
	}
}

// GetLock adapts the "GetLock" operation for WOPI.
// Only the request's context is needed in order to extract the WOPI context.
// The operation's response will be sent through the response writer and
// the headers according to the spec
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
	w.WriteHeader(http.StatusOK)
}

// Lock adapts the "Lock" and "UnlockAndRelock" operations for WOPI.
// The request's context is needed in order to extract the WOPI context. In
// addition, the "X-WOPI-Lock" and "X-WOPI-OldLock" headers might be needed"
// (check spec)
// The operation's response will be sent through the response writer and
// the headers according to the spec
func (h *HttpAdapter) Lock(w http.ResponseWriter, r *http.Request) {
	oldLockID := h.locks.ParseLock(r.Header.Get(HeaderWopiOldLock))
	lockID := h.locks.ParseLock(r.Header.Get(HeaderWopiLock))

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
	w.WriteHeader(http.StatusOK)
}

// RefreshLock adapts the "RefreshLock" operation for WOPI
// The request's context is needed in order to extract the WOPI context. In
// addition, the "X-WOPI-Lock" header is needed (check spec).
// The lock will be refreshed to last another 30 minutes. The value is
// hardcoded
// The operation's response will be sent through the response writer and
// the headers according to the spec
func (h *HttpAdapter) RefreshLock(w http.ResponseWriter, r *http.Request) {
	lockID := h.locks.ParseLock(r.Header.Get(HeaderWopiLock))

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
	w.WriteHeader(http.StatusOK)
}

// UnLock adapts the "Unlock" operation for WOPI
// The request's context is needed in order to extract the WOPI context. In
// addition, the "X-WOPI-Lock" header is needed (check spec).
// The operation's response will be sent through the response writer and
// the headers according to the spec
func (h *HttpAdapter) UnLock(w http.ResponseWriter, r *http.Request) {
	lockID := h.locks.ParseLock(r.Header.Get(HeaderWopiLock))

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
	w.WriteHeader(http.StatusOK)
}

// CheckFileInfo will retrieve the information of the file in json format
// Only the request's context is needed in order to extract the WOPI context.
// The operation's response will be sent through the response writer and
// the headers according to the spec
func (h *HttpAdapter) CheckFileInfo(w http.ResponseWriter, r *http.Request) {
	fileCon := h.con.GetFileConnector()

	w.Header().Set(HeaderContentType, "application/json")
	w.Header().Set(HeaderContentLength, "0")

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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set(HeaderContentLength, strconv.Itoa(len(jsonFileInfo)))
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

// GetFile will download the file
// Only the request's context is needed in order to extract the WOPI context.
// The file's content will be written in the response writer
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

// PutFile will upload the file
// The request's context and its body are needed (content length is also
// needed)
// The operation's response will be sent through the response writer and
// the headers according to the spec
func (h *HttpAdapter) PutFile(w http.ResponseWriter, r *http.Request) {
	lockID := h.locks.ParseLock(r.Header.Get(HeaderWopiLock))

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
	w.WriteHeader(http.StatusOK)
}

// PutRelativeFile will upload the file with a specific name. The name might be
// automatically adjusted depending on the request headers.
// Note that this method will also send a json body in the response.
// It has 2 mutually exclusive operation methods that are used based on the
// provided headers in the request.
// Note that this method won't used locks (not documented).
func (h *HttpAdapter) PutRelativeFile(w http.ResponseWriter, r *http.Request) {
	relativeTarget := r.Header.Get(HeaderWopiRT)
	suggestedTarget := r.Header.Get(HeaderWopiST)

	w.Header().Set(HeaderContentType, "application/json")
	w.Header().Set(HeaderContentLength, "0")

	if relativeTarget != "" && suggestedTarget != "" {
		// headers are mutually exclusive
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var response *PutRelativeResponse
	var headers *PutRelativeHeaders
	var putErr error
	fileCon := h.con.GetFileConnector()

	if suggestedTarget != "" {
		utf8Target, decErr := utf7.DecodeString(suggestedTarget)
		if decErr != nil || len(utf8Target) > 512 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		response, putErr = fileCon.PutRelativeFileSuggested(r.Context(), h.con.GetContentConnector(), r.Body, r.ContentLength, utf8Target)
	}

	if relativeTarget != "" {
		utf8Target, decErr := utf7.DecodeString(relativeTarget)
		if decErr != nil || len(utf8Target) > 512 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		response, headers, putErr = fileCon.PutRelativeFileRelative(r.Context(), h.con.GetContentConnector(), r.Body, r.ContentLength, utf8Target)
	}

	var conError *ConnectorError
	if putErr != nil && !errors.As(putErr, &conError) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger := zerolog.Ctx(r.Context())

	jsonFileInfo, err := json.Marshal(response)
	if err != nil {
		logger.Error().Err(err).Msg("PutRelativeFile: failed to marshal response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set(HeaderContentLength, strconv.Itoa(len(jsonFileInfo)))
	if conError != nil {
		if headers != nil {
			w.Header().Set(HeaderWopiValidRT, utf7.EncodeString(headers.ValidTarget))
			w.Header().Set(HeaderWopiLock, headers.LockID)
		}
		w.WriteHeader(conError.HttpCodeOut)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	bytes, err := w.Write(jsonFileInfo)

	if err != nil {
		logger.Error().
			Err(err).
			Int("TotalBytes", len(jsonFileInfo)).
			Int("WrittenBytes", bytes).
			Msg("PutRelativeFile: failed to write contents in the HTTP response")
	}
}

// DeleteFile will delete the provided file. If the file is locked and can't
// be deleted, a 409 conflict error will be returned with its corresponding
// lock.
func (h *HttpAdapter) DeleteFile(w http.ResponseWriter, r *http.Request) {
	lockID := r.Header.Get(HeaderWopiLock)

	fileCon := h.con.GetFileConnector()
	newLockID, err := fileCon.DeleteFile(r.Context(), lockID)
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

func (h *HttpAdapter) RenameFile(w http.ResponseWriter, r *http.Request) {
	lockID := r.Header.Get(HeaderWopiLock)
	requestedName := r.Header.Get(HeaderWopiRequestedName)

	utf8Target, decErr := utf7.DecodeString(requestedName)
	if decErr != nil || len(utf8Target) > 495 { // need space for the possible prefix and the extension
		w.Header().Set("X-WOPI-InvalidFileNameError", "Filename too long")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fileCon := h.con.GetFileConnector()
	response, newLockID, err := fileCon.RenameFile(r.Context(), lockID, utf8Target)
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

	// need to return a JSON response with the name if this is successful
	logger := zerolog.Ctx(r.Context())

	jsonFileInfo, err := json.Marshal(response)
	if err != nil {
		logger.Error().Err(err).Msg("RenameFile: failed to marshal response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set(HeaderContentLength, strconv.Itoa(len(jsonFileInfo)))
	w.WriteHeader(http.StatusOK)
	bytes, err := w.Write(jsonFileInfo)

	if err != nil {
		logger.Error().
			Err(err).
			Int("TotalBytes", len(jsonFileInfo)).
			Int("WrittenBytes", bytes).
			Msg("RenameFile: failed to write contents in the HTTP response")
	}
}
