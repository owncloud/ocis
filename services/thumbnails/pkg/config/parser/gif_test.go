package parser

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config/defaults"
)

func TestGifFirstFrameOnlyEnvironment(t *testing.T) {
	if defaults.DefaultConfig().Thumbnail.GifFirstFrameOnly {
		t.Fatal("GIF thumbnails must preserve animation by default")
	}
	for _, value := range []string{"true", "false"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("THUMBNAILS_GIF_FIRST_FRAME_ONLY", value)
			cfg := defaults.FullDefaultConfig()
			if err := envdecode.Decode(cfg); err != nil {
				t.Fatal(err)
			}
			if cfg.Thumbnail.GifFirstFrameOnly != (value == "true") {
				t.Fatalf("environment value %q was not applied", value)
			}
		})
	}
}
