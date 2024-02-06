package app

import (
	"encoding/json"
	"net/http"
	"path"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/google/uuid"
)

func WopiInfoHandler(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	// Logs for this endpoint will be covered by the access log. We can't extract
	// more info
	http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
}

// CheckFileInfo returns information about the requested file and capabilities of the wopi server
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/checkfileinfo
func CheckFileInfo(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	statRes, err := app.gwc.Stat(ctx, &providerv1beta1.StatRequest{
		Ref: &wopiContext.FileReference,
	})
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("CheckFileInfo: stat failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if statRes.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", statRes.Status.Code.String()).
			Str("StatusMsg", statRes.Status.Message).
			Msg("CheckFileInfo: stat failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fileInfo := FileInfo{
		OwnerID:           statRes.Info.Owner.OpaqueId + "@" + statRes.Info.Owner.Idp,
		Size:              int64(statRes.Info.Size),
		Version:           statRes.Info.Mtime.String(),
		BaseFileName:      path.Base(statRes.Info.Path),
		BreadcrumbDocName: path.Base(statRes.Info.Path),
		// to get the folder we actually need to do a GetPath() request
		//BreadcrumbFolderName: path.Dir(statRes.Info.Path),

		UserCanNotWriteRelative: true,

		HostViewUrl: wopiContext.ViewAppUrl,
		HostEditUrl: wopiContext.EditAppUrl,

		EnableOwnerTermination: true,

		SupportsExtendedLockLength: true,

		SupportsGetLock: true,
		SupportsLocks:   true,
	}

	switch wopiContext.ViewMode {
	case appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE:
		fileInfo.SupportsUpdate = true
		fileInfo.UserCanWrite = true

	case appproviderv1beta1.ViewMode_VIEW_MODE_READ_ONLY:
		// nothing special to do here for now

	case appproviderv1beta1.ViewMode_VIEW_MODE_VIEW_ONLY:
		fileInfo.DisableExport = true
		fileInfo.DisableCopy = true
		fileInfo.DisablePrint = true
	}

	// user logic from reva wopi driver #TODO: refactor
	var isPublicShare bool = false
	if wopiContext.User != nil {
		if wopiContext.User.Id.Type == userv1beta1.UserType_USER_TYPE_LIGHTWEIGHT {
			fileInfo.UserID = statRes.Info.Owner.OpaqueId + "@" + statRes.Info.Owner.Idp
		} else {
			fileInfo.UserID = wopiContext.User.Id.OpaqueId + "@" + wopiContext.User.Id.Idp
		}

		if wopiContext.User.Opaque != nil {
			if _, ok := wopiContext.User.Opaque.Map["public-share-role"]; ok {
				isPublicShare = true
			}
		}
		if !isPublicShare {
			fileInfo.UserFriendlyName = wopiContext.User.Username
			fileInfo.UserID = wopiContext.User.Id.OpaqueId + "@" + wopiContext.User.Id.Idp
		}
	}
	if wopiContext.User == nil || isPublicShare {
		randomID, _ := uuid.NewUUID()
		fileInfo.UserID = "guest-" + randomID.String()
		fileInfo.UserFriendlyName = "Guest " + randomID.String()
		fileInfo.IsAnonymousUser = true
	}

	jsonFileInfo, err := json.Marshal(fileInfo)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("CheckFileInfo: failed to marshal fileinfo")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Msg("CheckFileInfo: success")

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonFileInfo)
	w.WriteHeader(http.StatusOK)
}
