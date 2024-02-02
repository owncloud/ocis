// Package gocloak is a golang keycloak adaptor.
package gocloak

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"

	"github.com/Nerzal/gocloak/v13/pkg/jwx"
)

// GoCloak provides functionalities to talk to Keycloak.
type GoCloak struct {
	basePath    string
	certsCache  sync.Map
	certsLock   sync.Mutex
	restyClient *resty.Client
	Config      struct {
		CertsInvalidateTime time.Duration
		authAdminRealms     string
		authRealms          string
		tokenEndpoint       string
		revokeEndpoint      string
		logoutEndpoint      string
		openIDConnect       string
		attackDetection     string
	}
}

const (
	adminClientID string = "admin-cli"
	urlSeparator  string = "/"
)

func makeURL(path ...string) string {
	return strings.Join(path, urlSeparator)
}

// GetRequest returns a request for calling endpoints.
func (g *GoCloak) GetRequest(ctx context.Context) *resty.Request {
	var err HTTPErrorResponse
	return injectTracingHeaders(
		ctx, g.restyClient.R().
			SetContext(ctx).
			SetError(&err),
	)
}

// GetRequestWithBearerAuthNoCache returns a JSON base request configured with an auth token and no-cache header.
func (g *GoCloak) GetRequestWithBearerAuthNoCache(ctx context.Context, token string) *resty.Request {
	return g.GetRequest(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Cache-Control", "no-cache")
}

// GetRequestWithBearerAuth returns a JSON base request configured with an auth token.
func (g *GoCloak) GetRequestWithBearerAuth(ctx context.Context, token string) *resty.Request {
	return g.GetRequest(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json")
}

// GetRequestWithBearerAuthXMLHeader returns an XML base request configured with an auth token.
func (g *GoCloak) GetRequestWithBearerAuthXMLHeader(ctx context.Context, token string) *resty.Request {
	return g.GetRequest(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/xml;charset=UTF-8")
}

// GetRequestWithBasicAuth returns a form data base request configured with basic auth.
func (g *GoCloak) GetRequestWithBasicAuth(ctx context.Context, clientID, clientSecret string) *resty.Request {
	req := g.GetRequest(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded")
	// Public client doesn't require Basic Auth
	if len(clientID) > 0 && len(clientSecret) > 0 {
		httpBasicAuth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
		req.SetHeader("Authorization", "Basic "+httpBasicAuth)
	}

	return req
}

func (g *GoCloak) getRequestingParty(ctx context.Context, token string, realm string, options RequestingPartyTokenOptions, res interface{}) (*resty.Response, error) {
	return g.GetRequestWithBearerAuth(ctx, token).
		SetFormData(options.FormData()).
		SetFormDataFromValues(url.Values{"permission": PStringSlice(options.Permissions)}).
		SetResult(&res).
		Post(g.getRealmURL(realm, g.Config.tokenEndpoint))
}

func checkForError(resp *resty.Response, err error, errMessage string) error {
	if err != nil {
		return &APIError{
			Code:    0,
			Message: errors.Wrap(err, errMessage).Error(),
			Type:    ParseAPIErrType(err),
		}
	}

	if resp == nil {
		return &APIError{
			Message: "empty response",
			Type:    ParseAPIErrType(err),
		}
	}

	if resp.IsError() {
		var msg string

		if e, ok := resp.Error().(*HTTPErrorResponse); ok && e.NotEmpty() {
			msg = fmt.Sprintf("%s: %s", resp.Status(), e)
		} else {
			msg = resp.Status()
		}

		return &APIError{
			Code:    resp.StatusCode(),
			Message: msg,
			Type:    ParseAPIErrType(err),
		}
	}

	return nil
}

func getID(resp *resty.Response) string {
	header := resp.Header().Get("Location")
	splittedPath := strings.Split(header, urlSeparator)

	return splittedPath[len(splittedPath)-1]
}

func findUsedKey(usedKeyID string, keys []CertResponseKey) *CertResponseKey {
	for _, key := range keys {
		if *(key.Kid) == usedKeyID {
			return &key
		}
	}

	return nil
}

func injectTracingHeaders(ctx context.Context, req *resty.Request) *resty.Request {
	// look for span in context, do nothing if span is not found
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return req
	}

	// look for tracer in context, use global tracer if not found
	tracer, ok := ctx.Value(tracerContextKey).(opentracing.Tracer)
	if !ok || tracer == nil {
		tracer = opentracing.GlobalTracer()
	}

	// inject tracing header into request
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return req
	}

	return req
}

// ===============
// Keycloak client
// ===============

// NewClient creates a new Client
func NewClient(basePath string, options ...func(*GoCloak)) *GoCloak {
	c := GoCloak{
		basePath:    strings.TrimRight(basePath, urlSeparator),
		restyClient: resty.New(),
	}

	c.Config.CertsInvalidateTime = 10 * time.Minute
	c.Config.authAdminRealms = makeURL("admin", "realms")
	c.Config.authRealms = makeURL("realms")
	c.Config.tokenEndpoint = makeURL("protocol", "openid-connect", "token")
	c.Config.logoutEndpoint = makeURL("protocol", "openid-connect", "logout")
	c.Config.revokeEndpoint = makeURL("protocol", "openid-connect", "revoke")
	c.Config.openIDConnect = makeURL("protocol", "openid-connect")
	c.Config.attackDetection = makeURL("attack-detection", "brute-force")

	for _, option := range options {
		option(&c)
	}

	return &c
}

// RestyClient returns the internal resty g.
// This can be used to configure the g.
func (g *GoCloak) RestyClient() *resty.Client {
	return g.restyClient
}

// SetRestyClient overwrites the internal resty g.
func (g *GoCloak) SetRestyClient(restyClient *resty.Client) {
	g.restyClient = restyClient
}

func (g *GoCloak) getRealmURL(realm string, path ...string) string {
	path = append([]string{g.basePath, g.Config.authRealms, realm}, path...)
	return makeURL(path...)
}

func (g *GoCloak) getAdminRealmURL(realm string, path ...string) string {
	path = append([]string{g.basePath, g.Config.authAdminRealms, realm}, path...)
	return makeURL(path...)
}

func (g *GoCloak) getAttackDetectionURL(realm string, user string, path ...string) string {
	path = append([]string{g.basePath, g.Config.authAdminRealms, realm, g.Config.attackDetection, user}, path...)
	return makeURL(path...)
}

// ==== Functional Options ===

// SetLegacyWildFlySupport maintain legacy WildFly support.
func SetLegacyWildFlySupport() func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.authAdminRealms = makeURL("auth", "admin", "realms")
		g.Config.authRealms = makeURL("auth", "realms")
	}
}

// SetAuthRealms sets the auth realm
func SetAuthRealms(url string) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.authRealms = url
	}
}

// SetAuthAdminRealms sets the auth admin realm
func SetAuthAdminRealms(url string) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.authAdminRealms = url
	}
}

// SetTokenEndpoint sets the token endpoint
func SetTokenEndpoint(url string) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.tokenEndpoint = url
	}
}

// SetRevokeEndpoint sets the revoke endpoint
func SetRevokeEndpoint(url string) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.revokeEndpoint = url
	}
}

// SetLogoutEndpoint sets the logout
func SetLogoutEndpoint(url string) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.logoutEndpoint = url
	}
}

// SetOpenIDConnectEndpoint sets the logout
func SetOpenIDConnectEndpoint(url string) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.openIDConnect = url
	}
}

// SetCertCacheInvalidationTime sets the logout
func SetCertCacheInvalidationTime(duration time.Duration) func(g *GoCloak) {
	return func(g *GoCloak) {
		g.Config.CertsInvalidateTime = duration
	}
}

