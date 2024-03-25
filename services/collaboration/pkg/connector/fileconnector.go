package connector

import (
	"context"
	"encoding/hex"
	"path"
	"time"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/google/uuid"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/middleware"
	"github.com/rs/zerolog"
)

const (
	// WOPI Locks generally have a lock duration of 30 minutes and will be refreshed before expiration if needed
	// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/concepts#lock
	lockDuration time.Duration = 30 * time.Minute
)

type FileConnector struct {
	gwc gatewayv1beta1.GatewayAPIClient
	cfg *config.Config
}

func NewFileConnector(gwc gatewayv1beta1.GatewayAPIClient, cfg *config.Config) *FileConnector {
	return &FileConnector{
		gwc: gwc,
		cfg: cfg,
	}
}

// GetLock returns a lock or an empty string if no lock exists
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/getlock
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// The lock ID applied to the file reference in the context will be returned
// (if any). An error will be returned if something goes wrong. The error
// could be a ConnectorError
func (f *FileConnector) GetLock(ctx context.Context) (string, error) {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return "", err
	}

	logger := zerolog.Ctx(ctx)

	req := &providerv1beta1.GetLockRequest{
		Ref: &wopiContext.FileReference,
	}

	resp, err := f.gwc.GetLock(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("GetLock failed")
		return "", err
	}

	if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("GetLock failed with unexpected status")
		return "", NewConnectorError(404, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
	}

	lockID := ""
	if resp.GetLock() != nil {
		lockID = resp.GetLock().GetLockId()
	}

	// log the success at debug level
	logger.Debug().
		Str("LockID", lockID).
		Msg("GetLock success")

	return lockID, nil
}

