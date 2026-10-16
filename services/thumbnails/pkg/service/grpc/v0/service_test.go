package svc

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config/defaults"
)

func TestGifFirstFrameOnlyServiceConfig(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		cfg := defaults.FullDefaultConfig()
		cfg.Thumbnail.GifFirstFrameOnly = enabled
		svc := NewService(Config(cfg), Logger(log.NewLogger())).(Thumbnail)
		if svc.preprocessorOpts.GifFirstFrameOnly != enabled {
			t.Fatalf("first-frame setting %v was not passed to the service", enabled)
		}
	}
}
