package upload

import (
	"github.com/mitchellh/mapstructure"
)

// AsyncConf is how a service asks for async uploads: whether they are enabled,
// and the consumer subscription to use if they are.
type AsyncConf struct {
	Enabled       bool
	ConsumerGroup string
	NumConsumers  int
	// MountID is the storage id this provider answers for, used to drop
	// postprocessing events belonging to other storages.
	MountID string
}

// AsyncConfFromDriverConf reads the postprocessing settings off the driver's own
// config keys, so the coordinator and the driver cannot disagree about them.
func AsyncConfFromDriverConf(driverConf map[string]interface{}) AsyncConf {
	if driverConf == nil {
		return AsyncConf{}
	}
	var ac struct {
		AsyncFileUploads bool   `mapstructure:"asyncfileuploads"`
		MountID          string `mapstructure:"mount_id"`
		Events           struct {
			NumConsumers  int    `mapstructure:"numconsumers"`
			ConsumerGroup string `mapstructure:"consumer_group"`
		} `mapstructure:"events"`
	}
	_ = mapstructure.Decode(driverConf, &ac)
	group := ac.Events.ConsumerGroup
	if group == "" {
		// decomposedfs's default (options.go:177). The coordinator takes over the
		// driver's subscription, so it must land in the same group.
		group = "dcfs"
	}
	return AsyncConf{
		Enabled:       ac.AsyncFileUploads,
		ConsumerGroup: group,
		NumConsumers:  ac.Events.NumConsumers,
		MountID:       ac.MountID,
	}
}