// Lock returns a WOPI lock or performs an unlock and relock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/lock
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlockandrelock
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// Lock the file reference contained in the context with the provided lockID.
// The oldLockID is only used for the "unlock and relock" operation. The "lock"
// operation doesn't use the oldLockID and needs to be empty in this case.
//
// For the "lock" operation, if the operation is successful, an empty lock id
// will be returned without any error. In case of conflict, the current lock
// id will be returned along with a 409 ConnectorError. For any other error,
// the method will return an empty lock id.
//
// For the "unlock and relock" operation, the behavior will be the same.
func (f *FileConnector) Lock(ctx context.Context, lockID, oldLockID string) (string, error) {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return "", err
	}

	logger := zerolog.Ctx(ctx).With().
		Str("RequestedLockID", lockID).
		Str("RequestedOldLockID", oldLockID).
		Logger()

	if lockID == "" {
		logger.Error().Msg("Lock failed due to empty lockID")
		return "", NewConnectorError(400, "Requested lockID is empty")
	}

	var setOrRefreshStatus *rpcv1beta1.Status
	if oldLockID == "" {
		// If the oldLockID is empty, this is a "LOCK" request
		req := &providerv1beta1.SetLockRequest{
			Ref: &wopiContext.FileReference,
			Lock: &providerv1beta1.Lock{
				LockId:  lockID,
				AppName: f.cfg.App.LockName,
				Type:    providerv1beta1.LockType_LOCK_TYPE_WRITE,
				Expiration: &typesv1beta1.Timestamp{
					Seconds: uint64(time.Now().Add(lockDuration).Unix()),
				},
			},
		}

		resp, err := f.gwc.SetLock(ctx, req)
		if err != nil {
			logger.Error().Err(err).Msg("SetLock failed")
			return "", err
		}
		setOrRefreshStatus = resp.GetStatus()
	} else {
		// If the oldLockID isn't empty, this is a "UnlockAndRelock" request. We'll
		// do a "RefreshLock" in reva and provide the old lock
		req := &providerv1beta1.RefreshLockRequest{
			Ref: &wopiContext.FileReference,
			Lock: &providerv1beta1.Lock{
				LockId:  lockID,
				AppName: f.cfg.App.LockName,
				Type:    providerv1beta1.LockType_LOCK_TYPE_WRITE,
				Expiration: &typesv1beta1.Timestamp{
					Seconds: uint64(time.Now().Add(lockDuration).Unix()),
				},
			},
			ExistingLockId: oldLockID,
		}

		resp, err := f.gwc.RefreshLock(ctx, req)
		if err != nil {
			logger.Error().Err(err).Msg("UnlockAndRefresh failed")
			return "", err
		}
		setOrRefreshStatus = resp.GetStatus()
	}

	// we're checking the status of either the "SetLock" or "RefreshLock" operations
	switch setOrRefreshStatus.GetCode() {
	case rpcv1beta1.Code_CODE_OK:
		logger.Debug().Msg("SetLock successful")
		return "", nil

	case rpcv1beta1.Code_CODE_FAILED_PRECONDITION, rpcv1beta1.Code_CODE_ABORTED:
		// Code_CODE_FAILED_PRECONDITION -> Lock operation mismatched lock
		// Code_CODE_ABORTED -> UnlockAndRelock operation mismatched lock
		// In both cases, we need to get the current lock to return it in a
		// 409 response if needed
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := f.gwc.GetLock(ctx, req)
		if err != nil {
			logger.Error().Err(err).Msg("SetLock failed, fallback to GetLock failed too")
			return "", err
		}

		if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("SetLock failed, fallback to GetLock failed with unexpected status")
		}

		if resp.GetLock() != nil {
			if resp.GetLock().GetLockId() != lockID {
				// lockId is different -> return 409 with the current lockId
				logger.Warn().
					Str("LockID", resp.GetLock().GetLockId()).
					Msg("SetLock conflict")
				return resp.GetLock().GetLockId(), NewConnectorError(409, "Lock conflict")
			}

			// TODO: according to the spec we need to treat this as a RefreshLock
			// There was a problem with the lock, but the file has the same lockId now.
			// This should never happen unless there are race conditions.
			// Since the lockId matches now, we'll assume success for now.
			// As said in the todo, we probably should send a "RefreshLock" request here.
			logger.Warn().
				Str("LockID", resp.GetLock().GetLockId()).
				Msg("SetLock lock refreshed instead")
			return resp.GetLock().GetLockId(), nil
		}

		logger.Error().Msg("SetLock failed and could not refresh")
		return "", NewConnectorError(500, "Could not refresh the lock")

	case rpcv1beta1.Code_CODE_NOT_FOUND:
		logger.Error().Msg("SetLock failed, file not found")
		return "", NewConnectorError(404, "File not found")

	default:
		logger.Error().
			Str("StatusCode", setOrRefreshStatus.GetCode().String()).
			Str("StatusMsg", setOrRefreshStatus.GetMessage()).
			Msg("SetLock failed with unexpected status")
		return "", NewConnectorError(500, setOrRefreshStatus.GetCode().String()+" "+setOrRefreshStatus.GetMessage())
	}
}

// RefreshLock refreshes a provided lock for 30 minutes
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/refreshlock
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// If the operation is successful, an empty lock id will be returned without
// any error. In case of conflict, the current lock id will be returned
// along with a 409 ConnectorError. For any other error, the method will
// return an empty lock id.
// The conflict happens if the provided lockID doesn't match the one actually
// applied in the target file.
func (f *FileConnector) RefreshLock(ctx context.Context, lockID string) (string, error) {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return "", err
	}

	logger := zerolog.Ctx(ctx).With().
		Str("RequestedLockID", lockID).
		Logger()

	if lockID == "" {
		logger.Error().Msg("RefreshLock failed due to empty lockID")
		return "", NewConnectorError(400, "Requested lockID is empty")
	}

	req := &providerv1beta1.RefreshLockRequest{
		Ref: &wopiContext.FileReference,
		Lock: &providerv1beta1.Lock{
			LockId:  lockID,
			AppName: f.cfg.App.LockName,
			Type:    providerv1beta1.LockType_LOCK_TYPE_WRITE,
			Expiration: &typesv1beta1.Timestamp{
				Seconds: uint64(time.Now().Add(lockDuration).Unix()),
			},
		},
	}

	resp, err := f.gwc.RefreshLock(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("RefreshLock failed")
		return "", err
	}

	switch resp.GetStatus().GetCode() {
	case rpcv1beta1.Code_CODE_OK:
		logger.Debug().Msg("RefreshLock successful")
		return "", nil

	case rpcv1beta1.Code_CODE_NOT_FOUND:
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("RefreshLock failed, file reference not found")
		return "", NewConnectorError(404, "File reference not found")

	case rpcv1beta1.Code_CODE_ABORTED:
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("RefreshLock failed, lock mismatch")

		// Either the file is unlocked or there is no lock
		// We need to return 409 with the current lock
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := f.gwc.GetLock(ctx, req)
		if err != nil {
			logger.Error().Err(err).Msg("RefreshLock failed trying to get the current lock")
			return "", err
		}

		if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("RefreshLock failed, tried to get the current lock failed with unexpected status")
			return "", NewConnectorError(500, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
		}

		if resp.GetLock() == nil {
			logger.Error().
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("RefreshLock failed, no lock on file")
			return "", NewConnectorError(409, "No lock on file")
		} else {
			// lock is different than the one requested, otherwise we wouldn't reached this point
			logger.Error().
				Str("LockID", resp.GetLock().GetLockId()).
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("RefreshLock failed, lock mismatch")
			return resp.GetLock().GetLockId(), NewConnectorError(409, "Lock mismatch")
		}
	default:
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("RefreshLock failed with unexpected status")
		return "", NewConnectorError(500, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
	}
}

