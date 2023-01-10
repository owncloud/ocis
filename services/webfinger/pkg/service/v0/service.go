package service

import (
	"context"

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
	Webfinger(ctx context.Context, resource, rel string) (webfinger.JSONResourceDescriptor, error)
}

// New returns a new instance of Service
func New(opts ...Option) (Service, error) {
	options := newOptions(opts...)

	return svc{
		log:    options.Logger,
		config: options.Config,
	}, nil
}

type svc struct {
	config *config.Config
	log    log.Logger
}

// SpacesThumbnail is the endpoint for retrieving thumbnails inside of spaces.
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
func (s svc) Webfinger(ctx context.Context, resource, rel string) (webfinger.JSONResourceDescriptor, error) {

	// TODO query ldap server here and fetch all instances the user has access to
	// what is the domain for the instance?

	// TODO use another relation? more graph specific? nah
	return webfinger.JSONResourceDescriptor{
		Subject: resource,
		Links: []webfinger.Link{
			{
				Rel:  OwnCloudInstanceRel,
				Href: "https://instance.server...",
				Titles: map[string]string{
					"en": "Readable Instance name",
				},
			},
			{
				Rel:  OwnCloudInstanceRel,
				Href: "https://otherinstance.server...",
				Titles: map[string]string{
					"en": "Other readable Instance name",
				},
			},
			// and we can return the OpenID Connect
			{
				Rel:  OpenIDConnectRel,
				Href: "https://idp.server...",
				Titles: map[string]string{
					"en": "Readable Openid Connect IDP name",
				},
			},
			{
				Rel:  OpenIDConnectRel,
				Href: "https://otheridp.server...",
				Titles: map[string]string{
					"en": "Other readable Openid Connect IDP name",
				},
			},
			// FIXME but now the clients have no way of knowing whic idp belongs to which instance
			// we could mix like this:
			{
				Rel:  OwnCloudInstanceRel,
				Href: "https://otherinstance.server...",
				Titles: map[string]string{
					"en": "Other readable Instance name",
				},
				Properties: map[string]string{
					OpenIDConnectRel: "https://otheridp.server...",
				},
			},
		},
	}, nil
}
