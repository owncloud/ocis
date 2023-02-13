package relations

import (
	"context"
	"testing"

	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

func TestOpenidDiscovery(t *testing.T) {
	provider := OpenIDDiscovery("http://issuer.url")

	jrd := webfinger.JSONResourceDescriptor{}

	provider.Add(context.Background(), &jrd)

	if len(jrd.Links) != 1 {
		t.Errorf("provider returned wrong number of links: %v, expected 1", len(jrd.Links))
	}
	if jrd.Links[0].Href != "http://issuer.url" {
		t.Errorf("provider returned wrong issuer link href: %v, expected %v", jrd.Links[0].Href, "http://issuer.url")
	}
	if jrd.Links[0].Rel != "http://openid.net/specs/connect/1.0/issuer" {
		t.Errorf("provider returned wrong openid connect rel: %v, expected %v", jrd.Links[0].Href, OpenIDConnectRel)
	}
}
