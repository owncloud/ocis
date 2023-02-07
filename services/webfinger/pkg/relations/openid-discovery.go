package relations

import (
	"context"

	"github.com/owncloud/ocis/v2/services/webfinger/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

const (
	OpenIDConnectRel = "http://openid.net/specs/connect/1.0/issuer"
)

type openIDDiscovery struct {
	Href string
	next service.Webfinger
}

func OpenIDDiscovery(href string, next service.Webfinger) service.Webfinger {
	if next == nil {
		next = Noop{}
	}
	return &openIDDiscovery{
		Href: href,
		next: next,
	}
}

func (l *openIDDiscovery) Lookup(ctx context.Context, jrd *webfinger.JSONResourceDescriptor) {
	if jrd == nil {
		jrd = &webfinger.JSONResourceDescriptor{}
	}
	// TODO check if this relation was requested
	jrd.Links = append(jrd.Links, webfinger.Link{
		Rel:  OpenIDConnectRel,
		Href: l.Href,
		// Titles: , // TODO use , separated env var with : separated language -> title pairs
	})
	l.next.Lookup(ctx, jrd)
}
func (l *openIDDiscovery) Next(next service.Webfinger) {
	l.next = next
}
