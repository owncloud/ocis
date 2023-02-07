package service

import (
	"context"
	"net/url"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/config"
	"github.com/owncloud/ocis/v2/services/webfinger/pkg/webfinger"
)

const (
	OwnCloudInstanceRel = "http://webfinger.owncloud/rel/server-instance"
	OpenIDConnectRel    = "http://openid.net/specs/connect/1.0/issuer"
)

// Service defines the extension handlers.
type Service interface {
	// Webfinger is the endpoint for retrieving various href relations.
	//
	//	GET /.well-known/webfinger?
	//	     resource=acct%3Acarol%40example.com&
	//	     rel=http%3A%2F%2Fwebfinger.owncloud%rel%2Fserver-instance
	//	     HTTP/1.1
	//	Host: example.com
	//
	// The server might respond like this:
	//
	//	HTTP/1.1 200 OK
	//	Access-Control-Allow-Origin: *
	//	Content-Type: application/jrd+json
	//
	//	{
	//	  "subject" : "acct:carol@example.com",
	//	  "links" :
	//	  [
	//	    {
	//	      "rel" : "http://webfinger.owncloud/rel/server-instance",
	//	      "href" : "https://instance.example.com",
	//	      "titles": {
	//	        "en": "Readable Instance Name"
	//	      }
	//	    },
	//	    {
	//	      "rel" : "http://webfinger.owncloud/rel/server-instance",
	//	      "href" : "https://otherinstance.example.com",
	//	      "titles": {
	//	        "en": "Other Readable Instance Name"
	//	      }
	//	    }
	//	  ]
	//	}
	Webfinger(ctx context.Context, queryTarget *url.URL, rels []string) (webfinger.JSONResourceDescriptor, error)
}

type Webfinger interface {
	Lookup(ctx context.Context, jrd *webfinger.JSONResourceDescriptor)
	Next(next Webfinger)
}

type InstanceSelector interface {
	GetInstanceIds(ctx context.Context, account string) []string
}

type InstanceLookup interface {
	GetInstance(ctx context.Context, id string) Instance
	// get multiple instances at once?
}

type Instance struct {
	Href   string
	Titles map[string]string
}

type DefaultInstanceSelector struct{}

func (s DefaultInstanceSelector) GetInstanceIds(ctx context.Context, account string) []string {
	return []string{"default"}
}

type DefaultInstanceLookup struct{}

func (l DefaultInstanceLookup) GetInstance(ctx context.Context, id string) Instance {
	if id == "default" {
		return Instance{
			Href: ctx.Value("href").(string),
			Titles: map[string]string{
				"en": "ownCloud Infinite Scale",
			},
		}
	}
	return Instance{}
}

// New returns a new instance of Service
func New(opts ...Option) (Service, error) {
	options := newOptions(opts...)

	// TODO use fallback implementations of InstanceIdLookup and InstanceLookup?
	// The InstanceIdLookup may have to happen earlier?

	return svc{
		log:         options.Logger,
		config:      options.Config,
		lookupChain: options.LookupChain,
	}, nil
}

type svc struct {
	config      *config.Config
	log         log.Logger
	lookupChain Webfinger
}

// TODO implement different implementations:
// static one returning the href or a configureable domain
// regex one returning different instances based on the regex that matches
// claim one that reads a claim and then fetches the instance?
// that is actually two interfaces / steps:
// - one that determines the instances/schools id (read from claim, regex match)
// - one that looks up in instance by id (use template, read from json, read from ldap, read from graph)

// Webfinger implements the service interface
func (s svc) Webfinger(ctx context.Context, queryTarget *url.URL, rel []string) (webfinger.JSONResourceDescriptor, error) {

	jrd := webfinger.JSONResourceDescriptor{
		Subject: queryTarget.String(),
	}

	// TODO acct chain vs https: chain?
	switch queryTarget.Scheme {
	case "acct":
		s.lookupChain.Lookup(ctx, &jrd)
	case "http", "https":
		s.lookupChain.Lookup(ctx, &jrd)
	default:
		return jrd, ErrNotFound
	}

	return jrd, nil
	/*
	   instanceIds := s.instanceIdSelector.GetInstanceIds(ctx, strings.TrimPrefix(queryTarget.String(), "acct:"))

	   href := ctx.Value("href").(string)

	   links := make([]webfinger.Link, 0, len(instanceIds))
	   // TODO, make listing oidc configuration optional

	   	links = append(links, webfinger.Link{
	   		Rel:  OpenIDConnectRel,
	   		Href: href,
	   		Titles: map[string]string{
	   			"en": "ownCloud Infinite Scale OpenID Connect Identity Provider",
	   		},
	   	})

	   	for _, instanceId := range instanceIds {
	   		instance := s.instanceLookup.GetInstance(ctx, instanceId)
	   		links = append(links, webfinger.Link{
	   			Rel:    OwnCloudInstanceRel,
	   			Href:   instance.Href,
	   			Titles: instance.Titles,
	   		})
	   	}

	   	return webfinger.JSONResourceDescriptor{
	   		Subject: queryTarget.String(),
	   		Links:   links,
	   	}, nil
	*/
}