// UnLock removes a given lock from a file
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/unlock
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// If the operation is successful, an empty lock id will be returned without
// any error. In case of conflict, the current lock id will be returned
// along with a 409 ConnectorError. For any other error, the method will
// return an empty lock id.
// The conflict happens if the provided lockID doesn't match the one actually
// applied in the target file.
func (f *FileConnector) UnLock(ctx context.Context, lockID string) (string, error) {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return "", err
	}

	logger := zerolog.Ctx(ctx).With().
		Str("RequestedLockID", lockID).
		Logger()

	if lockID == "" {
		logger.Error().Msg("Unlock failed due to empty lockID")
		return "", NewConnectorError(400, "Requested lockID is empty")
	}

	req := &providerv1beta1.UnlockRequest{
		Ref: &wopiContext.FileReference,
		Lock: &providerv1beta1.Lock{
			LockId:  lockID,
			AppName: f.cfg.App.LockName,
		},
	}

	resp, err := f.gwc.Unlock(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("Unlock failed")
		return "", err
	}

	switch resp.GetStatus().GetCode() {
	case rpcv1beta1.Code_CODE_OK:
		logger.Debug().Msg("Unlock successful")
		return "", nil
	case rpcv1beta1.Code_CODE_ABORTED:
		// File isn't locked. Need to return 409 with empty lock
		logger.Error().Err(err).Msg("Unlock failed, file isn't locked")
		return "", NewConnectorError(409, "File is not locked")
	case rpcv1beta1.Code_CODE_LOCKED:
		// We need to return 409 with the current lock
		req := &providerv1beta1.GetLockRequest{
			Ref: &wopiContext.FileReference,
		}

		resp, err := f.gwc.GetLock(ctx, req)
		if err != nil {
			logger.Error().Err(err).Msg("Unlock failed trying to get the current lock")
			return "", err
		}

		if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
			logger.Error().
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("Unlock failed, tried to get the current lock failed with unexpected status")
			return "", NewConnectorError(500, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
		}

		var outLockId string
		if resp.GetLock() == nil {
			logger.Error().
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("Unlock failed, no lock on file")
			outLockId = ""
		} else {
			// lock is different than the one requested, otherwise we wouldn't reached this point
			logger.Error().
				Str("LockID", resp.GetLock().GetLockId()).
				Str("StatusCode", resp.GetStatus().GetCode().String()).
				Str("StatusMsg", resp.GetStatus().GetMessage()).
				Msg("Unlock failed, lock mismatch")
			outLockId = resp.GetLock().GetLockId()
		}
		return outLockId, NewConnectorError(409, "Lock mismatch")
	default:
		logger.Error().
			Str("StatusCode", resp.GetStatus().GetCode().String()).
			Str("StatusMsg", resp.GetStatus().GetMessage()).
			Msg("Unlock failed with unexpected status")
		return "", NewConnectorError(500, resp.GetStatus().GetCode().String()+" "+resp.GetStatus().GetMessage())
	}
}

