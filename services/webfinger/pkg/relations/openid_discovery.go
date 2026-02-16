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
}

// OpenIDDiscovery adds the Openid Connect issuer relation
func OpenIDDiscovery(href string) service.RelationProvider {
	return &openIDDiscovery{
		Href: href,
	}
}

func (l *openIDDiscovery) Add(_ context.Context, jrd *webfinger.JSONResourceDescriptor) {
	if jrd == nil {
		jrd = &webfinger.JSONResourceDescriptor{}
	}
	jrd.Links = append(jrd.Links, webfinger.Link{
		Rel:  OpenIDConnectRel,
		Href: l.Href,
	})
}
