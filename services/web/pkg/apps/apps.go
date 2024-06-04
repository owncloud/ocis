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
	// ErrInvalidApp is the error when an app is invalid
	ErrInvalidApp = errors.New("invalid app")

	// ErrMissingManifest is the error when the manifest is missing
	ErrMissingManifest = errors.New("missing manifest")

	// ErrInvalidManifest is the error when the manifest is invalid
	ErrInvalidManifest = errors.New("invalid manifest")

	// ErrEntrypointDoesNotExist is the error when the entrypoint does not exist or is not a file
	ErrEntrypointDoesNotExist = errors.New("entrypoint does not exist")

	validate = validator.New(validator.WithRequiredStructEnabled())
)

const (
	// _manifest is the name of the manifest file for an application
	_manifest = "manifest.json"
)

// Application contains the metadata of an application
type Application struct {
	// ID is the unique identifier of the application
	ID string `json:"-"`

	// Entrypoint is the entrypoint of the application within the bundle
	Entrypoint string `json:"entrypoint" validate:"required"`

	// Config contains the application-specific configuration
	Config map[string]interface{} `json:"config,omitempty"`
}

// ToExternal converts an Application to an ExternalApp configuration
func (a Application) ToExternal(entrypoint string) config.ExternalApp {
	return config.ExternalApp{
		ID:     a.ID,
		Path:   filepathx.JailJoin(entrypoint, a.Entrypoint),
		Config: a.Config,
	}
}

// List returns a list of applications from the given filesystems,
// individual filesystems are searched for applications, and the list is merged.
// Last finding gets priority in case of conflicts, so the order of the filesystems is important.
func List(logger log.Logger, data map[string]config.App, fSystems ...fs.FS) []Application {
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
			if data, ok := data[name]; ok {
				appData = data
			}

			if appData.Disabled {
				// if the app is disabled, skip it
				continue
			}

			application, err := build(fSystem, name, appData.Config)
			if err != nil {
				// if app creation fails, log the error and continue with the next app
				logger.Debug().Err(err).Str("path", entry.Name()).Msg("failed to load application")
				continue
			}

			// everything is fine, add the application to the list of applications
			registry[name] = application
		}
	}

	return maps.Values(registry)
}

func build(fSystem fs.FS, id string, conf map[string]any) (Application, error) {
	// skip non-directory listings, every app needs to be contained inside a directory
	entry, err := fs.Stat(fSystem, id)
	if err != nil || !entry.IsDir() {
		return Application{}, ErrInvalidApp
	}

	// read the manifest.json from the app directory.
	manifest := path.Join(id, _manifest)
	reader, err := fSystem.Open(manifest)
	if err != nil {
		// manifest.json is required
		return Application{}, errors.Join(err, ErrMissingManifest)
	}
	defer reader.Close()

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

	// the entrypoint is jailed to the app directory
	application.Entrypoint = filepathx.JailJoin(id, application.Entrypoint)
	info, err := fs.Stat(fSystem, application.Entrypoint)
	switch {
	case err != nil:
		return Application{}, errors.Join(err, ErrEntrypointDoesNotExist)
	case info.IsDir():
		return Application{}, ErrEntrypointDoesNotExist
	}

	application.ID = id

	return application, nil
}
