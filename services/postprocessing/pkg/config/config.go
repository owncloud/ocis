package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Log *Log `yaml:"log"`

	Postprocessing Postprocessing `yaml:"postprocessing"`

	Context context.Context `yaml:"-"`
}

// Postprocessing definces the config options for the postprocessing service.
type Postprocessing struct {
	Events          Events        `yaml:"events"`
	Virusscan       bool          `yaml:"virusscan" env:"POSTPROCESSING_VIRUSSCAN" desc:"should the system do a virusscan? Needs antivirus service"`
	FTSIndex        bool          `yaml:"fulltextsearch" env:"POSTPROCESSING_FTS" desc:"should the system index files for fts? Needs search service"`
	Delayprocessing time.Duration `yaml:"delayprocessing" env:"POSTPROCESSING_DELAY" desc:"the sytem sleeps for this time while postprocessing"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"POSTPROCESSING_EVENTS_ENDPOINT" desc:"Endpoint of the event system."`
	Cluster  string `yaml:"cluster" env:"POSTPROCESSING_EVENTS_CLUSTER" desc:"Cluster ID of the event system."`
}
