// Package store implements the go-micro store interface
package store

import (
	"encoding/json"
	"errors"
	"fmt"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

// ListValues reads all values that match the given bundleId and accountUUID.
// If the bundleId is empty, it's ignored for filtering.
// If the accountUUID is empty, only values with empty accountUUID are returned.
// If the accountUUID is not empty, values with an empty or with a matching accountUUID are returned.
func (s Store) ListValues(bundleID, accountUUID string) ([]*settingsmsg.Value, error) {
	return nil, errors.New("not implemented")
}

// ReadValue tries to find a value by the given valueId within the dataPath
func (s Store) ReadValue(valueID string) (*settingsmsg.Value, error) {
	b, err := s.mdc.SimpleDownload(nil, valuePath(valueID))
	if err != nil {
		return nil, err
	}
	val := &settingsmsg.Value{}
	return val, json.Unmarshal(b, val)
}

// ReadValueByUniqueIdentifiers tries to find a value given a set of unique identifiers
func (s Store) ReadValueByUniqueIdentifiers(accountUUID, settingID string) (*settingsmsg.Value, error) {
	return nil, errors.New("not implemented")
}

// WriteValue writes the given value into a file within the dataPath
func (s Store) WriteValue(value *settingsmsg.Value) (*settingsmsg.Value, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return value, s.mdc.SimpleUpload(nil, valuePath(value.Id), b)
}

func valuePath(id string) string {
	return fmt.Sprintf("%s/%s", valuesFolderLocation, id)
}
