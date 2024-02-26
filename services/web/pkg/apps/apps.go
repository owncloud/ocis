package apps

import (
	"encoding/json"
	"errors"
	"io/fs"
	"path"

	"dario.cat/mergo"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/maps"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/x/path/filepathx"
	"github.com/owncloud/ocis/v2/services/web/pkg/config"
)

var (
	validate           *validator.Validate
	logger             = log.Default()
	ErrInvalidApp      = errors.New("invalid app")
	ErrMissingManifest = errors.New("missing manifest")
	ErrInvalidManifest = errors.New("invalid manifest")
)

const (
	// _manifest is the name of the manifest file for an application
	_manifest = "manifest.json"
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// Application contains the metadata of an application
type Application struct {
	// ID is the unique identifier of the application
	ID string `json:"id" validate:"required"`

	// Entrypoint is the entrypoint of the application within the bundle
	Entrypoint string `json:"entrypoint" validate:"required"`

	// Config contains the application specific configuration
	Config map[string]interface{} `json:"config,omitempty"`
}

// ToExternal converts an Application to an ExternalApp configuration
func (a Application) ToExternal(entrypoint string) config.ExternalApp {
	return config.ExternalApp{
		ID:     a.ID,
		Path:   filepathx.JailJoin(entrypoint, a.ID, a.Entrypoint),
		Config: a.Config,
	}
}

// List returns a list of applications from the given filesystems,
// individual filesystems are searched for applications, and the list is merged.
// Last finding gets priority in case of conflicts, so the order of the filesystems is important.
func List(appsData map[string]config.App, fSystems ...fs.FS) []Application {
	registry := map[string]Application{}

	for _, fSystem := range fSystems {
		if fSystem == nil {
			continue
		}

		entries, err := fs.ReadDir(fSystem, ".")
		if err != nil {
			// skip non-directory listings, every app needs to be contained inside a directory
			continue
		}

		for _, entry := range entries {
			var appData config.App
			name := entry.Name()

			// configuration for the application is optional, if it is not present, the default configuration is used
			if data, ok := appsData[name]; ok {
				appData = data
			}

			if appData.Disabled {
				// if the app is disabled, skip it
				continue
			}

			application, err := Build(fSystem, name, appData.Config)
			if err != nil {
				// if app creation fails, log the error and continue with the next app
				logger.Debug().Err(err).Str("path", entry.Name()).Msg("failed to load application")
				continue
			}

			// everything is fine, add the application to the list of applications
			registry[application.ID] = application
		}
	}

	return maps.Values(registry)
}

func Build(fSystem fs.FS, name string, conf map[string]any) (Application, error) {
	// skip non-directory listings, every app needs to be contained inside a directory
	entry, err := fs.Stat(fSystem, name)
	if err != nil || !entry.IsDir() {
		return Application{}, ErrInvalidApp
	}

	// read the manifest.json from the app directory.
	manifest := path.Join(path.Base(name), _manifest)
	reader, err := fSystem.Open(manifest)
	if err != nil {
		// manifest.json is required
		return Application{}, errors.Join(err, ErrMissingManifest)
	}

	var application Application
	if json.NewDecoder(reader).Decode(&application) != nil {
		// a valid manifest.json is required
		return Application{}, errors.Join(err, ErrInvalidManifest)
	}

	if err := validate.Struct(application); err != nil {
		// the application is required to be valid
		return Application{}, errors.Join(err, ErrInvalidManifest)
	}

	// overload the default configuration with the application-specific configuration,
	// the application-specific configuration has priority, and failing is fine here
	_ = mergo.Merge(&application.Config, conf, mergo.WithOverride)

	return application, nil
}
