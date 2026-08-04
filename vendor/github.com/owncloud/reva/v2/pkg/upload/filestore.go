package upload

import (
	"context"
	"encoding/json"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

// TokenOptions carries the JWT-signing configuration needed to produce transfer
// URLs for the postprocessing service.
type TokenOptions struct {
	DownloadEndpoint     string
	DataGatewayEndpoint  string
	TransferSharedSecret string
	TransferExpires      int64
}

// SessionStore abstracts upload-session persistence for the Coordinator.
type SessionStore interface {
	New(ctx context.Context) Session
	Get(ctx context.Context, id string) (Session, error)
	List(ctx context.Context) ([]Session, error)
}

// FileStore is a filesystem-backed SessionStore. Sessions are stored as a pair
// of files under <root>/uploads/:
//
//   - <id>.info  — JSON-encoded tusd.FileInfo
//   - <id>       — staged binary bytes
//
// This is the same on-disk format used by OcisStore so existing sessions
// survive a rolling deploy that switches to FileStore.
type FileStore struct {
	root string
	opts TokenOptions
	log  *zerolog.Logger
}

// FileStoreFromDriverConf builds a FileStore from a reva driver config map.
// Returns nil if the config carries no root path (driver does not support
// coordinated uploads). Each service that mounts the same driver calls this
// independently.
func FileStoreFromDriverConf(driverConf map[string]interface{}, log *zerolog.Logger) *FileStore {
	if driverConf == nil {
		return nil
	}

	// storage_root is the ocm driver's spelling of root; it stages uploads under
	// the same <root>/uploads/ layout FileStore uses.
	type driverRootConf struct {
		Root            string `mapstructure:"root"`
		UploadDirectory string `mapstructure:"upload_directory"`
		StorageRoot     string `mapstructure:"storage_root"`
	}
	var rc driverRootConf
	_ = mapstructure.Decode(driverConf, &rc)

	root := rc.UploadDirectory
	if root == "" {
		root = rc.Root
	}
	if root == "" {
		root = rc.StorageRoot
	}
	if root == "" {
		return nil
	}

	return newFileStoreWithTokens(root, driverConf, log)
}

// NewFileStoreFromConfig builds a FileStore using uploadDir when set, falling
// back to root/upload_directory from the active driver config. This allows
// drivers that have no local root (e.g. KW) to still get a coordinator by
// setting upload_directory at the service level rather than inside the driver.
// Returns nil only when neither source resolves to a non-empty path.
func NewFileStoreFromConfig(uploadDir string, driverConf map[string]interface{}, log *zerolog.Logger) *FileStore {
	if uploadDir != "" {
		// Still take the tokens from the driver config: they sign the transfer URL
		// that postprocessing downloads the staged bytes from, and a service-level
		// upload directory says nothing about them.
		return newFileStoreWithTokens(uploadDir, driverConf, log)
	}
	return FileStoreFromDriverConf(driverConf, log)
}

// AsyncConf is how a service asks for async uploads: whether they are enabled,
// and the consumer subscription to use if they are.
type AsyncConf struct {
	Enabled       bool
	ConsumerGroup string
	NumConsumers  int
	// MountID is the storage id this provider answers for, used to drop
	// postprocessing events belonging to other storages.
	MountID string
}

// AsyncConfFromDriverConf reads the postprocessing settings off the driver config
// map the services already hand us.
//
// The keys are decomposedfs's (options.go: `asyncfileuploads`, `events`). Reading
// the driver's own keys rather than introducing service-level ones keeps a single
// source of truth: if the coordinator and the driver disagreed, uploads would
// either commit twice or never get scanned.
//
// The consumer group matters most. It is what makes retiring the driver's
// consumer a move rather than an addition: two consumers in one group take turns
// stealing each other's events, two in different groups both act and commit the
// same upload twice.
func AsyncConfFromDriverConf(driverConf map[string]interface{}) AsyncConf {
	if driverConf == nil {
		return AsyncConf{}
	}
	var ac struct {
		AsyncFileUploads bool   `mapstructure:"asyncfileuploads"`
		MountID          string `mapstructure:"mount_id"`
		Events           struct {
			NumConsumers  int    `mapstructure:"numconsumers"`
			ConsumerGroup string `mapstructure:"consumer_group"`
		} `mapstructure:"events"`
	}
	_ = mapstructure.Decode(driverConf, &ac)
	group := ac.Events.ConsumerGroup
	if group == "" {
		// decomposedfs's default (options.go:177). The coordinator takes over the
		// driver's subscription, so it must land in the same group.
		group = "dcfs"
	}
	return AsyncConf{
		Enabled:       ac.AsyncFileUploads,
		ConsumerGroup: group,
		NumConsumers:  ac.Events.NumConsumers,
		MountID:       ac.MountID,
	}
}

func newFileStoreWithTokens(root string, driverConf map[string]interface{}, log *zerolog.Logger) *FileStore {
	type tokenConf struct {
		DownloadEndpoint     string `mapstructure:"download_endpoint"`
		DataGatewayEndpoint  string `mapstructure:"datagateway_endpoint"`
		TransferSharedSecret string `mapstructure:"transfer_shared_secret"`
		TransferExpires      int64  `mapstructure:"transfer_expires"`
	}
	var tc tokenConf
	if tokens, ok := driverConf["tokens"]; ok {
		_ = mapstructure.Decode(tokens, &tc)
	}
	return NewFileStore(root, TokenOptions{
		DownloadEndpoint:     tc.DownloadEndpoint,
		DataGatewayEndpoint:  tc.DataGatewayEndpoint,
		TransferSharedSecret: tc.TransferSharedSecret,
		TransferExpires:      tc.TransferExpires,
	}, log)
}

// NewFileStore creates a FileStore rooted at root.
// root must be on a shared filesystem when multiple pods handle the same space.
func NewFileStore(root string, opts TokenOptions, log *zerolog.Logger) *FileStore {
	return &FileStore{root: root, opts: opts, log: log}
}

// Root returns the base directory of this FileStore.
func (fs *FileStore) Root() string {
	return fs.root
}

// Setup creates the uploads directory eagerly so permission problems are caught
// at startup rather than on the first upload.
func (fs *FileStore) Setup() error {
	return os.MkdirAll(filepath.Join(fs.root, "uploads"), 0700)
}

// New allocates a fresh session with a new UUID.
func (fs *FileStore) New(_ context.Context) Session {
	return &FileSession{
		store: fs,
		info: tusd.FileInfo{
			ID: uuid.New().String(),
			Storage: map[string]string{
				"Type": "OCISStore",
			},
			MetaData: tusd.MetaData{},
		},
	}
}

// Get loads the session with the given id from disk.
func (fs *FileStore) Get(ctx context.Context, id string) (Session, error) {
	infoPath := fileSessionPath(fs.root, id)

	data, err := os.ReadFile(infoPath)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ESTALE {
			return nil, tusd.ErrNotFound
		}
		if errors.Is(err, iofs.ErrNotExist) {
			return nil, tusd.ErrNotFound
		}
		return nil, err
	}

	var info tusd.FileInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	session := &FileSession{store: fs, info: info}

	stat, err := os.Stat(session.binPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, tusd.ErrNotFound
		}
		return nil, err
	}
	session.info.Offset = stat.Size()

	return session, nil
}

// List returns all sessions found under <root>/uploads/*.info.
func (fs *FileStore) List(ctx context.Context) ([]Session, error) {
	infoFiles, err := filepath.Glob(filepath.Join(fs.root, "uploads", "*.info"))
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(infoFiles))
	for _, path := range infoFiles {
		id := strings.TrimSuffix(filepath.Base(path), ".info")
		session, err := fs.Get(ctx, id)
		if err != nil {
			fs.log.Error().Str("path", path).Err(err).Msg("filestore: could not load session")
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
