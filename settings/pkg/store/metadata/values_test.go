package store

import (
	"testing"

	olog "github.com/owncloud/ocis/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
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
			BundleId:    bundle1,
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
}

func TestValues(t *testing.T) {
	mdc := NewMDC()
	s := Store{
		Logger: olog.NewLogger(
			olog.Color(true),
			olog.Pretty(true),
			olog.Level("info"),
		),
		mdc: mdc,
	}
	for i := range valueScenarios {
		index := i
		t.Run(valueScenarios[index].name, func(t *testing.T) {
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
