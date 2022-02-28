// Package store implements the go-micro store interface
package store

import (
	"encoding/json"
	"errors"
	"fmt"

	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

// ListBundles returns all bundles in the dataPath folder that match the given type.
func (s Store) ListBundles(bundleType settingsmsg.Bundle_Type, bundleIDs []string) ([]*settingsmsg.Bundle, error) {
	// TODO: this is needed for initialization - we need to find a better way to fix this
	if s.mdc == nil && len(bundleIDs) == 1 && bundleIDs[0] == "71881883-1768-46bd-a24d-a356a2afdf7f" {
		return []*settingsmsg.Bundle{{
			Id: "71881883-1768-46bd-a24d-a356a2afdf7f",
			Settings: []*settingsmsg.Setting{
				{
					Id:          "8e587774-d929-4215-910b-a317b1e80f73",
					Name:        "account-management",
					DisplayName: "Account Management",
					Description: "This permission gives full access to everything that is related to account management.",
					Resource: &settingsmsg.Resource{
						Type: settingsmsg.Resource_TYPE_USER,
						Id:   "all",
					},
					Value: &settingsmsg.Setting_PermissionValue{
						PermissionValue: &settingsmsg.Permission{
							Operation:  settingsmsg.Permission_OPERATION_READWRITE,
							Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
						},
					},
				},
			},
		}}, nil
	}
	s.Init()
	var bundles []*settingsmsg.Bundle
	for _, id := range bundleIDs {
		b, err := s.mdc.SimpleDownload(nil, bundlePath(id))
		if err != nil {
			return nil, err
		}

		bundle := &settingsmsg.Bundle{}
		err = json.Unmarshal(b, bundle)
		if err != nil {
			return nil, err
		}

		if bundle.Type == bundleType {
			bundles = append(bundles, bundle)
		}

	}
	return bundles, nil
}

// ReadBundle tries to find a bundle by the given id within the dataPath.
func (s Store) ReadBundle(bundleID string) (*settingsmsg.Bundle, error) {
	s.Init()
	b, err := s.mdc.SimpleDownload(nil, bundlePath(bundleID))
	if err != nil {
		return nil, err
	}

	bundle := &settingsmsg.Bundle{}
	return bundle, json.Unmarshal(b, bundle)
}

// ReadSetting tries to find a setting by the given id within the dataPath.
func (s Store) ReadSetting(settingID string) (*settingsmsg.Setting, error) {
	return nil, errors.New("not implemented")
}

// WriteBundle sends the givens record to the metadataclient. returns `record` for legacy reasons
func (s Store) WriteBundle(record *settingsmsg.Bundle) (*settingsmsg.Bundle, error) {
	s.Init()
	b, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	return record, s.mdc.SimpleUpload(nil, bundlePath(record.Id), b)
}

// AddSettingToBundle adds the given setting to the bundle with the given bundleID.
func (s Store) AddSettingToBundle(bundleID string, setting *settingsmsg.Setting) (*settingsmsg.Setting, error) {
	return nil, errors.New("not implemented")
}

// RemoveSettingFromBundle removes the setting from the bundle with the given ids.
func (s Store) RemoveSettingFromBundle(bundleID string, settingID string) error {
	return errors.New("not implemented")
}

func bundlePath(id string) string {
	return fmt.Sprintf("%s/%s", bundleFolderLocation, id)
}