// GetServerInfo fetches the server info.
func (g *GoCloak) GetServerInfo(ctx context.Context, accessToken string) (*ServerInfoRepresentation, error) {
	errMessage := "could not get server info"
	var result *ServerInfoRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(makeURL(g.basePath, "admin", "serverinfo"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetUserInfo calls the UserInfo endpoint
func (g *GoCloak) GetUserInfo(ctx context.Context, accessToken, realm string) (*UserInfo, error) {
	const errMessage = "could not get user info"

	var result UserInfo
	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getRealmURL(realm, g.Config.openIDConnect, "userinfo"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRawUserInfo calls the UserInfo endpoint and returns a raw json object
func (g *GoCloak) GetRawUserInfo(ctx context.Context, accessToken, realm string) (map[string]interface{}, error) {
	const errMessage = "could not get user info"

	var result map[string]interface{}
	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getRealmURL(realm, g.Config.openIDConnect, "userinfo"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

func (g *GoCloak) getNewCerts(ctx context.Context, realm string) (*CertResponse, error) {
	const errMessage = "could not get newCerts"

	var result CertResponse
	resp, err := g.GetRequest(ctx).
		SetResult(&result).
		Get(g.getRealmURL(realm, g.Config.openIDConnect, "certs"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCerts fetches certificates for the given realm from the public /open-id-connect/certs endpoint
func (g *GoCloak) GetCerts(ctx context.Context, realm string) (*CertResponse, error) {
	const errMessage = "could not get certs"

	if cert, ok := g.certsCache.Load(realm); ok {
		return cert.(*CertResponse), nil
	}

	g.certsLock.Lock()
	defer g.certsLock.Unlock()

	if cert, ok := g.certsCache.Load(realm); ok {
		return cert.(*CertResponse), nil
	}

	cert, err := g.getNewCerts(ctx, realm)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	g.certsCache.Store(realm, cert)
	time.AfterFunc(g.Config.CertsInvalidateTime, func() {
		g.certsCache.Delete(realm)
	})

	return cert, nil
}

// GetIssuer gets the issuer of the given realm
func (g *GoCloak) GetIssuer(ctx context.Context, realm string) (*IssuerResponse, error) {
	const errMessage = "could not get issuer"

	var result IssuerResponse
	resp, err := g.GetRequest(ctx).
		SetResult(&result).
		Get(g.getRealmURL(realm))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// RetrospectToken calls the openid-connect introspect endpoint
func (g *GoCloak) RetrospectToken(ctx context.Context, accessToken, clientID, clientSecret, realm string) (*IntroSpectTokenResult, error) {
	const errMessage = "could not introspect requesting party token"

	var result IntroSpectTokenResult
	resp, err := g.GetRequestWithBasicAuth(ctx, clientID, clientSecret).
		SetFormData(map[string]string{
			"token_type_hint": "requesting_party_token",
			"token":           accessToken,
		}).
		SetResult(&result).
		Post(g.getRealmURL(realm, g.Config.tokenEndpoint, "introspect"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

func (g *GoCloak) decodeAccessTokenWithClaims(ctx context.Context, accessToken, realm string, claims jwt.Claims) (*jwt.Token, error) {
	const errMessage = "could not decode access token"
	accessToken = strings.Replace(accessToken, "Bearer ", "", 1)

	decodedHeader, err := jwx.DecodeAccessTokenHeader(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	certResult, err := g.GetCerts(ctx, realm)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	if certResult.Keys == nil {
		return nil, errors.Wrap(errors.New("there is no keys to decode the token"), errMessage)
	}
	usedKey := findUsedKey(decodedHeader.Kid, *certResult.Keys)
	if usedKey == nil {
		return nil, errors.Wrap(errors.New("cannot find a key to decode the token"), errMessage)
	}

	if strings.HasPrefix(decodedHeader.Alg, "ES") {
		return jwx.DecodeAccessTokenECDSACustomClaims(accessToken, usedKey.X, usedKey.Y, usedKey.Crv, claims)
	} else if strings.HasPrefix(decodedHeader.Alg, "RS") {
		return jwx.DecodeAccessTokenRSACustomClaims(accessToken, usedKey.E, usedKey.N, claims)
	}
	return nil, fmt.Errorf("unsupported algorithm")
}

// DecodeAccessToken decodes the accessToken
func (g *GoCloak) DecodeAccessToken(ctx context.Context, accessToken, realm string) (*jwt.Token, *jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := g.decodeAccessTokenWithClaims(ctx, accessToken, realm, claims)
	if err != nil {
		return nil, nil, err
	}
	return token, &claims, nil
}

// DecodeAccessTokenCustomClaims decodes the accessToken and writes claims into the given claims
func (g *GoCloak) DecodeAccessTokenCustomClaims(ctx context.Context, accessToken, realm string, claims jwt.Claims) (*jwt.Token, error) {
	return g.decodeAccessTokenWithClaims(ctx, accessToken, realm, claims)
}

// GetToken uses TokenOptions to fetch a token.
func (g *GoCloak) GetToken(ctx context.Context, realm string, options TokenOptions) (*JWT, error) {
	const errMessage = "could not get token"

	var token JWT
	var req *resty.Request

	if !NilOrEmpty(options.ClientSecret) {
		req = g.GetRequestWithBasicAuth(ctx, *options.ClientID, *options.ClientSecret)
	} else {
		req = g.GetRequest(ctx)
	}

	resp, err := req.SetFormData(options.FormData()).
		SetResult(&token).
		Post(g.getRealmURL(realm, g.Config.tokenEndpoint))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &token, nil
}

// GetRequestingPartyToken returns a requesting party token with permissions granted by the server
func (g *GoCloak) GetRequestingPartyToken(ctx context.Context, token, realm string, options RequestingPartyTokenOptions) (*JWT, error) {
	const errMessage = "could not get requesting party token"

	var res JWT

	resp, err := g.getRequestingParty(ctx, token, realm, options, &res)
	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetRequestingPartyPermissions returns a requesting party permissions granted by the server
func (g *GoCloak) GetRequestingPartyPermissions(ctx context.Context, token, realm string, options RequestingPartyTokenOptions) (*[]RequestingPartyPermission, error) {
	const errMessage = "could not get requesting party token"

	var res []RequestingPartyPermission

	options.ResponseMode = StringP("permissions")

	resp, err := g.getRequestingParty(ctx, token, realm, options, &res)
	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetRequestingPartyPermissionDecision returns a requesting party permission decision granted by the server
func (g *GoCloak) GetRequestingPartyPermissionDecision(ctx context.Context, token, realm string, options RequestingPartyTokenOptions) (*RequestingPartyPermissionDecision, error) {
	const errMessage = "could not get requesting party token"

	var res RequestingPartyPermissionDecision

	options.ResponseMode = StringP("decision")

	resp, err := g.getRequestingParty(ctx, token, realm, options, &res)
	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &res, nil
}

// RefreshToken refreshes the given token.
// May return a *APIError with further details about the issue.
func (g *GoCloak) RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*JWT, error) {
	return g.GetToken(ctx, realm, TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		GrantType:    StringP("refresh_token"),
		RefreshToken: &refreshToken,
	})
}

// LoginAdmin performs a login with Admin client
func (g *GoCloak) LoginAdmin(ctx context.Context, username, password, realm string) (*JWT, error) {
	return g.GetToken(ctx, realm, TokenOptions{
		ClientID:  StringP(adminClientID),
		GrantType: StringP("password"),
		Username:  &username,
		Password:  &password,
	})
}

// LoginClient performs a login with client credentials
func (g *GoCloak) LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*JWT, error) {
	opts := TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		GrantType:    StringP("client_credentials"),
	}

	if len(scopes) > 0 {
		opts.Scope = &scopes[0]
	}

	return g.GetToken(ctx, realm, opts)
}

// LoginClientTokenExchange will exchange the presented token for a user's token
// Requires Token-Exchange is enabled: https://www.keycloak.org/docs/latest/securing_apps/index.html#_token-exchange
func (g *GoCloak) LoginClientTokenExchange(ctx context.Context, clientID, token, clientSecret, realm, targetClient, userID string) (*JWT, error) {
	tokenOptions := TokenOptions{
		ClientID:           &clientID,
		ClientSecret:       &clientSecret,
		GrantType:          StringP("urn:ietf:params:oauth:grant-type:token-exchange"),
		SubjectToken:       &token,
		RequestedTokenType: StringP("urn:ietf:params:oauth:token-type:refresh_token"),
		Audience:           &targetClient,
	}
	if userID != "" {
		tokenOptions.RequestedSubject = &userID
	}
	return g.GetToken(ctx, realm, tokenOptions)
}

// LoginClientSignedJWT performs a login with client credentials and signed jwt claims
func (g *GoCloak) LoginClientSignedJWT(
	ctx context.Context,
	clientID,
	realm string,
	key interface{},
	signedMethod jwt.SigningMethod,
	expiresAt *jwt.NumericDate,
) (*JWT, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: expiresAt,
		Issuer:    clientID,
		Subject:   clientID,
		ID:        ksuid.New().String(),
		Audience: jwt.ClaimStrings{
			g.getRealmURL(realm),
		},
	}
	assertion, err := jwx.SignClaims(claims, key, signedMethod)
	if err != nil {
		return nil, err
	}

	return g.GetToken(ctx, realm, TokenOptions{
		ClientID:            &clientID,
		GrantType:           StringP("client_credentials"),
		ClientAssertionType: StringP("urn:ietf:params:oauth:client-assertion-type:jwt-bearer"),
		ClientAssertion:     &assertion,
	})
}

// Login performs a login with user credentials and a client
func (g *GoCloak) Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*JWT, error) {
	return g.GetToken(ctx, realm, TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		GrantType:    StringP("password"),
		Username:     &username,
		Password:     &password,
		Scope:        StringP("openid"),
	})
}

// LoginOtp performs a login with user credentials and otp token
func (g *GoCloak) LoginOtp(ctx context.Context, clientID, clientSecret, realm, username, password, totp string) (*JWT, error) {
	return g.GetToken(ctx, realm, TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		GrantType:    StringP("password"),
		Username:     &username,
		Password:     &password,
		Totp:         &totp,
	})
}

// Logout logs out users with refresh token
func (g *GoCloak) Logout(ctx context.Context, clientID, clientSecret, realm, refreshToken string) error {
	const errMessage = "could not logout"

	resp, err := g.GetRequestWithBasicAuth(ctx, clientID, clientSecret).
		SetFormData(map[string]string{
			"client_id":     clientID,
			"refresh_token": refreshToken,
		}).
		Post(g.getRealmURL(realm, g.Config.logoutEndpoint))

	return checkForError(resp, err, errMessage)
}

// LogoutPublicClient performs a logout using a public client and the accessToken.
func (g *GoCloak) LogoutPublicClient(ctx context.Context, clientID, realm, accessToken, refreshToken string) error {
	const errMessage = "could not logout public client"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetFormData(map[string]string{
			"client_id":     clientID,
			"refresh_token": refreshToken,
		}).
		Post(g.getRealmURL(realm, g.Config.logoutEndpoint))

	return checkForError(resp, err, errMessage)
}

// LogoutAllSessions logs out all sessions of a user given an id.
func (g *GoCloak) LogoutAllSessions(ctx context.Context, accessToken, realm, userID string) error {
	const errMessage = "could not logout"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		Post(g.getAdminRealmURL(realm, "users", userID, "logout"))

	return checkForError(resp, err, errMessage)
}

// RevokeUserConsents revokes the given user consent.
func (g *GoCloak) RevokeUserConsents(ctx context.Context, accessToken, realm, userID, clientID string) error {
	const errMessage = "could not revoke consents"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		Delete(g.getAdminRealmURL(realm, "users", userID, "consents", clientID))

	return checkForError(resp, err, errMessage)
}

// LogoutUserSession logs out a single sessions of a user given a session id
func (g *GoCloak) LogoutUserSession(ctx context.Context, accessToken, realm, session string) error {
	const errMessage = "could not logout"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		Delete(g.getAdminRealmURL(realm, "sessions", session))

	return checkForError(resp, err, errMessage)
}

// ExecuteActionsEmail executes an actions email
func (g *GoCloak) ExecuteActionsEmail(ctx context.Context, token, realm string, params ExecuteActionsEmail) error {
	const errMessage = "could not execute actions email"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(params.Actions).
		SetQueryParams(queryParams).
		Put(g.getAdminRealmURL(realm, "users", *(params.UserID), "execute-actions-email"))

	return checkForError(resp, err, errMessage)
}

// SendVerifyEmail sends a verification e-mail to a user.
func (g *GoCloak) SendVerifyEmail(ctx context.Context, token, userID, realm string, params ...SendVerificationMailParams) error {
	const errMessage = "could not execute actions email"

	queryParams := map[string]string{}
	if params != nil {
		if params[0].ClientID != nil {
			queryParams["client_id"] = *params[0].ClientID
		}

		if params[0].RedirectURI != nil {
			queryParams["redirect_uri"] = *params[0].RedirectURI
		}
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetQueryParams(queryParams).
		Put(g.getAdminRealmURL(realm, "users", userID, "send-verify-email"))

	return checkForError(resp, err, errMessage)
}

// CreateGroup creates a new group.
func (g *GoCloak) CreateGroup(ctx context.Context, token, realm string, group Group) (string, error) {
	const errMessage = "could not create group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(group).
		Post(g.getAdminRealmURL(realm, "groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}
	return getID(resp), nil
}

// CreateChildGroup creates a new child group
func (g *GoCloak) CreateChildGroup(ctx context.Context, token, realm, groupID string, group Group) (string, error) {
	const errMessage = "could not create child group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(group).
		Post(g.getAdminRealmURL(realm, "groups", groupID, "children"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// CreateComponent creates the given component.
func (g *GoCloak) CreateComponent(ctx context.Context, token, realm string, component Component) (string, error) {
	const errMessage = "could not create component"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(component).
		Post(g.getAdminRealmURL(realm, "components"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// CreateClient creates the given g.
func (g *GoCloak) CreateClient(ctx context.Context, accessToken, realm string, newClient Client) (string, error) {
	const errMessage = "could not create client"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetBody(newClient).
		Post(g.getAdminRealmURL(realm, "clients"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// CreateClientRepresentation creates a new client representation
func (g *GoCloak) CreateClientRepresentation(ctx context.Context, token, realm string, newClient Client) (*Client, error) {
	const errMessage = "could not create client representation"

	var result Client

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(newClient).
		Post(g.getRealmURL(realm, "clients-registrations", "default"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateClientRole creates a new role for a client
func (g *GoCloak) CreateClientRole(ctx context.Context, token, realm, idOfClient string, role Role) (string, error) {
	const errMessage = "could not create client role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "roles"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// CreateClientScope creates a new client scope
func (g *GoCloak) CreateClientScope(ctx context.Context, token, realm string, scope ClientScope) (string, error) {
	const errMessage = "could not create client scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(scope).
		Post(g.getAdminRealmURL(realm, "client-scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// CreateClientScopeProtocolMapper creates a new protocolMapper under the given client scope
func (g *GoCloak) CreateClientScopeProtocolMapper(ctx context.Context, token, realm, scopeID string, protocolMapper ProtocolMappers) (string, error) {
	const errMessage = "could not create client scope protocol mapper"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(protocolMapper).
		Post(g.getAdminRealmURL(realm, "client-scopes", scopeID, "protocol-mappers", "models"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// UpdateGroup updates the given group.
func (g *GoCloak) UpdateGroup(ctx context.Context, token, realm string, updatedGroup Group) error {
	const errMessage = "could not update group"

	if NilOrEmpty(updatedGroup.ID) {
		return errors.Wrap(errors.New("ID of a group required"), errMessage)
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(updatedGroup).
		Put(g.getAdminRealmURL(realm, "groups", PString(updatedGroup.ID)))

	return checkForError(resp, err, errMessage)
}

// UpdateGroupManagementPermissions updates the given group management permissions
func (g *GoCloak) UpdateGroupManagementPermissions(ctx context.Context, accessToken, realm string, idOfGroup string, managementPermissions ManagementPermissionRepresentation) (*ManagementPermissionRepresentation, error) {
	const errMessage = "could not update group management permissions"

	var result ManagementPermissionRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		SetBody(managementPermissions).
		Put(g.getAdminRealmURL(realm, "groups", idOfGroup, "management", "permissions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateClient updates the given Client
func (g *GoCloak) UpdateClient(ctx context.Context, token, realm string, updatedClient Client) error {
	const errMessage = "could not update client"

	if NilOrEmpty(updatedClient.ID) {
		return errors.Wrap(errors.New("ID of a client required"), errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(updatedClient).
		Put(g.getAdminRealmURL(realm, "clients", PString(updatedClient.ID)))

	return checkForError(resp, err, errMessage)
}

// UpdateClientRepresentation updates the given client representation
func (g *GoCloak) UpdateClientRepresentation(ctx context.Context, accessToken, realm string, updatedClient Client) (*Client, error) {
	const errMessage = "could not update client representation"

	if NilOrEmpty(updatedClient.ID) {
		return nil, errors.Wrap(errors.New("ID of a client required"), errMessage)
	}

	var result Client

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		SetBody(updatedClient).
		Put(g.getRealmURL(realm, "clients-registrations", "default", PString(updatedClient.ClientID)))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateClientManagementPermissions updates the given client management permissions
func (g *GoCloak) UpdateClientManagementPermissions(ctx context.Context, accessToken, realm string, idOfClient string, managementPermissions ManagementPermissionRepresentation) (*ManagementPermissionRepresentation, error) {
	const errMessage = "could not update client management permissions"

	var result ManagementPermissionRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		SetBody(managementPermissions).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "management", "permissions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateRole updates the given role.
func (g *GoCloak) UpdateRole(ctx context.Context, token, realm, idOfClient string, role Role) error {
	const errMessage = "could not update role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "roles", PString(role.Name)))

	return checkForError(resp, err, errMessage)
}

// UpdateClientScope updates the given client scope.
func (g *GoCloak) UpdateClientScope(ctx context.Context, token, realm string, scope ClientScope) error {
	const errMessage = "could not update client scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(scope).
		Put(g.getAdminRealmURL(realm, "client-scopes", PString(scope.ID)))

	return checkForError(resp, err, errMessage)
}

// UpdateClientScopeProtocolMapper updates the given protocol mapper for a client scope
func (g *GoCloak) UpdateClientScopeProtocolMapper(ctx context.Context, token, realm, scopeID string, protocolMapper ProtocolMappers) error {
	const errMessage = "could not update client scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(protocolMapper).
		Put(g.getAdminRealmURL(realm, "client-scopes", scopeID, "protocol-mappers", "models", PString(protocolMapper.ID)))

	return checkForError(resp, err, errMessage)
}

// DeleteGroup deletes the group with the given groupID.
func (g *GoCloak) DeleteGroup(ctx context.Context, token, realm, groupID string) error {
	const errMessage = "could not delete group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "groups", groupID))

	return checkForError(resp, err, errMessage)
}

// DeleteClient deletes a given client
func (g *GoCloak) DeleteClient(ctx context.Context, token, realm, idOfClient string) error {
	const errMessage = "could not delete client"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// DeleteComponent deletes the component with the given id.
func (g *GoCloak) DeleteComponent(ctx context.Context, token, realm, componentID string) error {
	const errMessage = "could not delete component"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "components", componentID))

	return checkForError(resp, err, errMessage)
}

// DeleteClientRepresentation deletes a given client representation.
func (g *GoCloak) DeleteClientRepresentation(ctx context.Context, accessToken, realm, clientID string) error {
	const errMessage = "could not delete client representation"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		Delete(g.getRealmURL(realm, "clients-registrations", "default", clientID))

	return checkForError(resp, err, errMessage)
}

// DeleteClientRole deletes a given role.
func (g *GoCloak) DeleteClientRole(ctx context.Context, token, realm, idOfClient, roleName string) error {
	const errMessage = "could not delete client role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "roles", roleName))

	return checkForError(resp, err, errMessage)
}

// DeleteClientScope deletes the scope with the given id.
func (g *GoCloak) DeleteClientScope(ctx context.Context, token, realm, scopeID string) error {
	const errMessage = "could not delete client scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "client-scopes", scopeID))

	return checkForError(resp, err, errMessage)
}

// DeleteClientScopeProtocolMapper deletes the given protocol mapper from the client scope
func (g *GoCloak) DeleteClientScopeProtocolMapper(ctx context.Context, token, realm, scopeID, protocolMapperID string) error {
	const errMessage = "could not delete client scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "client-scopes", scopeID, "protocol-mappers", "models", protocolMapperID))

	return checkForError(resp, err, errMessage)
}

// GetClient returns a client
func (g *GoCloak) GetClient(ctx context.Context, token, realm, idOfClient string) (*Client, error) {
	const errMessage = "could not get client"

	var result Client

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientRepresentation returns a client representation
func (g *GoCloak) GetClientRepresentation(ctx context.Context, accessToken, realm, clientID string) (*Client, error) {
	const errMessage = "could not get client representation"

	var result Client

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getRealmURL(realm, "clients-registrations", "default", clientID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAdapterConfiguration returns a adapter configuration
func (g *GoCloak) GetAdapterConfiguration(ctx context.Context, accessToken, realm, clientID string) (*AdapterConfiguration, error) {
	const errMessage = "could not get adapter configuration"

	var result AdapterConfiguration

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getRealmURL(realm, "clients-registrations", "install", clientID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientsDefaultScopes returns a list of the client's default scopes
func (g *GoCloak) GetClientsDefaultScopes(ctx context.Context, token, realm, idOfClient string) ([]*ClientScope, error) {
	const errMessage = "could not get clients default scopes"

	var result []*ClientScope

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "default-client-scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// AddDefaultScopeToClient adds a client scope to the list of client's default scopes
func (g *GoCloak) AddDefaultScopeToClient(ctx context.Context, token, realm, idOfClient, scopeID string) error {
	const errMessage = "could not add default scope to client"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "default-client-scopes", scopeID))

	return checkForError(resp, err, errMessage)
}

// RemoveDefaultScopeFromClient removes a client scope from the list of client's default scopes
func (g *GoCloak) RemoveDefaultScopeFromClient(ctx context.Context, token, realm, idOfClient, scopeID string) error {
	const errMessage = "could not remove default scope from client"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "default-client-scopes", scopeID))

	return checkForError(resp, err, errMessage)
}

// GetClientsOptionalScopes returns a list of the client's optional scopes
func (g *GoCloak) GetClientsOptionalScopes(ctx context.Context, token, realm, idOfClient string) ([]*ClientScope, error) {
	const errMessage = "could not get clients optional scopes"

	var result []*ClientScope

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "optional-client-scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// AddOptionalScopeToClient adds a client scope to the list of client's optional scopes
func (g *GoCloak) AddOptionalScopeToClient(ctx context.Context, token, realm, idOfClient, scopeID string) error {
	const errMessage = "could not add optional scope to client"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "optional-client-scopes", scopeID))

	return checkForError(resp, err, errMessage)
}

// RemoveOptionalScopeFromClient deletes a client scope from the list of client's optional scopes
func (g *GoCloak) RemoveOptionalScopeFromClient(ctx context.Context, token, realm, idOfClient, scopeID string) error {
	const errMessage = "could not remove optional scope from client"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "optional-client-scopes", scopeID))

	return checkForError(resp, err, errMessage)
}

// GetDefaultOptionalClientScopes returns a list of default realm optional scopes
func (g *GoCloak) GetDefaultOptionalClientScopes(ctx context.Context, token, realm string) ([]*ClientScope, error) {
	const errMessage = "could not get default optional client scopes"

	var result []*ClientScope

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "default-optional-client-scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetDefaultDefaultClientScopes returns a list of default realm default scopes
func (g *GoCloak) GetDefaultDefaultClientScopes(ctx context.Context, token, realm string) ([]*ClientScope, error) {
	const errMessage = "could not get default client scopes"

	var result []*ClientScope

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "default-default-client-scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScope returns a clientscope
func (g *GoCloak) GetClientScope(ctx context.Context, token, realm, scopeID string) (*ClientScope, error) {
	const errMessage = "could not get client scope"

	var result ClientScope

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", scopeID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientScopes returns all client scopes
func (g *GoCloak) GetClientScopes(ctx context.Context, token, realm string) ([]*ClientScope, error) {
	const errMessage = "could not get client scopes"

	var result []*ClientScope

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopeProtocolMappers returns all protocol mappers of a client scope
func (g *GoCloak) GetClientScopeProtocolMappers(ctx context.Context, token, realm, scopeID string) ([]*ProtocolMappers, error) {
	const errMessage = "could not get client scope protocol mappers"

	var result []*ProtocolMappers

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", scopeID, "protocol-mappers", "models"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopeProtocolMapper returns a protocol mapper of a client scope
func (g *GoCloak) GetClientScopeProtocolMapper(ctx context.Context, token, realm, scopeID, protocolMapperID string) (*ProtocolMappers, error) {
	const errMessage = "could not get client scope protocol mappers"

	var result *ProtocolMappers

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", scopeID, "protocol-mappers", "models", protocolMapperID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopeMappings returns all scope mappings for the client
func (g *GoCloak) GetClientScopeMappings(ctx context.Context, token, realm, idOfClient string) (*MappingsRepresentation, error) {
	const errMessage = "could not get all scope mappings for the client"

	var result *MappingsRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopeMappingsRealmRoles returns realm-level roles associated with the client’s scope
func (g *GoCloak) GetClientScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClient string) ([]*Role, error) {
	const errMessage = "could not get realm-level roles with the client’s scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "realm"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopeMappingsRealmRolesAvailable returns realm-level roles that are available to attach to this client’s scope
func (g *GoCloak) GetClientScopeMappingsRealmRolesAvailable(ctx context.Context, token, realm, idOfClient string) ([]*Role, error) {
	const errMessage = "could not get available realm-level roles with the client’s scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "realm", "available"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateClientScopeMappingsRealmRoles create realm-level roles to the client’s scope
func (g *GoCloak) CreateClientScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClient string, roles []Role) error {
	const errMessage = "could not create realm-level roles to the client’s scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// DeleteClientScopeMappingsRealmRoles deletes realm-level roles from the client’s scope
func (g *GoCloak) DeleteClientScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClient string, roles []Role) error {
	const errMessage = "could not delete realm-level roles from the client’s scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// GetClientScopeMappingsClientRoles returns roles associated with a client’s scope
func (g *GoCloak) GetClientScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string) ([]*Role, error) {
	const errMessage = "could not get roles associated with a client’s scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "clients", idOfSelectedClient))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopeMappingsClientRolesAvailable returns available roles associated with a client’s scope
func (g *GoCloak) GetClientScopeMappingsClientRolesAvailable(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string) ([]*Role, error) {
	const errMessage = "could not get available roles associated with a client’s scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "clients", idOfSelectedClient, "available"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateClientScopeMappingsClientRoles creates client-level roles from the client’s scope
func (g *GoCloak) CreateClientScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string, roles []Role) error {
	const errMessage = "could not create client-level roles from the client’s scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "clients", idOfSelectedClient))

	return checkForError(resp, err, errMessage)
}

// DeleteClientScopeMappingsClientRoles deletes client-level roles from the client’s scope
func (g *GoCloak) DeleteClientScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string, roles []Role) error {
	const errMessage = "could not delete client-level roles from the client’s scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "scope-mappings", "clients", idOfSelectedClient))

	return checkForError(resp, err, errMessage)
}

// GetClientSecret returns a client's secret
func (g *GoCloak) GetClientSecret(ctx context.Context, token, realm, idOfClient string) (*CredentialRepresentation, error) {
	const errMessage = "could not get client secret"

	var result CredentialRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "client-secret"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientServiceAccount retrieves the service account "user" for a client if enabled
func (g *GoCloak) GetClientServiceAccount(ctx context.Context, token, realm, idOfClient string) (*User, error) {
	const errMessage = "could not get client service account"

	var result User
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "service-account-user"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// RegenerateClientSecret triggers the creation of the new client secret.
func (g *GoCloak) RegenerateClientSecret(ctx context.Context, token, realm, idOfClient string) (*CredentialRepresentation, error) {
	const errMessage = "could not regenerate client secret"

	var result CredentialRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "client-secret"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientOfflineSessions returns offline sessions associated with the client
func (g *GoCloak) GetClientOfflineSessions(ctx context.Context, token, realm, idOfClient string, params ...GetClientUserSessionsParams) ([]*UserSessionRepresentation, error) {
	const errMessage = "could not get client offline sessions"
	var res []*UserSessionRepresentation

	queryParams := map[string]string{}
	if len(params) > 0 {
		var err error

		queryParams, err = GetQueryParams(params[0])
		if err != nil {
			return nil, errors.Wrap(err, errMessage)
		}
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "offline-sessions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return res, nil
}

// GetClientUserSessions returns user sessions associated with the client
func (g *GoCloak) GetClientUserSessions(ctx context.Context, token, realm, idOfClient string, params ...GetClientUserSessionsParams) ([]*UserSessionRepresentation, error) {
	const errMessage = "could not get client user sessions"
	var res []*UserSessionRepresentation

	queryParams := map[string]string{}
	if len(params) > 0 {
		var err error

		queryParams, err = GetQueryParams(params[0])
		if err != nil {
			return nil, errors.Wrap(err, errMessage)
		}
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "user-sessions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateClientProtocolMapper creates a protocol mapper in client scope
func (g *GoCloak) CreateClientProtocolMapper(ctx context.Context, token, realm, idOfClient string, mapper ProtocolMapperRepresentation) (string, error) {
	const errMessage = "could not create client protocol mapper"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(mapper).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "protocol-mappers", "models"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// UpdateClientProtocolMapper updates a protocol mapper in client scope
func (g *GoCloak) UpdateClientProtocolMapper(ctx context.Context, token, realm, idOfClient, mapperID string, mapper ProtocolMapperRepresentation) error {
	const errMessage = "could not update client protocol mapper"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(mapper).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "protocol-mappers", "models", mapperID))

	return checkForError(resp, err, errMessage)
}

// DeleteClientProtocolMapper deletes a protocol mapper in client scope
func (g *GoCloak) DeleteClientProtocolMapper(ctx context.Context, token, realm, idOfClient, mapperID string) error {
	const errMessage = "could not delete client protocol mapper"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "protocol-mappers", "models", mapperID))

	return checkForError(resp, err, errMessage)
}

// GetKeyStoreConfig get keystoreconfig of the realm
func (g *GoCloak) GetKeyStoreConfig(ctx context.Context, token, realm string) (*KeyStoreConfig, error) {
	const errMessage = "could not get key store config"

	var result KeyStoreConfig
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "keys"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetComponents get all components in realm
func (g *GoCloak) GetComponents(ctx context.Context, token, realm string) ([]*Component, error) {
	const errMessage = "could not get components"

	var result []*Component
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "components"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetComponentsWithParams get all components in realm with query params
func (g *GoCloak) GetComponentsWithParams(ctx context.Context, token, realm string, params GetComponentsParams) ([]*Component, error) {
	const errMessage = "could not get components"
	var result []*Component

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "components"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetComponent get exactly one component by ID
func (g *GoCloak) GetComponent(ctx context.Context, token, realm string, componentID string) (*Component, error) {
	const errMessage = "could not get components"
	var result *Component

	componentURL := fmt.Sprintf("components/%s", componentID)

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, componentURL))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateComponent updates the given component
func (g *GoCloak) UpdateComponent(ctx context.Context, token, realm string, component Component) error {
	const errMessage = "could not update component"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(component).
		Put(g.getAdminRealmURL(realm, "components", PString(component.ID)))

	return checkForError(resp, err, errMessage)
}

// GetDefaultGroups returns a list of default groups
func (g *GoCloak) GetDefaultGroups(ctx context.Context, token, realm string) ([]*Group, error) {
	const errMessage = "could not get default groups"

	var result []*Group

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "default-groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// AddDefaultGroup adds group to the list of default groups
func (g *GoCloak) AddDefaultGroup(ctx context.Context, token, realm, groupID string) error {
	const errMessage = "could not add default group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Put(g.getAdminRealmURL(realm, "default-groups", groupID))

	return checkForError(resp, err, errMessage)
}

// RemoveDefaultGroup removes group from the list of default groups
func (g *GoCloak) RemoveDefaultGroup(ctx context.Context, token, realm, groupID string) error {
	const errMessage = "could not remove default group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "default-groups", groupID))

	return checkForError(resp, err, errMessage)
}

func (g *GoCloak) getRoleMappings(ctx context.Context, token, realm, path, objectID string) (*MappingsRepresentation, error) {
	const errMessage = "could not get role mappings"

	var result MappingsRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, path, objectID, "role-mappings"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRoleMappingByGroupID gets the role mappings by group
func (g *GoCloak) GetRoleMappingByGroupID(ctx context.Context, token, realm, groupID string) (*MappingsRepresentation, error) {
	return g.getRoleMappings(ctx, token, realm, "groups", groupID)
}

// GetRoleMappingByUserID gets the role mappings by user
func (g *GoCloak) GetRoleMappingByUserID(ctx context.Context, token, realm, userID string) (*MappingsRepresentation, error) {
	return g.getRoleMappings(ctx, token, realm, "users", userID)
}

// GetGroup get group with id in realm
func (g *GoCloak) GetGroup(ctx context.Context, token, realm, groupID string) (*Group, error) {
	const errMessage = "could not get group"

	var result Group

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetGroupByPath get group with path in realm
func (g *GoCloak) GetGroupByPath(ctx context.Context, token, realm, groupPath string) (*Group, error) {
	const errMessage = "could not get group"

	var result Group

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "group-by-path", groupPath))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetGroups get all groups in realm
func (g *GoCloak) GetGroups(ctx context.Context, token, realm string, params GetGroupsParams) ([]*Group, error) {
	const errMessage = "could not get groups"

	var result []*Group
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGroupManagementPermissions returns whether group Authorization permissions have been initialized or not and a reference
// to the managed permissions
func (g *GoCloak) GetGroupManagementPermissions(ctx context.Context, token, realm string, idOfGroup string) (*ManagementPermissionRepresentation, error) {
	const errMessage = "could not get management permissions"

	var result ManagementPermissionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", idOfGroup, "management", "permissions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetGroupsByRole gets groups assigned with a specific role of a realm
func (g *GoCloak) GetGroupsByRole(ctx context.Context, token, realm string, roleName string) ([]*Group, error) {
	const errMessage = "could not get groups"

	var result []*Group
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles", roleName, "groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGroupsByClientRole gets groups with specified roles assigned of given client within a realm
func (g *GoCloak) GetGroupsByClientRole(ctx context.Context, token, realm string, roleName string, clientID string) ([]*Group, error) {
	const errMessage = "could not get groups"

	var result []*Group
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", clientID, "roles", roleName, "groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGroupsCount gets the groups count in the realm
func (g *GoCloak) GetGroupsCount(ctx context.Context, token, realm string, params GetGroupsParams) (int, error) {
	const errMessage = "could not get groups count"

	var result GroupsCount
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return 0, errors.Wrap(err, errMessage)
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "groups", "count"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return -1, errors.Wrap(err, errMessage)
	}

	return result.Count, nil
}

// GetGroupMembers get a list of users of group with id in realm
func (g *GoCloak) GetGroupMembers(ctx context.Context, token, realm, groupID string, params GetGroupsParams) ([]*User, error) {
	const errMessage = "could not get group members"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "members"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientRoles get all roles for the given client in realm
func (g *GoCloak) GetClientRoles(ctx context.Context, token, realm, idOfClient string, params GetRoleParams) ([]*Role, error) {
	const errMessage = "could not get client roles"

	var result []*Role
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "roles"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientRoleByID gets role for the given client in realm using role ID
func (g *GoCloak) GetClientRoleByID(ctx context.Context, token, realm, roleID string) (*Role, error) {
	const errMessage = "could not get client role"

	var result Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClientRolesByUserID returns all client roles assigned to the given user
func (g *GoCloak) GetClientRolesByUserID(ctx context.Context, token, realm, idOfClient, userID string) ([]*Role, error) {
	const errMessage = "could not client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientRolesByGroupID returns all client roles assigned to the given group
func (g *GoCloak) GetClientRolesByGroupID(ctx context.Context, token, realm, idOfClient, groupID string) ([]*Role, error) {
	const errMessage = "could not get client roles by group id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "clients", idOfClient))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeClientRolesByRoleID returns all client composite roles associated with the given client role
func (g *GoCloak) GetCompositeClientRolesByRoleID(ctx context.Context, token, realm, idOfClient, roleID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by role id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites", "clients", idOfClient))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeClientRolesByUserID returns all client roles and composite roles assigned to the given user
func (g *GoCloak) GetCompositeClientRolesByUserID(ctx context.Context, token, realm, idOfClient, userID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient, "composite"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAvailableClientRolesByUserID returns all available client roles to the given user
func (g *GoCloak) GetAvailableClientRolesByUserID(ctx context.Context, token, realm, idOfClient, userID string) ([]*Role, error) {
	const errMessage = "could not get available client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient, "available"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAvailableClientRolesByGroupID returns all available roles to the given group
func (g *GoCloak) GetAvailableClientRolesByGroupID(ctx context.Context, token, realm, idOfClient, groupID string) ([]*Role, error) {
	const errMessage = "could not get available client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "clients", idOfClient, "available"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeClientRolesByGroupID returns all client roles and composite roles assigned to the given group
func (g *GoCloak) GetCompositeClientRolesByGroupID(ctx context.Context, token, realm, idOfClient, groupID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by group id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "clients", idOfClient, "composite"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientRole get a role for the given client in a realm by role name
func (g *GoCloak) GetClientRole(ctx context.Context, token, realm, idOfClient, roleName string) (*Role, error) {
	const errMessage = "could not get client role"

	var result Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "roles", roleName))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetClients gets all clients in realm
func (g *GoCloak) GetClients(ctx context.Context, token, realm string, params GetClientsParams) ([]*Client, error) {
	const errMessage = "could not get clients"

	var result []*Client
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientManagementPermissions returns whether client Authorization permissions have been initialized or not and a reference
// to the managed permissions
func (g *GoCloak) GetClientManagementPermissions(ctx context.Context, token, realm string, idOfClient string) (*ManagementPermissionRepresentation, error) {
	const errMessage = "could not get management permissions"

	var result ManagementPermissionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "management", "permissions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UserAttributeContains checks if the given attribute value is set
func UserAttributeContains(attributes map[string][]string, attribute, value string) bool {
	for _, item := range attributes[attribute] {
		if item == value {
			return true
		}
	}
	return false
}

// -----------
// Realm Roles
// -----------

// CreateRealmRole creates a role in a realm
func (g *GoCloak) CreateRealmRole(ctx context.Context, token string, realm string, role Role) (string, error) {
	const errMessage = "could not create realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Post(g.getAdminRealmURL(realm, "roles"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// GetRealmRole returns a role from a realm by role's name
func (g *GoCloak) GetRealmRole(ctx context.Context, token, realm, roleName string) (*Role, error) {
	const errMessage = "could not get realm role"

	var result Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles", roleName))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRealmRoleByID returns a role from a realm by role's ID
func (g *GoCloak) GetRealmRoleByID(ctx context.Context, token, realm, roleID string) (*Role, error) {
	const errMessage = "could not get realm role"

	var result Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRealmRoles get all roles of the given realm.
func (g *GoCloak) GetRealmRoles(ctx context.Context, token, realm string, params GetRoleParams) ([]*Role, error) {
	const errMessage = "could not get realm roles"

	var result []*Role
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "roles"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetRealmRolesByUserID returns all roles assigned to the given user
func (g *GoCloak) GetRealmRolesByUserID(ctx context.Context, token, realm, userID string) ([]*Role, error) {
	const errMessage = "could not get realm roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetRealmRolesByGroupID returns all roles assigned to the given group
func (g *GoCloak) GetRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) ([]*Role, error) {
	const errMessage = "could not get realm roles by group id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateRealmRole updates a role in a realm
func (g *GoCloak) UpdateRealmRole(ctx context.Context, token, realm, roleName string, role Role) error {
	const errMessage = "could not update realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Put(g.getAdminRealmURL(realm, "roles", roleName))

	return checkForError(resp, err, errMessage)
}

// UpdateRealmRoleByID updates a role in a realm by role's ID
func (g *GoCloak) UpdateRealmRoleByID(ctx context.Context, token, realm, roleID string, role Role) error {
	const errMessage = "could not update realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(role).
		Put(g.getAdminRealmURL(realm, "roles-by-id", roleID))

	return checkForError(resp, err, errMessage)
}

// DeleteRealmRole deletes a role in a realm by role's name
func (g *GoCloak) DeleteRealmRole(ctx context.Context, token, realm, roleName string) error {
	const errMessage = "could not delete realm role"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "roles", roleName))

	return checkForError(resp, err, errMessage)
}

// AddRealmRoleToUser adds realm-level role mappings
func (g *GoCloak) AddRealmRoleToUser(ctx context.Context, token, realm, userID string, roles []Role) error {
	const errMessage = "could not add realm role to user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// DeleteRealmRoleFromUser deletes realm-level role mappings
func (g *GoCloak) DeleteRealmRoleFromUser(ctx context.Context, token, realm, userID string, roles []Role) error {
	const errMessage = "could not delete realm role from user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// AddRealmRoleToGroup adds realm-level role mappings
func (g *GoCloak) AddRealmRoleToGroup(ctx context.Context, token, realm, groupID string, roles []Role) error {
	const errMessage = "could not add realm role to group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// DeleteRealmRoleFromGroup deletes realm-level role mappings
func (g *GoCloak) DeleteRealmRoleFromGroup(ctx context.Context, token, realm, groupID string, roles []Role) error {
	const errMessage = "could not delete realm role from group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// AddRealmRoleComposite adds a role to the composite.
func (g *GoCloak) AddRealmRoleComposite(ctx context.Context, token, realm, roleName string, roles []Role) error {
	const errMessage = "could not add realm role composite"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "roles", roleName, "composites"))

	return checkForError(resp, err, errMessage)
}

// DeleteRealmRoleComposite deletes a role from the composite.
func (g *GoCloak) DeleteRealmRoleComposite(ctx context.Context, token, realm, roleName string, roles []Role) error {
	const errMessage = "could not delete realm role composite"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "roles", roleName, "composites"))

	return checkForError(resp, err, errMessage)
}

// GetCompositeRealmRoles returns all realm composite roles associated with the given realm role
func (g *GoCloak) GetCompositeRealmRoles(ctx context.Context, token, realm, roleName string) ([]*Role, error) {
	const errMessage = "could not get composite realm roles by role"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles", roleName, "composites"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeRolesByRoleID returns all realm composite roles associated with the given client role
func (g *GoCloak) GetCompositeRolesByRoleID(ctx context.Context, token, realm, roleID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by role id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeRealmRolesByRoleID returns all realm composite roles associated with the given client role
func (g *GoCloak) GetCompositeRealmRolesByRoleID(ctx context.Context, token, realm, roleID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by role id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites", "realm"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeRealmRolesByUserID returns all realm roles and composite roles assigned to the given user
func (g *GoCloak) GetCompositeRealmRolesByUserID(ctx context.Context, token, realm, userID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm", "composite"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCompositeRealmRolesByGroupID returns all realm roles and composite roles assigned to the given group
func (g *GoCloak) GetCompositeRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) ([]*Role, error) {
	const errMessage = "could not get composite client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm", "composite"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAvailableRealmRolesByUserID returns all available realm roles to the given user
func (g *GoCloak) GetAvailableRealmRolesByUserID(ctx context.Context, token, realm, userID string) ([]*Role, error) {
	const errMessage = "could not get available client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "realm", "available"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAvailableRealmRolesByGroupID returns all available realm roles to the given group
func (g *GoCloak) GetAvailableRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) ([]*Role, error) {
	const errMessage = "could not get available client roles by user id"

	var result []*Role
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "realm", "available"))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// -----
// Realm
// -----

// GetRealm returns top-level representation of the realm
func (g *GoCloak) GetRealm(ctx context.Context, token, realm string) (*RealmRepresentation, error) {
	const errMessage = "could not get realm"

	var result RealmRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRealms returns top-level representation of all realms
func (g *GoCloak) GetRealms(ctx context.Context, token string) ([]*RealmRepresentation, error) {
	const errMessage = "could not get realms"

	var result []*RealmRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(""))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateRealm creates a realm
func (g *GoCloak) CreateRealm(ctx context.Context, token string, realm RealmRepresentation) (string, error) {
	const errMessage = "could not create realm"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(&realm).
		Post(g.getAdminRealmURL(""))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}
	return getID(resp), nil
}

// UpdateRealm updates a given realm
func (g *GoCloak) UpdateRealm(ctx context.Context, token string, realm RealmRepresentation) error {
	const errMessage = "could not update realm"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(realm).
		Put(g.getAdminRealmURL(PString(realm.Realm)))

	return checkForError(resp, err, errMessage)
}

// DeleteRealm removes a realm
func (g *GoCloak) DeleteRealm(ctx context.Context, token, realm string) error {
	const errMessage = "could not delete realm"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm))

	return checkForError(resp, err, errMessage)
}

// ClearRealmCache clears realm cache
func (g *GoCloak) ClearRealmCache(ctx context.Context, token, realm string) error {
	const errMessage = "could not clear realm cache"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Post(g.getAdminRealmURL(realm, "clear-realm-cache"))

	return checkForError(resp, err, errMessage)
}

// ClearUserCache clears realm cache
func (g *GoCloak) ClearUserCache(ctx context.Context, token, realm string) error {
	const errMessage = "could not clear user cache"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Post(g.getAdminRealmURL(realm, "clear-user-cache"))

	return checkForError(resp, err, errMessage)
}

// ClearKeysCache clears realm cache
func (g *GoCloak) ClearKeysCache(ctx context.Context, token, realm string) error {
	const errMessage = "could not clear keys cache"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Post(g.getAdminRealmURL(realm, "clear-keys-cache"))

	return checkForError(resp, err, errMessage)
}

// GetAuthenticationFlows get all authentication flows from a realm
func (g *GoCloak) GetAuthenticationFlows(ctx context.Context, token, realm string) ([]*AuthenticationFlowRepresentation, error) {
	const errMessage = "could not retrieve authentication flows"
	var result []*AuthenticationFlowRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "authentication", "flows"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAuthenticationFlow get an authentication flow with the given ID
func (g *GoCloak) GetAuthenticationFlow(ctx context.Context, token, realm string, authenticationFlowID string) (*AuthenticationFlowRepresentation, error) {
	const errMessage = "could not retrieve authentication flows"
	var result *AuthenticationFlowRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "authentication", "flows", authenticationFlowID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateAuthenticationFlow creates a new Authentication flow in a realm
func (g *GoCloak) CreateAuthenticationFlow(ctx context.Context, token, realm string, flow AuthenticationFlowRepresentation) error {
	const errMessage = "could not create authentication flows"
	var result []*AuthenticationFlowRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).SetBody(flow).
		Post(g.getAdminRealmURL(realm, "authentication", "flows"))

	return checkForError(resp, err, errMessage)
}

// UpdateAuthenticationFlow a given Authentication Flow
func (g *GoCloak) UpdateAuthenticationFlow(ctx context.Context, token, realm string, flow AuthenticationFlowRepresentation, authenticationFlowID string) (*AuthenticationFlowRepresentation, error) {
	const errMessage = "could not create authentication flows"
	var result *AuthenticationFlowRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).SetBody(flow).
		Put(g.getAdminRealmURL(realm, "authentication", "flows", authenticationFlowID))

	if err = checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteAuthenticationFlow deletes a flow in a realm with the given ID
func (g *GoCloak) DeleteAuthenticationFlow(ctx context.Context, token, realm, flowID string) error {
	const errMessage = "could not delete authentication flows"
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "authentication", "flows", flowID))

	return checkForError(resp, err, errMessage)
}

// GetAuthenticationExecutions retrieves all executions of a given flow
func (g *GoCloak) GetAuthenticationExecutions(ctx context.Context, token, realm, flow string) ([]*ModifyAuthenticationExecutionRepresentation, error) {
	const errMessage = "could not retrieve authentication flows"
	var result []*ModifyAuthenticationExecutionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "authentication", "flows", flow, "executions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateAuthenticationExecution creates a new execution for the given flow name in the given realm
func (g *GoCloak) CreateAuthenticationExecution(ctx context.Context, token, realm, flow string, execution CreateAuthenticationExecutionRepresentation) error {
	const errMessage = "could not create authentication execution"
	resp, err := g.GetRequestWithBearerAuth(ctx, token).SetBody(execution).
		Post(g.getAdminRealmURL(realm, "authentication", "flows", flow, "executions", "execution"))

	return checkForError(resp, err, errMessage)
}

// UpdateAuthenticationExecution updates an authentication execution for the given flow in the given realm
func (g *GoCloak) UpdateAuthenticationExecution(ctx context.Context, token, realm, flow string, execution ModifyAuthenticationExecutionRepresentation) error {
	const errMessage = "could not update authentication execution"
	resp, err := g.GetRequestWithBearerAuth(ctx, token).SetBody(execution).
		Put(g.getAdminRealmURL(realm, "authentication", "flows", flow, "executions"))

	return checkForError(resp, err, errMessage)
}

// DeleteAuthenticationExecution delete a single execution with the given ID
func (g *GoCloak) DeleteAuthenticationExecution(ctx context.Context, token, realm, executionID string) error {
	const errMessage = "could not delete authentication execution"
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "authentication", "executions", executionID))

	return checkForError(resp, err, errMessage)
}

// CreateAuthenticationExecutionFlow creates a new execution for the given flow name in the given realm
func (g *GoCloak) CreateAuthenticationExecutionFlow(ctx context.Context, token, realm, flow string, executionFlow CreateAuthenticationExecutionFlowRepresentation) error {
	const errMessage = "could not create authentication execution flow"
	resp, err := g.GetRequestWithBearerAuth(ctx, token).SetBody(executionFlow).
		Post(g.getAdminRealmURL(realm, "authentication", "flows", flow, "executions", "flow"))

	return checkForError(resp, err, errMessage)
}

// -----
// Users
// -----

// CreateUser creates the given user in the given realm and returns it's userID
// Note: Keycloak has not documented what members of the User object are actually being accepted, when creating a user.
// Things like RealmRoles must be attached using followup calls to the respective functions.
func (g *GoCloak) CreateUser(ctx context.Context, token, realm string, user User) (string, error) {
	const errMessage = "could not create user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(user).
		Post(g.getAdminRealmURL(realm, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// DeleteUser delete a given user
func (g *GoCloak) DeleteUser(ctx context.Context, token, realm, userID string) error {
	const errMessage = "could not delete user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID))

	return checkForError(resp, err, errMessage)
}

// GetUserByID fetches a user from the given realm with the given userID
func (g *GoCloak) GetUserByID(ctx context.Context, accessToken, realm, userID string) (*User, error) {
	const errMessage = "could not get user by id"

	if userID == "" {
		return nil, errors.Wrap(errors.New("userID shall not be empty"), errMessage)
	}

	var result User
	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUserCount gets the user count in the realm
func (g *GoCloak) GetUserCount(ctx context.Context, token string, realm string, params GetUsersParams) (int, error) {
	const errMessage = "could not get user count"

	var result int
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return 0, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "users", "count"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return -1, errors.Wrap(err, errMessage)
	}

	return result, nil
}

// GetUserGroups get all groups for user
func (g *GoCloak) GetUserGroups(ctx context.Context, token, realm, userID string, params GetGroupsParams) ([]*Group, error) {
	const errMessage = "could not get user groups"

	var result []*Group
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "users", userID, "groups"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetUsers get all users in realm
func (g *GoCloak) GetUsers(ctx context.Context, token, realm string, params GetUsersParams) ([]*User, error) {
	const errMessage = "could not get users"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetUsersByRoleName returns all users have a given role
func (g *GoCloak) GetUsersByRoleName(ctx context.Context, token, realm, roleName string, params GetUsersByRoleParams) ([]*User, error) {
	const errMessage = "could not get users by role name"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "roles", roleName, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetUsersByClientRoleName returns all users have a given client role
func (g *GoCloak) GetUsersByClientRoleName(ctx context.Context, token, realm, idOfClient, roleName string, params GetUsersByRoleParams) ([]*User, error) {
	const errMessage = "could not get users by client role name"

	var result []*User
	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, err
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "roles", roleName, "users"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// SetPassword sets a new password for the user with the given id. Needs elevated privileges
func (g *GoCloak) SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error {
	const errMessage = "could not set password"

	requestBody := SetPasswordRequest{Password: &password, Temporary: &temporary, Type: StringP("password")}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(requestBody).
		Put(g.getAdminRealmURL(realm, "users", userID, "reset-password"))

	return checkForError(resp, err, errMessage)
}

// UpdateUser updates a given user
func (g *GoCloak) UpdateUser(ctx context.Context, token, realm string, user User) error {
	const errMessage = "could not update user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(user).
		Put(g.getAdminRealmURL(realm, "users", PString(user.ID)))

	return checkForError(resp, err, errMessage)
}

// AddUserToGroup puts given user to given group
func (g *GoCloak) AddUserToGroup(ctx context.Context, token, realm, userID, groupID string) error {
	const errMessage = "could not add user to group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Put(g.getAdminRealmURL(realm, "users", userID, "groups", groupID))

	return checkForError(resp, err, errMessage)
}

// DeleteUserFromGroup deletes given user from given group
func (g *GoCloak) DeleteUserFromGroup(ctx context.Context, token, realm, userID, groupID string) error {
	const errMessage = "could not delete user from group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID, "groups", groupID))

	return checkForError(resp, err, errMessage)
}

// GetUserSessions returns user sessions associated with the user
func (g *GoCloak) GetUserSessions(ctx context.Context, token, realm, userID string) ([]*UserSessionRepresentation, error) {
	const errMessage = "could not get user sessions"

	var res []*UserSessionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		Get(g.getAdminRealmURL(realm, "users", userID, "sessions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserOfflineSessionsForClient returns offline sessions associated with the user and client
func (g *GoCloak) GetUserOfflineSessionsForClient(ctx context.Context, token, realm, userID, idOfClient string) ([]*UserSessionRepresentation, error) {
	const errMessage = "could not get user offline sessions for client"

	var res []*UserSessionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		Get(g.getAdminRealmURL(realm, "users", userID, "offline-sessions", idOfClient))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return res, nil
}

// AddClientRolesToUser adds client-level role mappings
func (g *GoCloak) AddClientRolesToUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) error {
	const errMessage = "could not add client role to user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// AddClientRoleToUser adds client-level role mappings
//
// Deprecated: replaced by AddClientRolesToUser
func (g *GoCloak) AddClientRoleToUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) error {
	return g.AddClientRolesToUser(ctx, token, realm, idOfClient, userID, roles)
}

// AddClientRolesToGroup adds a client role to the group
func (g *GoCloak) AddClientRolesToGroup(ctx context.Context, token, realm, idOfClient, groupID string, roles []Role) error {
	const errMessage = "could not add client role to group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// AddClientRoleToGroup adds a client role to the group
//
// Deprecated: replaced by AddClientRolesToGroup
func (g *GoCloak) AddClientRoleToGroup(ctx context.Context, token, realm, idOfClient, groupID string, roles []Role) error {
	return g.AddClientRolesToGroup(ctx, token, realm, idOfClient, groupID, roles)
}

// DeleteClientRolesFromUser adds client-level role mappings
func (g *GoCloak) DeleteClientRolesFromUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) error {
	const errMessage = "could not delete client role from user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "users", userID, "role-mappings", "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// DeleteClientRoleFromUser adds client-level role mappings
//
// Deprecated: replaced by DeleteClientRolesFrom
func (g *GoCloak) DeleteClientRoleFromUser(ctx context.Context, token, realm, idOfClient, userID string, roles []Role) error {
	return g.DeleteClientRolesFromUser(ctx, token, realm, idOfClient, userID, roles)
}

// DeleteClientRoleFromGroup removes a client role from from the group
func (g *GoCloak) DeleteClientRoleFromGroup(ctx context.Context, token, realm, idOfClient, groupID string, roles []Role) error {
	const errMessage = "could not client role from group"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "groups", groupID, "role-mappings", "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// AddClientRoleComposite adds roles as composite
func (g *GoCloak) AddClientRoleComposite(ctx context.Context, token, realm, roleID string, roles []Role) error {
	const errMessage = "could not add client role composite"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites"))

	return checkForError(resp, err, errMessage)
}

// DeleteClientRoleComposite deletes composites from a role
func (g *GoCloak) DeleteClientRoleComposite(ctx context.Context, token, realm, roleID string, roles []Role) error {
	const errMessage = "could not delete client role composite"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "roles-by-id", roleID, "composites"))

	return checkForError(resp, err, errMessage)
}

// GetUserFederatedIdentities gets all user federated identities
func (g *GoCloak) GetUserFederatedIdentities(ctx context.Context, token, realm, userID string) ([]*FederatedIdentityRepresentation, error) {
	const errMessage = "could not get user federated identities"

	var res []*FederatedIdentityRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&res).
		Get(g.getAdminRealmURL(realm, "users", userID, "federated-identity"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return res, err
}

// CreateUserFederatedIdentity creates an user federated identity
func (g *GoCloak) CreateUserFederatedIdentity(ctx context.Context, token, realm, userID, providerID string, federatedIdentityRep FederatedIdentityRepresentation) error {
	const errMessage = "could not create user federeated identity"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(federatedIdentityRep).
		Post(g.getAdminRealmURL(realm, "users", userID, "federated-identity", providerID))

	return checkForError(resp, err, errMessage)
}

// DeleteUserFederatedIdentity deletes an user federated identity
func (g *GoCloak) DeleteUserFederatedIdentity(ctx context.Context, token, realm, userID, providerID string) error {
	const errMessage = "could not delete user federeated identity"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID, "federated-identity", providerID))

	return checkForError(resp, err, errMessage)
}

// GetUserBruteForceDetectionStatus fetches a user status regarding brute force protection
func (g *GoCloak) GetUserBruteForceDetectionStatus(ctx context.Context, accessToken, realm, userID string) (*BruteForceStatus, error) {
	const errMessage = "could not brute force detection Status"
	var result BruteForceStatus

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getAttackDetectionURL(realm, "users", userID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// ------------------
// Identity Providers
// ------------------

// CreateIdentityProvider creates an identity provider in a realm
func (g *GoCloak) CreateIdentityProvider(ctx context.Context, token string, realm string, providerRep IdentityProviderRepresentation) (string, error) {
	const errMessage = "could not create identity provider"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(providerRep).
		Post(g.getAdminRealmURL(realm, "identity-provider", "instances"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// GetIdentityProviders returns list of identity providers in a realm
func (g *GoCloak) GetIdentityProviders(ctx context.Context, token, realm string) ([]*IdentityProviderRepresentation, error) {
	const errMessage = "could not get identity providers"

	var result []*IdentityProviderRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "identity-provider", "instances"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetIdentityProvider gets the identity provider in a realm
func (g *GoCloak) GetIdentityProvider(ctx context.Context, token, realm, alias string) (*IdentityProviderRepresentation, error) {
	const errMessage = "could not get identity provider"

	var result IdentityProviderRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "identity-provider", "instances", alias))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateIdentityProvider updates the identity provider in a realm
func (g *GoCloak) UpdateIdentityProvider(ctx context.Context, token, realm, alias string, providerRep IdentityProviderRepresentation) error {
	const errMessage = "could not update identity provider"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(providerRep).
		Put(g.getAdminRealmURL(realm, "identity-provider", "instances", alias))

	return checkForError(resp, err, errMessage)
}

// DeleteIdentityProvider deletes the identity provider in a realm
func (g *GoCloak) DeleteIdentityProvider(ctx context.Context, token, realm, alias string) error {
	const errMessage = "could not delete identity provider"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "identity-provider", "instances", alias))

	return checkForError(resp, err, errMessage)
}

// ExportIDPPublicBrokerConfig exports the broker config for a given alias
func (g *GoCloak) ExportIDPPublicBrokerConfig(ctx context.Context, token, realm, alias string) (*string, error) {
	const errMessage = "could not get public identity provider configuration"

	resp, err := g.GetRequestWithBearerAuthXMLHeader(ctx, token).
		Get(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "export"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	result := resp.String()
	return &result, nil
}

// ImportIdentityProviderConfig parses and returns the identity provider config at a given URL
func (g *GoCloak) ImportIdentityProviderConfig(ctx context.Context, token, realm, fromURL, providerID string) (map[string]string, error) {
	const errMessage = "could not import config"

	result := make(map[string]string)
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(map[string]string{
			"fromUrl":    fromURL,
			"providerId": providerID,
		}).
		Post(g.getAdminRealmURL(realm, "identity-provider", "import-config"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// ImportIdentityProviderConfigFromFile parses and returns the identity provider config from a given file
func (g *GoCloak) ImportIdentityProviderConfigFromFile(ctx context.Context, token, realm, providerID, fileName string, fileBody io.Reader) (map[string]string, error) {
	const errMessage = "could not import config"

	result := make(map[string]string)
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetFileReader("file", fileName, fileBody).
		SetFormData(map[string]string{
			"providerId": providerID,
		}).
		Post(g.getAdminRealmURL(realm, "identity-provider", "import-config"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateIdentityProviderMapper creates an instance of an identity provider mapper associated with the given alias
func (g *GoCloak) CreateIdentityProviderMapper(ctx context.Context, token, realm, alias string, mapper IdentityProviderMapper) (string, error) {
	const errMessage = "could not create mapper for identity provider"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(mapper).
		Post(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "mappers"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	return getID(resp), nil
}

// GetIdentityProviderMapper gets the mapper by id for the given identity provider alias in a realm
func (g *GoCloak) GetIdentityProviderMapper(ctx context.Context, token string, realm string, alias string, mapperID string) (*IdentityProviderMapper, error) {
	const errMessage = "could not get identity provider mapper"

	result := IdentityProviderMapper{}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "mappers", mapperID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteIdentityProviderMapper deletes an instance of an identity provider mapper associated with the given alias and mapper ID
func (g *GoCloak) DeleteIdentityProviderMapper(ctx context.Context, token, realm, alias, mapperID string) error {
	const errMessage = "could not delete mapper for identity provider"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "mappers", mapperID))

	return checkForError(resp, err, errMessage)
}

// GetIdentityProviderMappers returns list of mappers associated with an identity provider
func (g *GoCloak) GetIdentityProviderMappers(ctx context.Context, token, realm, alias string) ([]*IdentityProviderMapper, error) {
	const errMessage = "could not get identity provider mappers"

	var result []*IdentityProviderMapper
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "mappers"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetIdentityProviderMapperByID gets the mapper of an identity provider
func (g *GoCloak) GetIdentityProviderMapperByID(ctx context.Context, token, realm, alias, mapperID string) (*IdentityProviderMapper, error) {
	const errMessage = "could not get identity provider mappers"

	var result IdentityProviderMapper
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "mappers", mapperID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateIdentityProviderMapper updates mapper of an identity provider
func (g *GoCloak) UpdateIdentityProviderMapper(ctx context.Context, token, realm, alias string, mapper IdentityProviderMapper) error {
	const errMessage = "could not update identity provider mapper"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(mapper).
		Put(g.getAdminRealmURL(realm, "identity-provider", "instances", alias, "mappers", PString(mapper.ID)))

	return checkForError(resp, err, errMessage)
}

// ------------------
// Protection API
// ------------------

// GetResource returns a client's resource with the given id, using access token from admin
func (g *GoCloak) GetResource(ctx context.Context, token, realm, idOfClient, resourceID string) (*ResourceRepresentation, error) {
	const errMessage = "could not get resource"

	var result ResourceRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "resource", resourceID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetResourceClient returns a client's resource with the given id, using access token from client
func (g *GoCloak) GetResourceClient(ctx context.Context, token, realm, resourceID string) (*ResourceRepresentation, error) {
	const errMessage = "could not get resource"

	var result ResourceRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getRealmURL(realm, "authz", "protection", "resource_set", resourceID))

	// http://${host}:${port}/auth/realms/${realm_name}/authz/protection/resource_set/{resource_id}

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetResources returns resources associated with the client, using access token from admin
func (g *GoCloak) GetResources(ctx context.Context, token, realm, idOfClient string, params GetResourceParams) ([]*ResourceRepresentation, error) {
	const errMessage = "could not get resources"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, err
	}

	var result []*ResourceRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "resource"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetResourcesClient returns resources associated with the client, using access token from client
func (g *GoCloak) GetResourcesClient(ctx context.Context, token, realm string, params GetResourceParams) ([]*ResourceRepresentation, error) {
	const errMessage = "could not get resources"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, err
	}

	var result []*ResourceRepresentation
	var resourceIDs []string
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&resourceIDs).
		SetQueryParams(queryParams).
		Get(g.getRealmURL(realm, "authz", "protection", "resource_set"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	for _, resourceID := range resourceIDs {
		resource, err := g.GetResourceClient(ctx, token, realm, resourceID)
		if err == nil {
			result = append(result, resource)
		}
	}

	return result, nil
}

// GetResourceServer returns resource server settings.
// The access token must have the realm view_clients role on its service
// account to be allowed to call this endpoint.
func (g *GoCloak) GetResourceServer(ctx context.Context, token, realm, idOfClient string) (*ResourceServerRepresentation, error) {
	const errMessage = "could not get resource server settings"

	var result *ResourceServerRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "settings"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateResource updates a resource associated with the client, using access token from admin
func (g *GoCloak) UpdateResource(ctx context.Context, token, realm, idOfClient string, resource ResourceRepresentation) error {
	const errMessage = "could not update resource"

	if NilOrEmpty(resource.ID) {
		return errors.New("ID of a resource required")
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(resource).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "resource", *(resource.ID)))

	return checkForError(resp, err, errMessage)
}

// UpdateResourceClient updates a resource associated with the client, using access token from client
func (g *GoCloak) UpdateResourceClient(ctx context.Context, token, realm string, resource ResourceRepresentation) error {
	const errMessage = "could not update resource"

	if NilOrEmpty(resource.ID) {
		return errors.New("ID of a resource required")
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(resource).
		Put(g.getRealmURL(realm, "authz", "protection", "resource_set", *(resource.ID)))

	return checkForError(resp, err, errMessage)
}

// CreateResource creates a resource associated with the client, using access token from admin
func (g *GoCloak) CreateResource(ctx context.Context, token, realm string, idOfClient string, resource ResourceRepresentation) (*ResourceRepresentation, error) {
	const errMessage = "could not create resource"

	var result ResourceRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(resource).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "resource"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateResourceClient creates a resource associated with the client, using access token from client
func (g *GoCloak) CreateResourceClient(ctx context.Context, token, realm string, resource ResourceRepresentation) (*ResourceRepresentation, error) {
	const errMessage = "could not create resource"

	var result ResourceRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(resource).
		Post(g.getRealmURL(realm, "authz", "protection", "resource_set"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteResource deletes a resource associated with the client (using an admin token)
func (g *GoCloak) DeleteResource(ctx context.Context, token, realm, idOfClient, resourceID string) error {
	const errMessage = "could not delete resource"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "resource", resourceID))

	return checkForError(resp, err, errMessage)
}

// DeleteResourceClient deletes a resource associated with the client (using a client token)
func (g *GoCloak) DeleteResourceClient(ctx context.Context, token, realm, resourceID string) error {
	const errMessage = "could not delete resource"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getRealmURL(realm, "authz", "protection", "resource_set", resourceID))

	return checkForError(resp, err, errMessage)
}

// GetScope returns a client's scope with the given id
func (g *GoCloak) GetScope(ctx context.Context, token, realm, idOfClient, scopeID string) (*ScopeRepresentation, error) {
	const errMessage = "could not get scope"

	var result ScopeRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "scope", scopeID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetScopes returns scopes associated with the client
func (g *GoCloak) GetScopes(ctx context.Context, token, realm, idOfClient string, params GetScopeParams) ([]*ScopeRepresentation, error) {
	const errMessage = "could not get scopes"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, err
	}
	var result []*ScopeRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "scope"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateScope creates a scope associated with the client
func (g *GoCloak) CreateScope(ctx context.Context, token, realm, idOfClient string, scope ScopeRepresentation) (*ScopeRepresentation, error) {
	const errMessage = "could not create scope"

	var result ScopeRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(scope).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "scope"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPermissionScope gets the permission scope associated with the client
func (g *GoCloak) GetPermissionScope(ctx context.Context, token, realm, idOfClient string, idOfScope string) (*PolicyRepresentation, error) {
	const errMessage = "could not get permission scope"

	var result PolicyRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", "scope", idOfScope))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdatePermissionScope updates a permission scope associated with the client
func (g *GoCloak) UpdatePermissionScope(ctx context.Context, token, realm, idOfClient string, idOfScope string, policy PolicyRepresentation) error {
	const errMessage = "could not create permission scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(policy).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", "scope", idOfScope))

	return checkForError(resp, err, errMessage)
}

// UpdateScope updates a scope associated with the client
func (g *GoCloak) UpdateScope(ctx context.Context, token, realm, idOfClient string, scope ScopeRepresentation) error {
	const errMessage = "could not update scope"

	if NilOrEmpty(scope.ID) {
		return errors.New("ID of a scope required")
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(scope).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "scope", *(scope.ID)))

	return checkForError(resp, err, errMessage)
}

// DeleteScope deletes a scope associated with the client
func (g *GoCloak) DeleteScope(ctx context.Context, token, realm, idOfClient, scopeID string) error {
	const errMessage = "could not delete scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "scope", scopeID))

	return checkForError(resp, err, errMessage)
}

// GetPolicy returns a client's policy with the given id
func (g *GoCloak) GetPolicy(ctx context.Context, token, realm, idOfClient, policyID string) (*PolicyRepresentation, error) {
	const errMessage = "could not get policy"

	var result PolicyRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", policyID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPolicies returns policies associated with the client
func (g *GoCloak) GetPolicies(ctx context.Context, token, realm, idOfClient string, params GetPolicyParams) ([]*PolicyRepresentation, error) {
	const errMessage = "could not get policies"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	path := []string{"clients", idOfClient, "authz", "resource-server", "policy"}
	if !NilOrEmpty(params.Type) {
		path = append(path, *params.Type)
	}

	var result []*PolicyRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, path...))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreatePolicy creates a policy associated with the client
func (g *GoCloak) CreatePolicy(ctx context.Context, token, realm, idOfClient string, policy PolicyRepresentation) (*PolicyRepresentation, error) {
	const errMessage = "could not create policy"

	if NilOrEmpty(policy.Type) {
		return nil, errors.New("type of a policy required")
	}

	var result PolicyRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(policy).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", *(policy.Type)))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdatePolicy updates a policy associated with the client
func (g *GoCloak) UpdatePolicy(ctx context.Context, token, realm, idOfClient string, policy PolicyRepresentation) error {
	const errMessage = "could not update policy"

	if NilOrEmpty(policy.ID) {
		return errors.New("ID of a policy required")
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(policy).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", *(policy.Type), *(policy.ID)))

	return checkForError(resp, err, errMessage)
}

// DeletePolicy deletes a policy associated with the client
func (g *GoCloak) DeletePolicy(ctx context.Context, token, realm, idOfClient, policyID string) error {
	const errMessage = "could not delete policy"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", policyID))

	return checkForError(resp, err, errMessage)
}

// GetAuthorizationPolicyAssociatedPolicies returns a client's associated policies of specific policy with the given policy id, using access token from admin
func (g *GoCloak) GetAuthorizationPolicyAssociatedPolicies(ctx context.Context, token, realm, idOfClient, policyID string) ([]*PolicyRepresentation, error) {
	const errMessage = "could not get policy associated policies"

	var result []*PolicyRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", policyID, "associatedPolicies"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAuthorizationPolicyResources returns a client's resources of specific policy with the given policy id, using access token from admin
func (g *GoCloak) GetAuthorizationPolicyResources(ctx context.Context, token, realm, idOfClient, policyID string) ([]*PolicyResourceRepresentation, error) {
	const errMessage = "could not get policy resources"

	var result []*PolicyResourceRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", policyID, "resources"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAuthorizationPolicyScopes returns a client's scopes of specific policy with the given policy id, using access token from admin
func (g *GoCloak) GetAuthorizationPolicyScopes(ctx context.Context, token, realm, idOfClient, policyID string) ([]*PolicyScopeRepresentation, error) {
	const errMessage = "could not get policy scopes"

	var result []*PolicyScopeRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", policyID, "scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetResourcePolicy updates a permission for a specific resource, using token obtained by Resource Owner Password Credentials Grant or Token exchange
func (g *GoCloak) GetResourcePolicy(ctx context.Context, token, realm, permissionID string) (*ResourcePolicyRepresentation, error) {
	const errMessage = "could not get resource policy"

	var result ResourcePolicyRepresentation
	resp, err := g.GetRequestWithBearerAuthNoCache(ctx, token).
		SetResult(&result).
		Get(g.getRealmURL(realm, "authz", "protection", "uma-policy", permissionID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetResourcePolicies returns resources associated with the client, using token obtained by Resource Owner Password Credentials Grant or Token exchange
func (g *GoCloak) GetResourcePolicies(ctx context.Context, token, realm string, params GetResourcePoliciesParams) ([]*ResourcePolicyRepresentation, error) {
	const errMessage = "could not get resource policies"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, err
	}

	var result []*ResourcePolicyRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getRealmURL(realm, "authz", "protection", "uma-policy"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateResourcePolicy associates a permission with a specific resource, using token obtained by Resource Owner Password Credentials Grant or Token exchange
func (g *GoCloak) CreateResourcePolicy(ctx context.Context, token, realm, resourceID string, policy ResourcePolicyRepresentation) (*ResourcePolicyRepresentation, error) {
	const errMessage = "could not create resource policy"

	var result ResourcePolicyRepresentation
	resp, err := g.GetRequestWithBearerAuthNoCache(ctx, token).
		SetResult(&result).
		SetBody(policy).
		Post(g.getRealmURL(realm, "authz", "protection", "uma-policy", resourceID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateResourcePolicy updates a permission for a specific resource, using token obtained by Resource Owner Password Credentials Grant or Token exchange
func (g *GoCloak) UpdateResourcePolicy(ctx context.Context, token, realm, permissionID string, policy ResourcePolicyRepresentation) error {
	const errMessage = "could not update resource policy"

	resp, err := g.GetRequestWithBearerAuthNoCache(ctx, token).
		SetBody(policy).
		Put(g.getRealmURL(realm, "authz", "protection", "uma-policy", permissionID))

	return checkForError(resp, err, errMessage)
}

// DeleteResourcePolicy deletes a permission for a specific resource, using token obtained by Resource Owner Password Credentials Grant or Token exchange
func (g *GoCloak) DeleteResourcePolicy(ctx context.Context, token, realm, permissionID string) error {
	const errMessage = "could not  delete resource policy"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getRealmURL(realm, "authz", "protection", "uma-policy", permissionID))

	return checkForError(resp, err, errMessage)
}

// GetPermission returns a client's permission with the given id
func (g *GoCloak) GetPermission(ctx context.Context, token, realm, idOfClient, permissionID string) (*PermissionRepresentation, error) {
	const errMessage = "could not get permission"

	var result PermissionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", permissionID))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDependentPermissions returns a client's permission with the given policy id
func (g *GoCloak) GetDependentPermissions(ctx context.Context, token, realm, idOfClient, policyID string) ([]*PermissionRepresentation, error) {
	const errMessage = "could not get permission"

	var result []*PermissionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "policy", policyID, "dependentPolicies"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPermissionResources returns a client's resource attached for the given permission id
func (g *GoCloak) GetPermissionResources(ctx context.Context, token, realm, idOfClient, permissionID string) ([]*PermissionResource, error) {
	const errMessage = "could not get permission resource"

	var result []*PermissionResource
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", permissionID, "resources"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPermissionScopes returns a client's scopes configured for the given permission id
func (g *GoCloak) GetPermissionScopes(ctx context.Context, token, realm, idOfClient, permissionID string) ([]*PermissionScope, error) {
	const errMessage = "could not get permission scopes"

	var result []*PermissionScope
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", permissionID, "scopes"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPermissions returns permissions associated with the client
func (g *GoCloak) GetPermissions(ctx context.Context, token, realm, idOfClient string, params GetPermissionParams) ([]*PermissionRepresentation, error) {
	const errMessage = "could not get permissions"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	path := []string{"clients", idOfClient, "authz", "resource-server", "permission"}
	if !NilOrEmpty(params.Type) {
		path = append(path, *params.Type)
	}

	var result []*PermissionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, path...))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// checkPermissionTicketParams checks that mandatory fields are present
func checkPermissionTicketParams(permissions []CreatePermissionTicketParams) error {
	if len(permissions) == 0 {
		return errors.New("at least one permission ticket must be requested")
	}

	for _, pt := range permissions {

		if NilOrEmpty(pt.ResourceID) {
			return errors.New("resourceID required for permission ticket")
		}
		if NilOrEmptyArray(pt.ResourceScopes) {
			return errors.New("at least one resourceScope required for permission ticket")
		}
	}

	return nil
}

// CreatePermissionTicket creates a permission ticket, using access token from client
func (g *GoCloak) CreatePermissionTicket(ctx context.Context, token, realm string, permissions []CreatePermissionTicketParams) (*PermissionTicketResponseRepresentation, error) {
	const errMessage = "could not create permission ticket"

	err := checkPermissionTicketParams(permissions)
	if err != nil {
		return nil, err
	}

	var result PermissionTicketResponseRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(permissions).
		Post(g.getRealmURL(realm, "authz", "protection", "permission"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// checkPermissionGrantParams checks for mandatory fields
func checkPermissionGrantParams(permission PermissionGrantParams) error {
	if NilOrEmpty(permission.RequesterID) {
		return errors.New("requesterID required to grant user permission")
	}
	if NilOrEmpty(permission.ResourceID) {
		return errors.New("resourceID required to grant user permission")
	}
	if NilOrEmpty(permission.ScopeName) {
		return errors.New("scopeName required to grant user permission")
	}

	return nil
}

// GrantUserPermission lets resource owner grant permission for specific resource ID to specific user ID
func (g *GoCloak) GrantUserPermission(ctx context.Context, token, realm string, permission PermissionGrantParams) (*PermissionGrantResponseRepresentation, error) {
	const errMessage = "could not grant user permission"

	err := checkPermissionGrantParams(permission)
	if err != nil {
		return nil, err
	}

	permission.Granted = BoolP(true)

	var result PermissionGrantResponseRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(permission).
		Post(g.getRealmURL(realm, "authz", "protection", "permission", "ticket"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// checkPermissionUpdateParams
func checkPermissionUpdateParams(permission PermissionGrantParams) error {
	err := checkPermissionGrantParams(permission)
	if err != nil {
		return err
	}

	if permission.Granted == nil {
		return errors.New("granted required to update user permission")
	}
	return nil
}

// UpdateUserPermission updates user permissions.
func (g *GoCloak) UpdateUserPermission(ctx context.Context, token, realm string, permission PermissionGrantParams) (*PermissionGrantResponseRepresentation, error) {
	const errMessage = "could not update user permission"

	err := checkPermissionUpdateParams(permission)
	if err != nil {
		return nil, err
	}

	var result PermissionGrantResponseRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(permission).
		Put(g.getRealmURL(realm, "authz", "protection", "permission", "ticket"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNoContent { // permission updated to 'not granted' removes permission
		return nil, nil
	}

	return &result, nil
}

// GetUserPermissions gets granted permissions according query parameters
func (g *GoCloak) GetUserPermissions(ctx context.Context, token, realm string, params GetUserPermissionParams) ([]*PermissionGrantResponseRepresentation, error) {
	const errMessage = "could not get user permissions"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, err
	}

	var result []*PermissionGrantResponseRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getRealmURL(realm, "authz", "protection", "permission", "ticket"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteUserPermission revokes permissions according query parameters
func (g *GoCloak) DeleteUserPermission(ctx context.Context, token, realm, ticketID string) error {
	const errMessage = "could not delete user permission"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getRealmURL(realm, "authz", "protection", "permission", "ticket", ticketID))

	return checkForError(resp, err, errMessage)
}

// CreatePermission creates a permission associated with the client
func (g *GoCloak) CreatePermission(ctx context.Context, token, realm, idOfClient string, permission PermissionRepresentation) (*PermissionRepresentation, error) {
	const errMessage = "could not create permission"

	if NilOrEmpty(permission.Type) {
		return nil, errors.New("type of a permission required")
	}

	var result PermissionRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetBody(permission).
		Post(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", *(permission.Type)))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdatePermission updates a permission associated with the client
func (g *GoCloak) UpdatePermission(ctx context.Context, token, realm, idOfClient string, permission PermissionRepresentation) error {
	const errMessage = "could not update permission"

	if NilOrEmpty(permission.ID) {
		return errors.New("ID of a permission required")
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(permission).
		Put(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", *permission.Type, *permission.ID))

	return checkForError(resp, err, errMessage)
}

// DeletePermission deletes a policy associated with the client
func (g *GoCloak) DeletePermission(ctx context.Context, token, realm, idOfClient, permissionID string) error {
	const errMessage = "could not delete permission"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "clients", idOfClient, "authz", "resource-server", "permission", permissionID))

	return checkForError(resp, err, errMessage)
}

// ---------------
// Credentials API
// ---------------

// GetCredentialRegistrators returns credentials registrators
func (g *GoCloak) GetCredentialRegistrators(ctx context.Context, token, realm string) ([]string, error) {
	const errMessage = "could not get user credential registrators"

	var result []string
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "credential-registrators"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetConfiguredUserStorageCredentialTypes returns credential types, which are provided by the user storage where user is stored
func (g *GoCloak) GetConfiguredUserStorageCredentialTypes(ctx context.Context, token, realm, userID string) ([]string, error) {
	const errMessage = "could not get user credential registrators"

	var result []string
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "configured-user-storage-credential-types"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCredentials returns credentials available for a given user
func (g *GoCloak) GetCredentials(ctx context.Context, token, realm, userID string) ([]*CredentialRepresentation, error) {
	const errMessage = "could not get user credentials"

	var result []*CredentialRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users", userID, "credentials"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteCredentials deletes the given credential for a given user
func (g *GoCloak) DeleteCredentials(ctx context.Context, token, realm, userID, credentialID string) error {
	const errMessage = "could not delete user credentials"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "users", userID, "credentials", credentialID))

	return checkForError(resp, err, errMessage)
}

// UpdateCredentialUserLabel updates label for the given credential for the given user
func (g *GoCloak) UpdateCredentialUserLabel(ctx context.Context, token, realm, userID, credentialID, userLabel string) error {
	const errMessage = "could not update credential label for a user"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetHeader("Content-Type", "text/plain").
		SetBody(userLabel).
		Put(g.getAdminRealmURL(realm, "users", userID, "credentials", credentialID, "userLabel"))

	return checkForError(resp, err, errMessage)
}

// DisableAllCredentialsByType disables all credentials for a user of a specific type
func (g *GoCloak) DisableAllCredentialsByType(ctx context.Context, token, realm, userID string, types []string) error {
	const errMessage = "could not update disable credentials"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(types).
		Put(g.getAdminRealmURL(realm, "users", userID, "disable-credential-types"))

	return checkForError(resp, err, errMessage)
}

// MoveCredentialBehind move a credential to a position behind another credential
func (g *GoCloak) MoveCredentialBehind(ctx context.Context, token, realm, userID, credentialID, newPreviousCredentialID string) error {
	const errMessage = "could not move credential"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Post(g.getAdminRealmURL(realm, "users", userID, "credentials", credentialID, "moveAfter", newPreviousCredentialID))

	return checkForError(resp, err, errMessage)
}

// MoveCredentialToFirst move a credential to a first position in the credentials list of the user
func (g *GoCloak) MoveCredentialToFirst(ctx context.Context, token, realm, userID, credentialID string) error {
	const errMessage = "could not move credential"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Post(g.getAdminRealmURL(realm, "users", userID, "credentials", credentialID, "moveToFirst"))

	return checkForError(resp, err, errMessage)
}

// GetEvents returns events
func (g *GoCloak) GetEvents(ctx context.Context, token string, realm string, params GetEventsParams) ([]*EventRepresentation, error) {
	const errMessage = "could not get events"

	queryParams, err := GetQueryParams(params)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	var result []*EventRepresentation
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		SetQueryParams(queryParams).
		Get(g.getAdminRealmURL(realm, "events"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopesScopeMappingsRealmRolesAvailable returns realm-level roles that are available to attach to this client scope
func (g *GoCloak) GetClientScopesScopeMappingsRealmRolesAvailable(ctx context.Context, token, realm, clientScopeID string) ([]*Role, error) {
	const errMessage = "could not get available realm-level roles with the client-scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", clientScopeID, "scope-mappings", "realm", "available"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopesScopeMappingsRealmRoles returns roles associated with a client-scope
func (g *GoCloak) GetClientScopesScopeMappingsRealmRoles(ctx context.Context, token, realm, clientScopeID string) ([]*Role, error) {
	const errMessage = "could not get realm-level roles with the client-scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", clientScopeID, "scope-mappings", "realm"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteClientScopesScopeMappingsRealmRoles deletes realm-level roles from the client-scope
func (g *GoCloak) DeleteClientScopesScopeMappingsRealmRoles(ctx context.Context, token, realm, clientScopeID string, roles []Role) error {
	const errMessage = "could not delete realm-level roles from the client-scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "client-scopes", clientScopeID, "scope-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// CreateClientScopesScopeMappingsRealmRoles creates realm-level roles to the client scope
func (g *GoCloak) CreateClientScopesScopeMappingsRealmRoles(ctx context.Context, token, realm, clientScopeID string, roles []Role) error {
	const errMessage = "could not create realm-level roles to the client-scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "client-scopes", clientScopeID, "scope-mappings", "realm"))

	return checkForError(resp, err, errMessage)
}

// RegisterRequiredAction creates a required action for a given realm
func (g *GoCloak) RegisterRequiredAction(ctx context.Context, token string, realm string, requiredAction RequiredActionProviderRepresentation) error {
	const errMessage = "could not create required action"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(requiredAction).
		Post(g.getAdminRealmURL(realm, "authentication", "register-required-action"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return err
	}

	return err
}

// GetRequiredActions gets a list of required actions for a given realm
func (g *GoCloak) GetRequiredActions(ctx context.Context, token string, realm string) ([]*RequiredActionProviderRepresentation, error) {
	const errMessage = "could not get required actions"
	var result []*RequiredActionProviderRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "authentication", "required-actions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, err
}

// GetRequiredAction gets a required action for a given realm
func (g *GoCloak) GetRequiredAction(ctx context.Context, token string, realm string, alias string) (*RequiredActionProviderRepresentation, error) {
	const errMessage = "could not get required action"
	var result RequiredActionProviderRepresentation

	if alias == "" {
		return nil, errors.New("alias is required for getting a required action")
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "authentication", "required-actions", alias))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, err
}

// UpdateRequiredAction updates a required action for a given realm
func (g *GoCloak) UpdateRequiredAction(ctx context.Context, token string, realm string, requiredAction RequiredActionProviderRepresentation) error {
	const errMessage = "could not update required action"

	if NilOrEmpty(requiredAction.ProviderID) {
		return errors.New("providerId is required for updating a required action")
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(requiredAction).
		Put(g.getAdminRealmURL(realm, "authentication", "required-actions", *requiredAction.ProviderID))

	return checkForError(resp, err, errMessage)
}

// DeleteRequiredAction updates a required action for a given realm
func (g *GoCloak) DeleteRequiredAction(ctx context.Context, token string, realm string, alias string) error {
	const errMessage = "could not delete required action"

	if alias == "" {
		return errors.New("alias is required for deleting a required action")
	}
	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		Delete(g.getAdminRealmURL(realm, "authentication", "required-actions", alias))

	if err := checkForError(resp, err, errMessage); err != nil {
		return err
	}

	return err
}

// CreateClientScopesScopeMappingsClientRoles attaches a client role to a client scope (not client's scope)
func (g *GoCloak) CreateClientScopesScopeMappingsClientRoles(
	ctx context.Context, token, realm, idOfClientScope, idOfClient string, roles []Role,
) error {
	const errMessage = "could not create client-level roles to the client-scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Post(g.getAdminRealmURL(realm, "client-scopes", idOfClientScope, "scope-mappings", "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// GetClientScopesScopeMappingsClientRolesAvailable returns available (i.e. not attached via
// CreateClientScopesScopeMappingsClientRoles) client roles for a specific client, for a client scope
// (not client's scope).
func (g *GoCloak) GetClientScopesScopeMappingsClientRolesAvailable(ctx context.Context, token, realm, idOfClientScope, idOfClient string) ([]*Role, error) {
	const errMessage = "could not get available client-level roles with the client-scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", idOfClientScope, "scope-mappings", "clients", idOfClient, "available"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// GetClientScopesScopeMappingsClientRoles returns attached client roles for a specific client, for a client scope
// (not client's scope).
func (g *GoCloak) GetClientScopesScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClientScope, idOfClient string) ([]*Role, error) {
	const errMessage = "could not get client-level roles with the client-scope"

	var result []*Role

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "client-scopes", idOfClientScope, "scope-mappings", "clients", idOfClient))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteClientScopesScopeMappingsClientRoles removes attachment of client roles from a client scope
// (not client's scope).
func (g *GoCloak) DeleteClientScopesScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClientScope, idOfClient string, roles []Role) error {
	const errMessage = "could not delete client-level roles from the client-scope"

	resp, err := g.GetRequestWithBearerAuth(ctx, token).
		SetBody(roles).
		Delete(g.getAdminRealmURL(realm, "client-scopes", idOfClientScope, "scope-mappings", "clients", idOfClient))

	return checkForError(resp, err, errMessage)
}

// RevokeToken revokes the passed token. The token can either be an access or refresh token.
func (g *GoCloak) RevokeToken(ctx context.Context, realm, clientID, clientSecret, refreshToken string) error {
	const errMessage = "could not revoke token"

	resp, err := g.GetRequestWithBasicAuth(ctx, clientID, clientSecret).
		SetFormData(map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"token":         refreshToken,
		}).
		Post(g.getRealmURL(realm, g.Config.revokeEndpoint))

	return checkForError(resp, err, errMessage)
}

// UpdateUsersManagementPermissions updates the management permissions for users
func (g *GoCloak) UpdateUsersManagementPermissions(ctx context.Context, accessToken, realm string, managementPermissions ManagementPermissionRepresentation) (*ManagementPermissionRepresentation, error) {
	const errMessage = "could not update users management permissions"

	var result ManagementPermissionRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		SetBody(managementPermissions).
		Put(g.getAdminRealmURL(realm, "users-management-permissions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUsersManagementPermissions returns the management permissions for users
func (g *GoCloak) GetUsersManagementPermissions(ctx context.Context, accessToken, realm string) (*ManagementPermissionRepresentation, error) {
	const errMessage = "could not get users management permissions"

	var result ManagementPermissionRepresentation

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetResult(&result).
		Get(g.getAdminRealmURL(realm, "users-management-permissions"))

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}
