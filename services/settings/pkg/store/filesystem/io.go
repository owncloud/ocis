package store

import (
	"fmt"
	"io"
	"os"

	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Unmarshal file into record
func (s Store) parseRecordFromFile(record proto.Message, filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("%q: %w", filePath, settings.ErrNotFound)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(b) == 0 {
		return fmt.Errorf("%q: %w", filePath, settings.ErrNotFound)
	}

	if err := protojson.Unmarshal(b, record); err != nil {
		return err
	}
	return nil
}

// Marshal record into file
func (s Store) writeRecordToFile(record proto.Message, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	v, err := protojson.Marshal(record)
	if err != nil {
		return err
	}

	_, err = file.Write(v)
	if err != nil {
		return err
	}

	return nil
}
