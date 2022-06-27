/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kcc // import "stash.kopano.io/kgol/kcc-go"

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	// DefaultURI is the default Kopano server URI to be used when no URI is
	// given when constructing a KCC or SOAP instance.
	DefaultURI = "http://127.0.0.1:236"
	// DefaultAppName is the default client app name as sent to the server.
	DefaultAppName = "kcc-go"
	// Version specifies the version string of this client implementation.
	Version = "0.0.0-dev"
	// ClientVersion specifies the version of this clients API implementation,
	ClientVersion = 8
)

var debug = false

func init() {
	uri := os.Getenv("KOPANO_SERVER_DEFAULT_URI")
	if uri != "" {
		DefaultURI = uri
	}
	debug = os.Getenv("KCC_GO_DEBUG") != ""
}

// A KCC is the client implementation base object containing the HTTP connection
// pool and other references to interface with a Kopano server via SOAP.
type KCC struct {
	Client       SOAPClient
	Capabilities KCFlag

	app [2]string
}

// NewKCC constructs a KCC instance with the provided URI. If no URI is passed,
// the current DefaultURI value will tbe used. Returns nil when the provided URI
// cannot be parsed or is not a valid kcc URI.
func NewKCC(uri *url.URL) *KCC {
	c, _ := NewKCCFromURI(uri)
	return c
}

// NewKCCWithClient constructs a KCC instance using the provided SOAPClient. Use
// this function if you need specific SOAPClient settings.
func NewKCCWithClient(client SOAPClient) *KCC {
	c := &KCC{
		app: [2]string{DefaultAppName, Version},

		Client:       client,
		Capabilities: DefaultClientCapabilities,
	}

	return c
}

// NewKCCFromURI constructs a KCC instance with the provided URI. If no URI is
// passed, the current DefaultURI value will tbe used.
func NewKCCFromURI(uri *url.URL) (*KCC, error) {
	if uri == nil {
		uri, _ = url.Parse(DefaultURI)
	}
	soap, err := NewSOAPClient(uri)
	if err != nil {
		return nil, err
	}

	return NewKCCWithClient(soap), nil
}

func (c *KCC) String() string {
	return fmt.Sprintf("KCC(%s)", c.Client)
}

// SetClientApp sets the clients app details as sent with requests to the
// accociated server.
func (c *KCC) SetClientApp(name, version string) error {
	c.app = [2]string{name, version}
	return nil
}

// Logon creates a session with the Kopano server using the provided credentials.
func (c *KCC) Logon(ctx context.Context, username, password string, logonFlags KCFlag) (*LogonResponse, error) {
	var b strings.Builder
	b.WriteString("<ns:logon><szUsername>")
	b.WriteString(xmlCharData(username).Escape())
	b.WriteString("</szUsername><szPassword>")
	b.WriteString(xmlCharData(password).Escape())
	b.WriteString("</szPassword><szImpersonateUser/><ulCapabilities>")
	b.WriteString(c.Capabilities.String())
	b.WriteString("</ulCapabilities><ulFlags>")
	b.WriteString(logonFlags.String())
	b.WriteString("</ulFlags><szClientApp>")
	b.WriteString(xmlCharData(c.app[0]).Escape())
	b.WriteString("</szClientApp><szClientAppVersion>")
	b.WriteString(xmlCharData(c.app[1]).Escape())
	b.WriteString("</szClientAppVersion><clientVersion>")
	b.WriteString(strconv.FormatInt(int64(ClientVersion), 10))
	b.WriteString("</clientVersion></ns:logon>")
	payload := b.String()

	var logonResponse LogonResponse
	err := c.Client.DoRequest(ctx, &payload, &logonResponse)

	return &logonResponse, err
}

// SSOLogon creates a session with the Kopano server using the provided credentials.
func (c *KCC) SSOLogon(ctx context.Context, prefix SSOType, username string, input []byte, sessionID KCSessionID, logonFlags KCFlag) (*LogonResponse, error) {
	if logonFlags != 0 {
		return nil, fmt.Errorf("logon flags are not support by sso logon")
	}

	// Add prefix value.
	lpInput := make([]byte, 0, len(prefix)+len(input))
	lpInput = append(lpInput, prefix.String()...)
	lpInput = append(lpInput, input...)

	// NOTE(longsleep): There is currently no way to specify flags when using
	// SSOLogon. This means, a new session is created when none was given and
	// the call will fail with error if the given session does not exist.
	var b strings.Builder
	b.WriteString("<ns:ssoLogon><szUsername>")
	b.WriteString(xmlCharData(username).Escape())
	b.WriteString("</szUsername><lpInput>")
	b.WriteString(base64.StdEncoding.EncodeToString(lpInput))
	b.WriteString("</lpInput><szImpersonateUser/><clientCaps>")
	b.WriteString(c.Capabilities.String())
	b.WriteString("</clientCaps><szClientApp>")
	b.WriteString(xmlCharData(c.app[0]).Escape())
	b.WriteString("</szClientApp><szClientAppVersion>")
	b.WriteString(xmlCharData(c.app[1]).Escape())
	b.WriteString("</szClientAppVersion><clientVersion>")
	b.WriteString(strconv.FormatInt(int64(ClientVersion), 10))
	b.WriteString("</clientVersion><ulSessionId>")
	b.WriteString(sessionID.String())
	b.WriteString("</ulSessionId></ns:ssoLogon>")
	payload := b.String()

	var logonResponse LogonResponse
	err := c.Client.DoRequest(ctx, &payload, &logonResponse)

	return &logonResponse, err
}

