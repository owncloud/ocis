package store

import (
	"testing"

	olog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/stretchr/testify/assert"
)

var bundleScenarios = []struct {
	name   string
	bundle *settingsmsg.Bundle
}{
	{
		name: "generic-test-file-resource",
		bundle: &settingsmsg.Bundle{
			Id:          bundle1,
			Type:        settingsmsg.Bundle_TYPE_DEFAULT,
			Extension:   extension1,
			DisplayName: "test1",
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_FILE,
				Id:   "beep",
			},
			Settings: []*settingsmsg.Setting{
				{
					Id:          setting1,
					Description: "test-desc-1",
					DisplayName: "test-displayname-1",
					Resource: &settingsmsg.Resource{
						Type: settingsmsg.Resource_TYPE_FILE,
						Id:   "bleep",
					},
					Value: &settingsmsg.Setting_IntValue{
						IntValue: &settingsmsg.Int{
							Min: 0,
							Max: 42,
						},
					},
				},
			},
		},
	},
	{
		name: "generic-test-system-resource",
		bundle: &settingsmsg.Bundle{
			Id:          bundle2,
			Type:        settingsmsg.Bundle_TYPE_DEFAULT,
			Extension:   extension2,
			DisplayName: "test1",
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_SYSTEM,
			},
			Settings: []*settingsmsg.Setting{
				{
					Id:          setting2,
					Description: "test-desc-2",
					DisplayName: "test-displayname-2",
					Resource: &settingsmsg.Resource{
						Type: settingsmsg.Resource_TYPE_SYSTEM,
					},
					Value: &settingsmsg.Setting_IntValue{
						IntValue: &settingsmsg.Int{
							Min: 0,
							Max: 42,
						},
					},
				},
			},
		},
	},
	{
		name: "generic-test-role-bundle",
		bundle: &settingsmsg.Bundle{
			Id:          bundle3,
			Type:        settingsmsg.Bundle_TYPE_ROLE,
			Extension:   extension1,
			DisplayName: "Role1",
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_SYSTEM,
			},
			Settings: []*settingsmsg.Setting{
				{
					Id:          setting3,
					Description: "test-desc-3",
					DisplayName: "test-displayname-3",
					Resource: &settingsmsg.Resource{
						Type: settingsmsg.Resource_TYPE_SETTING,
						Id:   setting1,
					},
					Value: &settingsmsg.Setting_PermissionValue{
						PermissionValue: &settingsmsg.Permission{
							Operation:  settingsmsg.Permission_OPERATION_READ,
							Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
						},
					},
				},
			},
		},
	},
}

func TestBundles(t *testing.T) {
	s := Store{
		dataPath: dataRoot,
		Logger: olog.NewLogger(
			olog.Color(true),
			olog.Pretty(true),
			olog.Level("info"),
		),
	}

	// write bundles
	for i := range bundleScenarios {
		index := i
		t.Run(bundleScenarios[index].name, func(t *testing.T) {
			filePath := s.buildFilePathForBundle(bundleScenarios[index].bundle.Id, true)
			if err := s.writeRecordToFile(bundleScenarios[index].bundle, filePath); err != nil {
				t.Error(err)
			}
			assert.FileExists(t, filePath)
		})
	}

	// check that ListBundles only returns bundles with type DEFAULT
	bundles, err := s.ListBundles(settingsmsg.Bundle_TYPE_DEFAULT, []string{})
	if err != nil {
		t.Error(err)
	}
	for i := range bundles {
		assert.Equal(t, settingsmsg.Bundle_TYPE_DEFAULT, bundles[i].Type)
	}

	// check that ListBundles filtered by an id only returns that bundle
	filteredBundles, err := s.ListBundles(settingsmsg.Bundle_TYPE_DEFAULT, []string{bundle2})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(filteredBundles))
	if len(filteredBundles) == 1 {
		assert.Equal(t, bundle2, filteredBundles[0].Id)
	}

	// check that ListRoles only returns bundles with type ROLE
	roles, err := s.ListBundles(settingsmsg.Bundle_TYPE_ROLE, []string{})
	if err != nil {
		t.Error(err)
	}
	for i := range roles {
		assert.Equal(t, settingsmsg.Bundle_TYPE_ROLE, roles[i].Type)
	}

	burnRoot()
}
