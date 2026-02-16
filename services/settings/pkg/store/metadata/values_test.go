package store

import (
	"testing"

	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/stretchr/testify/require"
)

var valueScenarios = []struct {
	name  string
	value *settingsmsg.Value
}{
	{
		name: "generic-test-with-system-resource",
		value: &settingsmsg.Value{
			Id:          value1,
			BundleId:    bundle1,
			SettingId:   setting1,
			AccountUuid: accountUUID1,
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_SYSTEM,
			},
			Value: &settingsmsg.Value_StringValue{
				StringValue: "lalala",
			},
		},
	},
	{
		name: "generic-test-with-file-resource",
		value: &settingsmsg.Value{
			Id:          value2,
			BundleId:    bundle2,
			SettingId:   setting2,
			AccountUuid: accountUUID1,
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_FILE,
				Id:   "adfba82d-919a-41c3-9cd1-5a3f83b2bf76",
			},
			Value: &settingsmsg.Value_StringValue{
				StringValue: "tralala",
			},
		},
	},
	{
		name: "value without accountUUID",
		value: &settingsmsg.Value{
			Id:          value3,
			BundleId:    bundle3,
			SettingId:   setting2,
			AccountUuid: "",
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_FILE,
				Id:   "adfba82d-919a-41c3-9cd1-5a3f83b2bf76",
			},
			Value: &settingsmsg.Value_StringValue{
				StringValue: "tralala",
			},
		},
	},
}

func TestValues(t *testing.T) {
	for i := range valueScenarios {
		index := i
		t.Run(valueScenarios[index].name, func(t *testing.T) {
			s := initStore()
			setupRoles(s)
			value := valueScenarios[index].value
			v, err := s.WriteValue(value)
			require.NoError(t, err)
			require.Equal(t, value, v)

			v, err = s.ReadValue(value.Id)
			require.NoError(t, err)
			require.Equal(t, value, v)
		})
	}
}

func TestListValues(t *testing.T) {
	s := initStore()
	setupRoles(s)
	for _, v := range valueScenarios {
		_, err := s.WriteValue(v.value)
		require.NoError(t, err)
	}

	// empty accountid returns only values with empty accountud
	vs, err := s.ListValues("", "")
	require.NoError(t, err)
	require.Len(t, vs, 1)

	// filled accountid returns matching and empty accountUUID values
	vs, err = s.ListValues("", accountUUID1)
	require.NoError(t, err)
	require.Len(t, vs, 3)

	// filled bundleid only returns matching values
	vs, err = s.ListValues(bundle3, accountUUID1)
	require.NoError(t, err)
	require.Len(t, vs, 1)

}
