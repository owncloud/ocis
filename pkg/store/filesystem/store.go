// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/owncloud/ocis-settings/pkg/settings"
)

var (
	// StoreName is the default name for the settings store
	StoreName     = "ocis-settings-store"
	managerName   = "filesystem"
	emptyKeyError = "key cannot be empty"
)

// Store interacts with the filesystem to manage settings information
type Store struct {
	mountPath string
	Logger    olog.Logger
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	s := Store{}

	dest := path.Join(cfg.Storage.RootMountPath, StoreName)
	if _, err := os.Stat(dest); err != nil {
		s.Logger.Info().Msgf("creating container on %v", dest)
		err := os.MkdirAll(dest, 0700)
		if err != nil {
			s.Logger.Err(err).Msgf("providing container on %v", dest)
		}
	}

	s.mountPath = dest
	return &s
}

// List returns all the bundles in the mountPath folder
func (s Store) List() ([]*proto.SettingsBundle, error) {
	records := []*proto.SettingsBundle{}
	bundles, err := ioutil.ReadDir(s.mountPath)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", s.mountPath)
		return nil, err
	}

	s.Logger.Info().Msg("listing bundles")
	for _, v := range bundles {
		records = append(records, parseFileName(v.Name()))
	}

	return records, nil
}

// ListByExtension returns all bundles in the mountPath folder belonging to the given extension
func (s Store) ListByExtension(extension string) ([]*proto.SettingsBundle, error) {
	records := []*proto.SettingsBundle{}
	bundles, err := ioutil.ReadDir(s.mountPath)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", s.mountPath)
		return nil, err
	}

	s.Logger.Info().Msgf("listing bundles by extension %v", extension)
	for _, v := range bundles {
		record := parseFileName(v.Name())
		if record.Extension == extension {
			records = append(records, record)
		}
	}

	return records, nil
}

// Read tries to find a bundle by the given extension and key within the mountPath
func (s Store) Read(extension string, key string) (*proto.SettingsBundle, error) {
	fileName := buildFileNameFromData(extension, key)
	contents, err := os.Open(path.Join(s.mountPath, fileName))
	if err != nil {
		s.Logger.Err(err).Msgf("error reading contents for extension %v and key %v: file not found", extension, key)
		return nil, err
	}

	record := proto.SettingsBundle{}
	if err = jsonpb.Unmarshal(contents, &record); err != nil {
		s.Logger.Err(err).Msg("error unmarshalling record")
		return nil, err
	}

	return &record, nil
}

// Write writes the given record into a file within the mountPath
func (s Store) Write(rec *proto.SettingsBundle) (*proto.SettingsBundle, error) {
	if len(rec.Key) < 1 {
		s.Logger.Error().Msg("key cannot be empty")
		return nil, fmt.Errorf(emptyKeyError)
	}

	marshaler := jsonpb.Marshaler{}
	recordPath := path.Join(s.mountPath, buildFileNameFromBundle(rec))
	if err := ioutil.WriteFile(recordPath, []byte{}, 0644); err != nil {
		return nil, err
	}

	fd, err := os.OpenFile(recordPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		s.Logger.Err(err).
			Str(
				"finding file",
				fmt.Sprintf("file `%v` not found on store: `%v`", recordPath, s.mountPath),
			)
	}

	if err = marshaler.Marshal(fd, rec); err != nil {
		s.Logger.Err(err).
			Str(
				"marshaling record",
				fmt.Sprintf("error marshaling record: %+v", rec),
			)
	}

	s.Logger.Info().Msgf("request contents written to file: %v", recordPath)
	return rec, nil
}

// Builds a unique file name from the given bundle
func buildFileNameFromBundle(bundle *proto.SettingsBundle) string {
	return buildFileNameFromData(bundle.Extension, bundle.Key)
}

// Builds a unique file name from the given params
func buildFileNameFromData(extension string, key string) string {
	return extension + "__" + key + ".json"
}

// Extracts extension and key from the given fileName and builds a (minimalistic) bundle from it
func parseFileName(fileName string) *proto.SettingsBundle {
	parts := strings.Split(strings.Replace(fileName, ".json", "", 1), "__")
	return &proto.SettingsBundle{
		Key:       parts[1],
		Extension: parts[0],
	}
}

func init() {
	settings.Registry[managerName] = New
}
