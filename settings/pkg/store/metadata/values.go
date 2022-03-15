// Package store implements the go-micro store interface
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

// ListValues reads all values that match the given bundleId and accountUUID.
// If the bundleId is empty, it's ignored for filtering.
// If the accountUUID is empty, only values with empty accountUUID are returned.
// If the accountUUID is not empty, values with an empty or with a matching accountUUID are returned.
func (s *Store) ListValues(bundleID, accountUUID string) ([]*settingsmsg.Value, error) {
	s.Init()
	ctx := context.TODO()

	vIDs, err := s.mdc.ReadDir(ctx, valuesFolderLocation)
	if err != nil {
		return nil, err
	}

	// TODO: refine logic not to spam metadata service
	var values []*settingsmsg.Value
	for _, vid := range vIDs {
		b, err := s.mdc.SimpleDownload(ctx, valuePath(vid))
		if err != nil {
			return nil, err
		}

		v := &settingsmsg.Value{}
		err = json.Unmarshal(b, v)
		if err != nil {
			return nil, err
		}

		if bundleID != "" && v.BundleId != bundleID {
			continue
		}

		if v.AccountUuid == "" {
			values = append(values, v)
			continue
		}

		if v.AccountUuid == accountUUID {
			values = append(values, v)
			continue
		}
	}
	return values, nil
}

// ReadValue tries to find a value by the given valueId within the dataPath
func (s *Store) ReadValue(valueID string) (*settingsmsg.Value, error) {
	s.Init()
	ctx := context.TODO()

	b, err := s.mdc.SimpleDownload(ctx, valuePath(valueID))
	if err != nil {
		return nil, err
	}
	val := &settingsmsg.Value{}
	return val, json.Unmarshal(b, val)
}

// ReadValueByUniqueIdentifiers tries to find a value given a set of unique identifiers
func (s *Store) ReadValueByUniqueIdentifiers(accountUUID, settingID string) (*settingsmsg.Value, error) {
	fmt.Println("ReadValueByUniqueIdentifiers not implemented")
	return nil, errors.New("not implemented")
}

// WriteValue writes the given value into a file within the dataPath
func (s *Store) WriteValue(value *settingsmsg.Value) (*settingsmsg.Value, error) {
	s.Init()
	ctx := context.TODO()

	if value.Id == "" {
		value.Id = uuid.Must(uuid.NewV4()).String()
	}
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return value, s.mdc.SimpleUpload(ctx, valuePath(value.Id), b)
}

func valuePath(id string) string {
	return fmt.Sprintf("%s/%s", valuesFolderLocation, id)
}
