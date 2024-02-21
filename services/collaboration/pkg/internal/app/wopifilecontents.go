package app

import (
	"io"
	"net/http"
	"strconv"

	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/internal/helpers"
)

// GetFile downloads the file from the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/getfile
func GetFile(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// download the file
	resp, err := helpers.DownloadFile(
		ctx,
		&wopiContext.FileReference,
		app.gwc,
		wopiContext.AccessToken,
		app.Config.CS3Api.DataGateway.Insecure,
		app.Logger,
	)

	if err != nil || resp.StatusCode != http.StatusOK {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Int("HttpCode", resp.StatusCode).
			Msg("GetFile: downloading the file failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// read the file from the body
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("GetFile: copying the file content to the response body failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Msg("GetFile: success")
}

// PutFile uploads the file to the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/putfile
func PutFile(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// read the file from the body
	defer r.Body.Close()

	// We need a stat call on the target file in order to get both the lock
	// (if any) and the current size of the file
	statRes, err := app.gwc.Stat(ctx, &providerv1beta1.StatRequest{
		Ref: &wopiContext.FileReference,
	})
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", r.Header.Get(HeaderWopiLock)).
			Str("UploadLength", strconv.FormatInt(r.ContentLength, 10)).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("PutFile: stat failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if statRes.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", r.Header.Get(HeaderWopiLock)).
			Str("UploadLength", strconv.FormatInt(r.ContentLength, 10)).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", statRes.Status.Code.String()).
			Str("StatusMsg", statRes.Status.Message).
			Msg("PutFile: stat failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// If there is a lock and it mismatches, return 409
	if statRes.Info.Lock != nil && statRes.Info.Lock.LockId != r.Header.Get(HeaderWopiLock) {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", r.Header.Get(HeaderWopiLock)).
			Str("UploadLength", strconv.FormatInt(r.ContentLength, 10)).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("LockID", statRes.Info.Lock.LockId).
			Msg("PutFile: wrong lock")
		// onlyoffice says it's required to send the current lockId, MS doesn't say anything
		w.Header().Add(HeaderWopiLock, statRes.Info.Lock.LockId)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	// only unlocked uploads can go through if the target file is empty,
	// otherwise the X-WOPI-Lock header is required even if there is no lock on the file
	// This is part of the onlyoffice documentation (https://api.onlyoffice.com/editors/wopi/restapi/putfile)
	// Wopivalidator fails some tests if we don't also check for the X-WOPI-Lock header.
	if r.Header.Get(HeaderWopiLock) == "" && statRes.Info.Lock == nil && statRes.Info.Size > 0 {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", r.Header.Get(HeaderWopiLock)).
			Str("UploadLength", strconv.FormatInt(r.ContentLength, 10)).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("PutFile: file must be locked first")
		// onlyoffice says to send an empty string if the file is unlocked, MS doesn't say anything
		w.Header().Add(HeaderWopiLock, "")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	// upload the file
	err = helpers.UploadFile(
		ctx,
		r.Body,
		r.ContentLength,
		&wopiContext.FileReference,
		app.gwc,
		wopiContext.AccessToken,
		r.Header.Get(HeaderWopiLock),
		app.Config.CS3Api.DataGateway.Insecure,
		app.Logger,
	)

	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", r.Header.Get(HeaderWopiLock)).
			Str("UploadLength", strconv.FormatInt(r.ContentLength, 10)).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("PutFile: uploading the file failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("RequestedLockID", r.Header.Get(HeaderWopiLock)).
		Str("UploadLength", strconv.FormatInt(r.ContentLength, 10)).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Msg("PutFile: success")
}
