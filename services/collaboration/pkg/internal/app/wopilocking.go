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
}

// Lock returns a WOPI lock or performs an unlock and relock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/lock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlockandrelock
func Lock(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	oldLockID := r.Header.Get(HeaderWopiOldLock)
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

	var setOrRefreshStatus *rpcv1beta1.Status
	if oldLockID == "" {
		// If the oldLockID is empty, this is a "LOCK" request
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

		resp, err := app.gwc.SetLock(ctx, req)
		if err != nil {
			app.Logger.Error().
				Err(err).
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Msg("SetLock failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		setOrRefreshStatus = resp.Status
	} else {
		// If the oldLockID isn't empty, this is a "UnlockAndRelock" request. We'll
		// do a "RefreshLock" in reva and provide the old lock
		req := &providerv1beta1.RefreshLockRequest{
			Ref: &wopiContext.FileReference,
			Lock: &providerv1beta1.Lock{
				LockId:  lockID,
				AppName: app.Config.App.LockName,
				Type:    providerv1beta1.LockType_LOCK_TYPE_WRITE,
				Expiration: &typesv1beta1.Timestamp{
					Seconds: uint64(time.Now().Add(lockDuration).Unix()),
				},
			},
			ExistingLockId: oldLockID,
		}

		resp, err := app.gwc.RefreshLock(ctx, req)
		if err != nil {
			app.Logger.Error().
				Err(err).
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("RequestedOldLockID", oldLockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Msg("UnlockAndRefresh failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		setOrRefreshStatus = resp.Status
	}

	// we're checking the status of either the "SetLock" or "RefreshLock" operations
	switch setOrRefreshStatus.Code {
	case rpcv1beta1.Code_CODE_OK:
		app.Logger.Debug().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("SetLock successful")
		return

	case rpcv1beta1.Code_CODE_FAILED_PRECONDITION, rpcv1beta1.Code_CODE_ABORTED:
		// Code_CODE_FAILED_PRECONDITION -> Lock operation mismatched lock
		// Code_CODE_ABORTED -> UnlockAndRelock operation mismatched lock
		// In both cases, we need to get the current lock to return it in a
		// 409 response if needed
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := app.gwc.GetLock(ctx, req)
		if err != nil {
			app.Logger.Error().
				Err(err).
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Msg("SetLock failed, fallback to GetLock failed too")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("SetLock failed, fallback to GetLock failed with unexpected status")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if resp.Lock != nil {
			if resp.Lock.LockId != lockID {
				// lockId is different -> return 409 with the current lockId
				app.Logger.Warn().
					Str("FileReference", wopiContext.FileReference.String()).
					Str("RequestedLockID", lockID).
					Str("ViewMode", wopiContext.ViewMode.String()).
					Str("Requester", wopiContext.User.GetId().String()).
					Str("LockID", resp.Lock.LockId).
					Msg("SetLock conflict")
				w.Header().Set(HeaderWopiLock, resp.Lock.LockId)
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			// TODO: according to the spec we need to treat this as a RefreshLock
			// There was a problem with the lock, but the file has the same lockId now.
			// This should never happen unless there are race conditions.
			// Since the lockId matches now, we'll assume success for now.
			// As said in the todo, we probably should send a "RefreshLock" request here.
			app.Logger.Warn().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("LockID", resp.Lock.LockId).
				Msg("SetLock lock refreshed instead")
			return
		}

		// TODO: Is this the right error code?
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("SetLock failed and could not refresh")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return

	default:
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", setOrRefreshStatus.Code.String()).
			Str("StatusMsg", setOrRefreshStatus.Message).
			Msg("SetLock failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// RefreshLock refreshes a provided lock for 30 minutes
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/refreshlock
func RefreshLock(app *DemoApp, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wopiContext, _ := WopiContextFromCtx(ctx)

	lockID := r.Header.Get(HeaderWopiLock)
	if lockID == "" {
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("RefreshLock failed due to empty lockID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req := &providerv1beta1.RefreshLockRequest{
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

	resp, err := app.gwc.RefreshLock(ctx, req)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("RefreshLock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	switch resp.Status.Code {
	case rpcv1beta1.Code_CODE_OK:
		app.Logger.Debug().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("RefreshLock successful")
		return

	case rpcv1beta1.Code_CODE_NOT_FOUND:
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("RefreshLock failed, file reference not found")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return

	case rpcv1beta1.Code_CODE_ABORTED:
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("RefreshLock failed, lock mismatch")

		// Either the file is unlocked or there is no lock
		// We need to return 409 with the current lock
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := app.gwc.GetLock(ctx, req)
		if err != nil {
			app.Logger.Error().
				Err(err).
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Msg("RefreshLock failed trying to get the current lock")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("RefreshLock failed, tried to get the current lock failed with unexpected status")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Lock == nil {
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("RefreshLock failed, no lock on file")
			w.Header().Set(HeaderWopiLock, "")
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		} else {
			// lock is different than the one requested, otherwise we wouldn't reached this point
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("LockID", resp.Lock.LockId).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("RefreshLock failed, lock mismatch")
			w.Header().Set(HeaderWopiLock, resp.Lock.LockId)
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	default:
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("RefreshLock failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("Unlock failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	switch resp.Status.Code {
	case rpcv1beta1.Code_CODE_OK:
		app.Logger.Debug().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("Unlock successful")
		return
	case rpcv1beta1.Code_CODE_ABORTED:
		// File isn't locked. Need to return 409 with empty lock
		app.Logger.Error().
			Err(err).
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Msg("Unlock failed, file isn't locked")
		w.Header().Set(HeaderWopiLock, "")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	case rpcv1beta1.Code_CODE_LOCKED:
		// We need to return 409 with the current lock
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := app.gwc.GetLock(ctx, req)
		if err != nil {
			app.Logger.Error().
				Err(err).
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Msg("Unlock failed trying to get the current lock")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("Unlock failed, tried to get the current lock failed with unexpected status")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.Lock == nil {
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("Unlock failed, no lock on file")
			w.Header().Set(HeaderWopiLock, "")
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		} else {
			// lock is different than the one requested, otherwise we wouldn't reached this point
			app.Logger.Error().
				Str("FileReference", wopiContext.FileReference.String()).
				Str("RequestedLockID", lockID).
				Str("ViewMode", wopiContext.ViewMode.String()).
				Str("Requester", wopiContext.User.GetId().String()).
				Str("LockID", resp.Lock.LockId).
				Str("StatusCode", resp.Status.Code.String()).
				Str("StatusMsg", resp.Status.Message).
				Msg("Unlock failed, lock mismatch")
			w.Header().Set(HeaderWopiLock, resp.Lock.LockId)
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		}
		return
	default:
		app.Logger.Error().
			Str("FileReference", wopiContext.FileReference.String()).
			Str("RequestedLockID", lockID).
			Str("ViewMode", wopiContext.ViewMode.String()).
			Str("Requester", wopiContext.User.GetId().String()).
			Str("StatusCode", resp.Status.Code.String()).
			Str("StatusMsg", resp.Status.Message).
			Msg("Unlock failed with unexpected status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