// CheckFileInfo returns information about the requested file and capabilities of the wopi server
// https://docs.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/checkfileinfo
//
// The context MUST have a WOPI context, otherwise an error will be returned.
// You can pass a pre-configured zerologger instance through the context that
// will be used to log messages.
//
// If the operation is successful, a "FileInfo" instance will be returned,
// otherwise the "FileInfo" will be empty and an error will be returned.
func (f *FileConnector) CheckFileInfo(ctx context.Context) (FileInfo, error) {
	wopiContext, err := middleware.WopiContextFromCtx(ctx)
	if err != nil {
		return FileInfo{}, err
	}

	logger := zerolog.Ctx(ctx)

	statRes, err := f.gwc.Stat(ctx, &providerv1beta1.StatRequest{
		Ref: &wopiContext.FileReference,
	})
	if err != nil {
		logger.Error().Err(err).Msg("CheckFileInfo: stat failed")
		return FileInfo{}, err
	}

	if statRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		logger.Error().
			Str("StatusCode", statRes.GetStatus().GetCode().String()).
			Str("StatusMsg", statRes.GetStatus().GetMessage()).
			Msg("CheckFileInfo: stat failed with unexpected status")
		return FileInfo{}, NewConnectorError(500, statRes.GetStatus().GetCode().String()+" "+statRes.GetStatus().GetMessage())
	}

	fileInfo := FileInfo{
		// OwnerId must use only alphanumeric chars (https://learn.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/checkfileinfo/checkfileinfo-response#requirements-for-user-identity-properties)
		OwnerId:           hex.EncodeToString([]byte(statRes.GetInfo().GetOwner().GetOpaqueId() + "@" + statRes.GetInfo().GetOwner().GetIdp())),
		Size:              int64(statRes.GetInfo().GetSize()),
		Version:           statRes.GetInfo().GetMtime().String(),
		BaseFileName:      path.Base(statRes.GetInfo().GetPath()),
		BreadcrumbDocName: path.Base(statRes.GetInfo().GetPath()),
		// to get the folder we actually need to do a GetPath() request
		//BreadcrumbFolderName: path.Dir(statRes.Info.Path),

		UserCanNotWriteRelative: true,

		HostViewUrl: wopiContext.ViewAppUrl,
		HostEditUrl: wopiContext.EditAppUrl,

		//EnableOwnerTermination: true,  // enable only for collabora? wopivalidator is complaining
		EnableOwnerTermination: false,

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
		// UserId must use only alphanumeric chars (https://learn.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/files/checkfileinfo/checkfileinfo-response#requirements-for-user-identity-properties)
		if wopiContext.User.GetId().GetType() == userv1beta1.UserType_USER_TYPE_LIGHTWEIGHT {
			fileInfo.UserId = hex.EncodeToString([]byte(statRes.GetInfo().GetOwner().GetOpaqueId() + "@" + statRes.GetInfo().GetOwner().GetIdp()))
		} else {
			fileInfo.UserId = hex.EncodeToString([]byte(wopiContext.User.GetId().GetOpaqueId() + "@" + wopiContext.User.GetId().GetIdp()))
		}

		if wopiContext.User.GetOpaque() != nil {
			if _, ok := wopiContext.User.GetOpaque().GetMap()["public-share-role"]; ok {
				isPublicShare = true
			}
		}
		if !isPublicShare {
			fileInfo.UserFriendlyName = wopiContext.User.GetUsername()
			fileInfo.UserId = hex.EncodeToString([]byte(wopiContext.User.GetId().GetOpaqueId() + "@" + wopiContext.User.GetId().GetIdp()))
		}
	}
	if wopiContext.User == nil || isPublicShare {
		randomID, _ := uuid.NewUUID()
		fileInfo.UserId = hex.EncodeToString([]byte("guest-" + randomID.String()))
		fileInfo.UserFriendlyName = "Guest " + randomID.String()
		fileInfo.IsAnonymousUser = true
	}

	logger.Debug().Msg("CheckFileInfo: success")
	return fileInfo, nil
}
