// Package store implements the go-micro store interface
package store

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gofrs/uuid"
	olog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
)

var (
	// Name is the default name for the settings store
	Name                   = "ocis-settings"
	managerName            = "metadata"
	settingsSpaceID        = "f1bdd61a-da7c-49fc-8203-0558109d1b4f" // uuid.Must(uuid.NewV4()).String()
	rootFolderLocation     = "settings"
	bundleFolderLocation   = "settings/bundles"
	accountsFolderLocation = "settings/accounts"
	valuesFolderLocation   = "settings/values"
)

// MetadataClient is the interface to talk to metadata service
type MetadataClient interface {
	SimpleDownload(ctx context.Context, id string) ([]byte, error)
	SimpleUpload(ctx context.Context, id string, content []byte) error
	Delete(ctx context.Context, id string) error
	ReadDir(ctx context.Context, id string) ([]string, error)
	MakeDirIfNotExist(ctx context.Context, id string) error
	Init(ctx context.Context, id string) error
}

// Store interacts with the filesystem to manage settings information
type Store struct {
	Logger olog.Logger

	mdc MetadataClient
	cfg *config.Config

	l *sync.Mutex
}

// Init initialize the store once, later calls are noops
func (s *Store) Init() {
	if s.mdc != nil {
		return
	}

	s.l.Lock()
	defer s.l.Unlock()

	if s.mdc != nil {
		return
	}

	mdc := &CachedMDC{
		next:   NewMetadataClient(s.cfg.Metadata),
		cfg:    s.cfg,
		logger: s.Logger,
	}
	if err := s.initMetadataClient(mdc); err != nil {
		s.Logger.Error().Err(err).Msg("error initializing metadata client")
	}
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	s := Store{
		Logger: olog.NewLogger(
			olog.Color(cfg.Log.Color),
			olog.Pretty(cfg.Log.Pretty),
			olog.Level(cfg.Log.Level),
			olog.File(cfg.Log.File),
		),
		cfg: cfg,
		l:   &sync.Mutex{},
	}

	return &s
}

// NewMetadataClient returns the MetadataClient
func NewMetadataClient(cfg config.Metadata) MetadataClient {
	mdc, err := metadata.NewCS3Storage(cfg.GatewayAddress, cfg.StorageAddress, cfg.SystemUserID, cfg.SystemUserIDP, cfg.SystemUserAPIKey)
	if err != nil {
		log.Fatal("error connecting to mdc:", err)
	}
	return mdc

}

// we need to lazy initialize the MetadataClient because metadata service might not be ready
func (s *Store) initMetadataClient(mdc MetadataClient) error {
	ctx := context.TODO()
	err := mdc.Init(ctx, settingsSpaceID)
	if err != nil {
		return err
	}

	for _, p := range []string{
		rootFolderLocation,
		accountsFolderLocation,
		bundleFolderLocation,
		valuesFolderLocation,
	} {
		err = mdc.MakeDirIfNotExist(ctx, p)
		if err != nil {
			return err
		}
	}

	for _, p := range s.cfg.Bundles {
		b, err := json.Marshal(p)
		if err != nil {
			return err
		}
		err = mdc.SimpleUpload(ctx, bundlePath(p.Id), b)
		if err != nil {
			return err
		}
	}

	for _, p := range defaults.DefaultRoleAssignments(s.cfg) {
		accountUUID := p.AccountUuid
		roleID := p.RoleId
		err = mdc.MakeDirIfNotExist(ctx, accountPath(accountUUID))
		if err != nil {
			return err
		}

		assIDs, err := mdc.ReadDir(ctx, accountPath(accountUUID))
		if err != nil {
			return err
		}

		adminUserID := accountUUID == s.cfg.AdminUserID
		if len(assIDs) > 0 && !adminUserID {
			// There is already a role assignment for this ID, skip to the next
			continue
		}
		// for the adminUserID we need to check if the user has the admin role every time
		if adminUserID {
			err = s.userMustHaveAdminRole(accountUUID, assIDs, mdc)
			if err != nil {
				return err
			}
			continue
		}

		ass := &settingsmsg.UserRoleAssignment{
			Id:          uuid.Must(uuid.NewV4()).String(),
			AccountUuid: accountUUID,
			RoleId:      roleID,
		}

		b, err := json.Marshal(ass)
		if err != nil {
			return err
		}
		err = mdc.SimpleUpload(ctx, assignmentPath(accountUUID, ass.Id), b)
		if err != nil {
			return err
		}
	}

	s.mdc = mdc
	return nil
}

func (s *Store) userMustHaveAdminRole(accountUUID string, assIDs []string, mdc MetadataClient) error {
	ctx := context.TODO()
	var hasAdminRole bool

	// load the assignments from the store and check if the admin role is already assigned
	for _, assID := range assIDs {
		b, err := mdc.SimpleDownload(ctx, assignmentPath(accountUUID, assID))
		switch err.(type) {
		case nil:
			// continue
		case errtypes.NotFound:
			continue
		default:
			return err
		}

		a := &settingsmsg.UserRoleAssignment{}
		err = json.Unmarshal(b, a)
		if err != nil {
			return err
		}

		if a.RoleId == defaults.BundleUUIDRoleAdmin {
			hasAdminRole = true
		}
	}

	// delete old role assignment and set admin role
	if !hasAdminRole {
		err := mdc.Delete(ctx, accountPath(accountUUID))
		switch err.(type) {
		case nil:
			// continue
		case errtypes.NotFound:
			// already gone, continue
		default:
			return err
		}

		err = mdc.MakeDirIfNotExist(ctx, accountPath(accountUUID))
		if err != nil {
			return err
		}

		ass := &settingsmsg.UserRoleAssignment{
			Id:          uuid.Must(uuid.NewV4()).String(),
			AccountUuid: accountUUID,
			RoleId:      defaults.BundleUUIDRoleAdmin,
		}
		b, err := json.Marshal(ass)
		if err != nil {
			return err
		}
		return mdc.SimpleUpload(ctx, assignmentPath(accountUUID, ass.Id), b)
	}
	return nil
}

func init() {
	settings.Registry[managerName] = New
}
