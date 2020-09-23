package store

import (
	"testing"

	olog "github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
)

var valueScenarios = []struct {
	name  string
	value *proto.Value
}{
	{
		name: "generic-test-with-system-resource",
		value: &proto.Value{
			Id:          value1,
			BundleId:    bundle1,
			SettingId:   setting1,
			AccountUuid: accountUUID1,
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_SYSTEM,
			},
			Value: &proto.Value_StringValue{
				StringValue: "lalala",
			},
		},
	},
	{
		name: "generic-test-with-file-resource",
		value: &proto.Value{
			Id:          value2,
			BundleId:    bundle1,
			SettingId:   setting2,
			AccountUuid: accountUUID1,
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_FILE,
				Id:   "adfba82d-919a-41c3-9cd1-5a3f83b2bf76",
			},
			Value: &proto.Value_StringValue{
				StringValue: "tralala",
			},
		},
	},
}

func TestValues(t *testing.T) {
	s := Store{
		dataPath: dataRoot,
		Logger: olog.NewLogger(
			olog.Color(true),
			olog.Pretty(true),
			olog.Level("info"),
		),
	}
	for i := range valueScenarios {
		index := i
		t.Run(valueScenarios[index].name, func(t *testing.T) {

			filePath := s.buildFilePathForValue(valueScenarios[index].value.Id, true)
			if err := s.writeRecordToFile(valueScenarios[index].value, filePath); err != nil {
				t.Error(err)
			}
			assert.FileExists(t, filePath)
		})
	}

	burnRoot()
}
