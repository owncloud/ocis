package app

import (
	"io"
	"net/http"

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
	http.Error(w, "", http.StatusOK)
}

// PutFile uploads the file to the storage
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/putfile
func PutFile(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// read the file from the body
	defer r.Body.Close()

	// upload the file
	err := helpers.UploadFile(
		ctx,
		r.Body,
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
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("PutFile: uploading the file failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Msg("PutFile: success")
	http.Error(w, "", http.StatusOK)
}