// Logoff terminates the provided session with the Kopano server.
func (c *KCC) Logoff(ctx context.Context, sessionID KCSessionID) (*LogoffResponse, error) {
	payload := "<ns:logoff><ulSessionId>" +
		sessionID.String() +
		"</ulSessionId></ns:logoff>"

	var logoffResponse LogoffResponse
	err := c.Client.DoRequest(ctx, &payload, &logoffResponse)

	return &logoffResponse, err
}

// ResolveUsername looks up the user ID of the provided username using the
// provided session.
func (c *KCC) ResolveUsername(ctx context.Context, username string, sessionID KCSessionID) (*ResolveUserResponse, error) {
	var b strings.Builder
	b.WriteString("<ns:resolveUsername><lpszUsername>")
	b.WriteString(xmlCharData(username).Escape())
	b.WriteString("</lpszUsername><ulSessionId>")
	b.WriteString(sessionID.String())
	b.WriteString("</ulSessionId></ns:resolveUsername>")
	payload := b.String()

	var resolveUserResponse ResolveUserResponse
	err := c.Client.DoRequest(ctx, &payload, &resolveUserResponse)

	return &resolveUserResponse, err
}

// GetUser fetches a user's detail meta data of the provided user Entry
// ID using the provided session.
func (c *KCC) GetUser(ctx context.Context, userEntryID string, sessionID KCSessionID) (*GetUserResponse, error) {
	var b strings.Builder
	b.WriteString("<ns:getUser><sUserId>")
	b.WriteString(userEntryID)
	b.WriteString("</sUserId><ulSessionId>")
	b.WriteString(sessionID.String())
	b.WriteString("</ulSessionId></ns:getUser>")
	payload := b.String()

	var getUserResponse GetUserResponse
	err := c.Client.DoRequest(ctx, &payload, &getUserResponse)

	return &getUserResponse, err
}

// ABResolveNames searches the AB for the provided props using the provided
// request data and flags.
func (c *KCC) ABResolveNames(ctx context.Context, props []PT, request map[PT]interface{}, requestFlags ABFlag, sessionID KCSessionID, resolveNamesFlags KCFlag) (*ABResolveNamesResponse, error) {
	var b strings.Builder
	b.WriteString("<ns:abResolveNames>")
	b.WriteString("<ulSessionId>")
	b.WriteString(sessionID.String())
	b.WriteString("</ulSessionId>")
	b.WriteString("<lpaPropTag SOAP-ENC:arrayType=\"xsd:unsignedInt[")
	b.WriteString(strconv.FormatUint(uint64(len(props)), 64))
	b.WriteString("]\">")
	for _, prop := range props {
		b.WriteString("<item>")
		b.WriteString(prop.String())
		b.WriteString("</item>")
	}
	b.WriteString("</lpaPropTag>")
	b.WriteString("<lpsRowSet SOAP-ENC:arrayType=\"propVal[][")
	b.WriteString(strconv.FormatUint(uint64(len(request)), 64))
	b.WriteString("]\">")
	for prop, value := range request {
		b.WriteString("<item SOAP-ENC:arrayType=\"propVal[1]\">")
		b.WriteString("<item>")
		b.WriteString("<ulPropTag>")
		b.WriteString(prop.String())
		b.WriteString("</ulPropTag>")
		switch tv := value.(type) {
		case string:
			b.WriteString("<lpszA>")
			b.WriteString(xmlCharData(tv).Escape())
			b.WriteString("</lpszA>")
		default:
			return nil, fmt.Errorf("unsupported type in request map value: %v", value)
		}
		b.WriteString("</item>")
		b.WriteString("</item>")
	}
	b.WriteString("</lpsRowSet>")
	b.WriteString("<lpaFlags>")
	b.WriteString("<item>")
	b.WriteString(requestFlags.String())
	b.WriteString("</item>")
	b.WriteString("</lpaFlags>")
	b.WriteString("<ulFlags>")
	b.WriteString(resolveNamesFlags.String())
	b.WriteString("</ulFlags>")
	b.WriteString("</ns:abResolveNames>")
	payload := b.String()

	var abResolveNamesResponse ABResolveNamesResponse
	err := c.Client.DoRequest(ctx, &payload, &abResolveNamesResponse)

	return &abResolveNamesResponse, err
}
