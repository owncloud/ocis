package relations

import (
	"context"

	"github.com/owncloud/ocis/v2/services/webfinger/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

type Noop struct{}

func (l Noop) Lookup(ctx context.Context, jrd *webfinger.JSONResourceDescriptor) {
}
func (l Noop) Next(next service.Webfinger) {
}
