package middleware

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
)

const publicFilesEndpoint = "/remote.php/dav/public-files/"

// BasicAuth provides a middleware to check if BasicAuth is provided
func BasicAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	if options.EnableBasicAuth {
		options.Logger.Warn().Msg("basic auth enabled, use only for testing or development")
	}

	h := basicAuth{
		logger:       logger,
		enabled:      options.EnableBasicAuth,
		userProvider: options.UserProvider,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				if h.isPublicLink(req) || !h.isBasicAuth(req) {
					if !h.isPublicLink(req) {
						userAgentAuthenticateLockIn(w, req, options.CredentialsByUserAgent, "basic")
					}
					next.ServeHTTP(w, req)
					return
				}

				removeSuperfluousAuthenticate(w)
				login, password, _ := req.BasicAuth()
				user, err := h.userProvider.Authenticate(req.Context(), login, password)

				// touch is a user agent locking guard, when touched changes to true it indicates the User-Agent on the
				// request is configured to support only one challenge, it it remains untouched, there are no considera-
				// tions and we should write all available authentication challenges to the response.
				touch := false

				if err != nil {
					for k, v := range options.CredentialsByUserAgent {
						if strings.Contains(k, req.UserAgent()) {
							removeSuperfluousAuthenticate(w)
							w.Header().Add("Www-Authenticate", fmt.Sprintf("%v realm=\"%s\", charset=\"UTF-8\"", strings.Title(v), req.Host))
							touch = true
							break
						}
					}

					// if the request is not bound to any user agent, write all available challenges
					if !touch {
						writeSupportedAuthenticateHeader(w, req)
					}

					// if the request is a PROPFIND return a WebDAV error code.
					// TODO: The proxy has to be smart enough to detect when a request is directed towards a webdav server
					// and react accordingly.

					w.WriteHeader(http.StatusUnauthorized)

					b, err := Marshal(exception{
						code:    SabredavPermissionDenied,
						message: "Authentication error",
					})

					HandleWebdavError(w, b, err)
					return
				}

				claims := &oidc.StandardClaims{
					OcisID:            user.Id.OpaqueId,
					Iss:               user.Id.Idp,
					PreferredUsername: user.Username,
					Email:             user.Mail,
				}

				next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), claims)))
			},
		)
	}
}

type basicAuth struct {
	logger       log.Logger
	enabled      bool
	userProvider backend.UserBackend
}

func (m basicAuth) isPublicLink(req *http.Request) bool {
	login, _, ok := req.BasicAuth()
	return ok && login == "public" && strings.HasPrefix(req.URL.Path, publicFilesEndpoint)
}

func (m basicAuth) isBasicAuth(req *http.Request) bool {
	_, _, ok := req.BasicAuth()
	return m.enabled && ok
}

type code int

const (
	// SabredavBadRequest maps to HTTP 400
	SabredavBadRequest code = iota
	// SabredavMethodNotAllowed maps to HTTP 405
	SabredavMethodNotAllowed
	// SabredavNotAuthenticated maps to HTTP 401
	SabredavNotAuthenticated
	// SabredavPreconditionFailed maps to HTTP 412
	SabredavPreconditionFailed
	// SabredavPermissionDenied maps to HTTP 403
	SabredavPermissionDenied
	// SabredavNotFound maps to HTTP 404
	SabredavNotFound
	// SabredavConflict maps to HTTP 409
	SabredavConflict
)

var (
	codesEnum = []string{
		"Sabre\\DAV\\Exception\\BadRequest",
		"Sabre\\DAV\\Exception\\MethodNotAllowed",
		"Sabre\\DAV\\Exception\\NotAuthenticated",
		"Sabre\\DAV\\Exception\\PreconditionFailed",
		"Sabre\\DAV\\Exception\\PermissionDenied",
		"Sabre\\DAV\\Exception\\NotFound",
		"Sabre\\DAV\\Exception\\Conflict",
	}
)

type exception struct {
	code    code
	message string
	header  string
}

// Marshal just calls the xml marshaller for a given exception.
func Marshal(e exception) ([]byte, error) {
	xmlstring, err := xml.Marshal(&errorXML{
		Xmlnsd:    "DAV",
		Xmlnss:    "http://sabredav.org/ns",
		Exception: codesEnum[e.code],
		Message:   e.message,
		Header:    e.header,
	})
	if err != nil {
		return []byte(""), err
	}
	return []byte(xml.Header + string(xmlstring)), err
}

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_error
type errorXML struct {
	XMLName   xml.Name `xml:"d:error"`
	Xmlnsd    string   `xml:"xmlns:d,attr"`
	Xmlnss    string   `xml:"xmlns:s,attr"`
	Exception string   `xml:"s:exception"`
	Message   string   `xml:"s:message"`
	InnerXML  []byte   `xml:",innerxml"`
	Header    string   `xml:"s:header,omitempty"`
}

func HandleWebdavError(w http.ResponseWriter, b []byte, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
	}
}
