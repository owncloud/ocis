package app

import (
	"net/http"
	"time"

	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
)

const (
	// WOPI Locks generally have a lock duration of 30 minutes and will be refreshed before expiration if needed
	// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/concepts#lock
	lockDuration time.Duration = 30 * time.Minute
)

// GetLock returns a lock or an empty string if no lock exists
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/getlock
func GetLock(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	req := &providerv1beta1.GetLockRequest{
		Ref: &wopiContext.FileReference,
	}

	resp, err := app.gwc.GetLock(ctx, req)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("GetLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("GetLock failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	lockID := ""
	if resp.Lock != nil {
		lockID = resp.Lock.LockId
	}

	// log the success at debug level
	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Str("StatusCode", resp.Status.Code.String()).
		Str("LockID", lockID).
		Msg("GetLock success")

	w.Header().Set(HeaderWopiLock, lockID)
	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}

// Lock returns a WOPI lock or performs an unlock and relock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/lock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlockandrelock
func Lock(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	// TODO: handle un- and relock

	lockID := r.Header.Get(HeaderWopiLock)
	if lockID == "" {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("Lock failed due to empty lockID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := &providerv1beta1.SetLockRequest{
		Ref: &wopiContext.FileReference,
		Lock: &providerv1beta1.Lock{
			LockId:  lockID,
			AppName: app.Config.App.LockName,
			Type:    providerv1beta1.LockType_LOCK_TYPE_WRITE,
			Expiration: &typesv1beta1.Timestamp{
				Seconds: uint64(time.Now().Add(lockDuration).Unix()),
			},
		},
	}

	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Str("RequestedLockID", lockID).
		Msg("Performing SetLock")

	resp, err := app.gwc.SetLock(ctx, req)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("RequestedLockID", lockID).
			Msg("SetLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	switch resp.Status.Code {
	case rpcv1beta1.Code_CODE_OK:
		app.Logger.Debug().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("RequestedLockID", lockID).
			Msg("SetLock successful")
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return

	case rpcv1beta1.Code_CODE_FAILED_PRECONDITION:
		// already locked
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := app.gwc.GetLock(ctx, req)
		if err != nil {
			app.Logger.Error().
				Err(err).
				Str("FileReference", wopiContext.FileReference.String()).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("RequestedLockID", lockID).
				Msg("SetLock failed, fallback to GetLock failed too")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("RequestedLockID", lockID).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("SetLock failed, fallback to GetLock failed with unexpected status")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if resp.Lock != nil {
			if resp.Lock.LockId != lockID {
				app.Logger.Warn().
					Str("FileReference", wopiContext.FileReference.String()).
					Str("ViewMode", wopiContext.ViewMode.String()).
					Str("Requester", wopiContext.User.GetId().String()).
					Str("RequestedLockID", lockID).
					Str("LockID", resp.Lock.LockId).
					Msg("SetLock conflict")
				w.Header().Set(HeaderWopiLock, resp.Lock.LockId)
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			// TODO: according to the spec we need to treat this as a RefreshLock

			app.Logger.Warn().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("RequestedLockID", lockID).
				Str("LockID", resp.Lock.LockId).
				Msg("SetLock lock refreshed instead")
			http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
			return
		}

		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("RequestedLockID", lockID).
			Msg("SetLock failed and could not refresh")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return

	default:
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("RequestedLockID", lockID).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("SetLock failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// RefreshLock refreshes a provided lock for 30 minutes
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/refreshlock
func RefreshLock(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// UnLock removes a given lock from a file
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlock
func UnLock(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	lockID := r.Header.Get(HeaderWopiLock)
	if lockID == "" {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("Unlock failed due to empty lockID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := &providerv1beta1.UnlockRequest{
		Ref: &wopiContext.FileReference,
		Lock: &providerv1beta1.Lock{
			LockId:  lockID,
			AppName: app.Config.App.LockName,
		},
	}

	resp, err := app.gwc.Unlock(ctx, req)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("RequestedLockID", lockID).
			Msg("Unlock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("RequestedLockID", lockID).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("Unlock failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.Logger.Debug().
		Str("FileReference", wopiContext.FileReference.String()).
		Str("ViewMode", wopiContext.ViewMode.String()).
		Str("Requester", wopiContext.User.GetId().String()).
		Str("RequestedLockID", lockID).
		Msg("Unlock successful")
	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}
