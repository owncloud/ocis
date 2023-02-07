package relations

import (
	"context"
	"net/url"
	"strings"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

const (
	LibreGraphIDProp             = "http://libregraph.org/prop/user/id"
	LibreGraphSamAccountNameProp = "http://libregraph.org/prop/user/onPremisesSamAccountName"
	LibreGraphMailProp           = "http://libregraph.org/prop/user/mail"
	LibreGraphDisplayNameProp    = "http://libregraph.org/prop/user/displayName"
)

type ownCloudAccount struct {
	subject url.URL
	next    service.Webfinger
}

func OwnCloudAccount(url url.URL, next service.Webfinger) service.Webfinger {
	if next == nil {
		next = Noop{}
	}

	return &ownCloudAccount{
		subject: url,
		next:    next,
	}
}

func (l *ownCloudAccount) Lookup(ctx context.Context, jrd *webfinger.JSONResourceDescriptor) {
	if jrd == nil {
		jrd = &webfinger.JSONResourceDescriptor{}
	}
	if strings.HasPrefix("acct:me", jrd.Subject) {
		// TODO check if this relation was requested
		if u, ok := revactx.ContextGetUser(ctx); ok {
			// return correct account based on id
			jrd.Subject = "acct:" + u.GetId().GetOpaqueId() + "@" + l.subject.Host + l.subject.Path
			jrd.Properties[LibreGraphIDProp] = u.GetId().GetOpaqueId()
			jrd.Properties[LibreGraphSamAccountNameProp] = u.GetUsername()
			jrd.Properties[LibreGraphMailProp] = u.GetMail()
			jrd.Properties[LibreGraphDisplayNameProp] = u.GetDisplayName()
		} else {
			// todo if we don't know the user return a 404, well, in this case a 401
		}
	}
	l.next.Lookup(ctx, jrd)
}
func (l *ownCloudAccount) Next(next service.Webfinger) {
	l.next = next
}
