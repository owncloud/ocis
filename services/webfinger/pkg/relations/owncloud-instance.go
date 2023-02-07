package relations

import (
	"context"
	"regexp"

	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/config"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

const (
	OwnCloudInstanceRel = "http://webfinger.owncloud/rel/server-instance"
)

type compiledInstance struct {
	config.Instance
	compiledRegex *regexp.Regexp
}

type ownCloudInstance struct {
	next      service.Webfinger
	instances []compiledInstance
}

func OwnCloudInstance(instances []config.Instance, next service.Webfinger) service.Webfinger {
	if next == nil {
		next = Noop{}
	}
	compiledInstances := make([]compiledInstance, 0, len(instances))
	var err error
	for _, instance := range instances {
		compiled := compiledInstance{Instance: instance}
		compiled.compiledRegex, err = regexp.Compile(instance.Regex)
		if err != nil {
			// TODO return error
		}
	}

	return &ownCloudInstance{
		instances: compiledInstances,
		next:      next,
	}
}

func (l *ownCloudInstance) Lookup(ctx context.Context, jrd *webfinger.JSONResourceDescriptor) {
	if jrd == nil {
		jrd = &webfinger.JSONResourceDescriptor{}
	}
	if claims := oidc.FromContext(ctx); claims != nil {
		for _, instance := range l.instances {
			if value, ok := claims[instance.Claim].(string); ok && instance.compiledRegex.MatchString(value) {
				jrd.Links = append(jrd.Links, webfinger.Link{
					Rel:    OpenIDConnectRel,
					Href:   instance.Href, // allow a template?
					Titles: instance.Titles,
				})
			}
		}
	}
	l.next.Lookup(ctx, jrd)
}

func (l *ownCloudInstance) Next(next service.Webfinger) {
	l.next = next
}
