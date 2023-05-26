// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package pool

import (
	"crypto/tls"
	"fmt"
	"sync"

	appprovider "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	appregistry "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	applicationauth "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	authprovider "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	authregistry "github.com/cs3org/go-cs3apis/cs3/auth/registry/v1beta1"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmcore "github.com/cs3org/go-cs3apis/cs3/ocm/core/v1beta1"
	invitepb "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	preferences "github.com/cs3org/go-cs3apis/cs3/preferences/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageregistry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	datatx "github.com/cs3org/go-cs3apis/cs3/tx/v1beta1"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	rtrace "github.com/cs3org/reva/v2/pkg/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type provider struct {
	m    sync.Mutex
	conn map[string]interface{}
}

func newProvider() provider {
	return provider{
		sync.Mutex{},
		make(map[string]interface{}),
	}
}

// TLSMode represents TLS mode for the clients
type TLSMode int

const (
	// TLSOff completely disables transport security
	TLSOff TLSMode = iota
	// TLSOn enables transport security
	TLSOn
	// TLSInsecure enables transport security, but disables the verification of the
	// server certificate
	TLSInsecure
)

// ClientOptions represent additional options (e.g. tls settings) for the grpc clients
type ClientOptions struct {
	tlsMode        TLSMode
	caCert         string
	tracerProvider trace.TracerProvider
}

// Option is used to pass client options
type Option func(opts *ClientOptions)

// TODO(labkode): is concurrent access to the maps safe?
// var storageProviders = map[string]storageprovider.ProviderAPIClient{}
var (
	storageProviders       = newProvider()
	authProviders          = newProvider()
	appAuthProviders       = newProvider()
	authRegistries         = newProvider()
	userShareProviders     = newProvider()
	ocmShareProviders      = newProvider()
	ocmInviteManagers      = newProvider()
	ocmProviderAuthorizers = newProvider()
	ocmCores               = newProvider()
	publicShareProviders   = newProvider()
	preferencesProviders   = newProvider()
	permissionsProviders   = newProvider()
	appRegistries          = newProvider()
	appProviders           = newProvider()
	storageRegistries      = newProvider()
	gatewayProviders       = newProvider()
	userProviders          = newProvider()
	groupProviders         = newProvider()
	dataTxs                = newProvider()
	maxCallRecvMsgSize     = 10240000
)

// StringToTLSMode converts the supply string into the equivalent TLSMode constant
func StringToTLSMode(m string) (TLSMode, error) {
	switch m {
	case "off", "":
		return TLSOff, nil
	case "insecure":
		return TLSInsecure, nil
	case "on":
		return TLSOn, nil
	default:
		return TLSOff, fmt.Errorf("unknown TLS mode: '%s'. Valid values are 'on', 'off' and 'insecure'", m)
	}
}

func (o *ClientOptions) init() error {
	// default to shared settings
	sharedOpt := sharedconf.GRPCClientOptions()
	var err error

	if o.tlsMode, err = StringToTLSMode(sharedOpt.TLSMode); err != nil {
		return err
	}
	o.caCert = sharedOpt.CACertFile
	o.tracerProvider = rtrace.DefaultProvider()
	return nil
}

// WithTLSMode allows to set the TLSMode option for grpc clients
func WithTLSMode(v TLSMode) Option {
	return func(o *ClientOptions) {
		o.tlsMode = v
	}
}

// WithTLSCACert allows to set the CA Certificate for grpc clients
func WithTLSCACert(v string) Option {
	return func(o *ClientOptions) {
		o.caCert = v
	}
}

// WithTracerProvider allows to set the opentelemetry tracer provider for grpc clients
func WithTracerProvider(v trace.TracerProvider) Option {
	return func(o *ClientOptions) {
		o.tracerProvider = v
	}
}

// NewConn creates a new connection to a grpc server
// with open census tracing support.
// TODO(labkode): make grpc tls configurable.
// TODO make maxCallRecvMsgSize configurable, raised from the default 4MB to be able to list 10k files
func NewConn(endpoint string, opts ...Option) (*grpc.ClientConn, error) {

	options := ClientOptions{}
	if err := options.init(); err != nil {
		return nil, err
	}

	// then overwrite with supplied options
	for _, opt := range opts {
		opt(&options)
	}

	var cred credentials.TransportCredentials
	switch options.tlsMode {
	case TLSOff:
		cred = insecure.NewCredentials()
	case TLSInsecure:
		tlsConfig := tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		}
		cred = credentials.NewTLS(&tlsConfig)
	case TLSOn:
		if options.caCert != "" {
			var err error
			if cred, err = credentials.NewClientTLSFromFile(options.caCert, ""); err != nil {
				return nil, err
			}
		} else {
			// Use system's cert pool
			cred = credentials.NewTLS(&tls.Config{})
		}
	}

	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(cred),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize),
		),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(
			otelgrpc.WithTracerProvider(
				options.tracerProvider,
			),
			otelgrpc.WithPropagators(
				rtrace.Propagator,
			),
		)),
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(
					options.tracerProvider,
				),
				otelgrpc.WithPropagators(
					rtrace.Propagator,
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// GetGatewayServiceClient returns a GatewayServiceClient.
func GetGatewayServiceClient(endpoint string, opts ...Option) (gateway.GatewayAPIClient, error) {
	gatewayProviders.m.Lock()
	defer gatewayProviders.m.Unlock()

	if val, ok := gatewayProviders.conn[endpoint]; ok {
		return val.(gateway.GatewayAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := gateway.NewGatewayAPIClient(conn)
	gatewayProviders.conn[endpoint] = v

	return v, nil
}

// GetUserProviderServiceClient returns a UserProviderServiceClient.
func GetUserProviderServiceClient(endpoint string, opts ...Option) (user.UserAPIClient, error) {
	userProviders.m.Lock()
	defer userProviders.m.Unlock()

	if val, ok := userProviders.conn[endpoint]; ok {
		return val.(user.UserAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := user.NewUserAPIClient(conn)
	userProviders.conn[endpoint] = v
	return v, nil
}

// GetGroupProviderServiceClient returns a GroupProviderServiceClient.
func GetGroupProviderServiceClient(endpoint string, opts ...Option) (group.GroupAPIClient, error) {
	groupProviders.m.Lock()
	defer groupProviders.m.Unlock()

	if val, ok := groupProviders.conn[endpoint]; ok {
		return val.(group.GroupAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := group.NewGroupAPIClient(conn)
	groupProviders.conn[endpoint] = v
	return v, nil
}

// GetStorageProviderServiceClient returns a StorageProviderServiceClient.
func GetStorageProviderServiceClient(endpoint string, opts ...Option) (storageprovider.ProviderAPIClient, error) {
	storageProviders.m.Lock()
	defer storageProviders.m.Unlock()

	if c, ok := storageProviders.conn[endpoint]; ok {
		return c.(storageprovider.ProviderAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := storageprovider.NewProviderAPIClient(conn)
	storageProviders.conn[endpoint] = v
	return v, nil
}

// GetAuthRegistryServiceClient returns a new AuthRegistryServiceClient.
func GetAuthRegistryServiceClient(endpoint string, opts ...Option) (authregistry.RegistryAPIClient, error) {
	authRegistries.m.Lock()
	defer authRegistries.m.Unlock()

	// if there is already a connection to this node, use it.
	if c, ok := authRegistries.conn[endpoint]; ok {
		return c.(authregistry.RegistryAPIClient), nil
	}

	// if not, create a new connection
	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	// and memoize it
	v := authregistry.NewRegistryAPIClient(conn)
	authRegistries.conn[endpoint] = v
	return v, nil
}

// GetAuthProviderServiceClient returns a new AuthProviderServiceClient.
func GetAuthProviderServiceClient(endpoint string, opts ...Option) (authprovider.ProviderAPIClient, error) {
	authProviders.m.Lock()
	defer authProviders.m.Unlock()

	if c, ok := authProviders.conn[endpoint]; ok {
		return c.(authprovider.ProviderAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := authprovider.NewProviderAPIClient(conn)
	authProviders.conn[endpoint] = v
	return v, nil
}

// GetAppAuthProviderServiceClient returns a new AppAuthProviderServiceClient.
func GetAppAuthProviderServiceClient(endpoint string, opts ...Option) (applicationauth.ApplicationsAPIClient, error) {
	appAuthProviders.m.Lock()
	defer appAuthProviders.m.Unlock()

	if c, ok := appAuthProviders.conn[endpoint]; ok {
		return c.(applicationauth.ApplicationsAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := applicationauth.NewApplicationsAPIClient(conn)
	appAuthProviders.conn[endpoint] = v
	return v, nil
}

// GetUserShareProviderClient returns a new UserShareProviderClient.
func GetUserShareProviderClient(endpoint string, opts ...Option) (collaboration.CollaborationAPIClient, error) {
	userShareProviders.m.Lock()
	defer userShareProviders.m.Unlock()

	if c, ok := userShareProviders.conn[endpoint]; ok {
		return c.(collaboration.CollaborationAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := collaboration.NewCollaborationAPIClient(conn)
	userShareProviders.conn[endpoint] = v
	return v, nil
}

// GetOCMShareProviderClient returns a new OCMShareProviderClient.
func GetOCMShareProviderClient(endpoint string, opts ...Option) (ocm.OcmAPIClient, error) {
	ocmShareProviders.m.Lock()
	defer ocmShareProviders.m.Unlock()

	if c, ok := ocmShareProviders.conn[endpoint]; ok {
		return c.(ocm.OcmAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := ocm.NewOcmAPIClient(conn)
	ocmShareProviders.conn[endpoint] = v
	return v, nil
}

// GetOCMInviteManagerClient returns a new OCMInviteManagerClient.
func GetOCMInviteManagerClient(endpoint string, opts ...Option) (invitepb.InviteAPIClient, error) {
	ocmInviteManagers.m.Lock()
	defer ocmInviteManagers.m.Unlock()

	if c, ok := ocmInviteManagers.conn[endpoint]; ok {
		return c.(invitepb.InviteAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := invitepb.NewInviteAPIClient(conn)
	ocmInviteManagers.conn[endpoint] = v
	return v, nil
}

// GetPublicShareProviderClient returns a new PublicShareProviderClient.
func GetPublicShareProviderClient(endpoint string, opts ...Option) (link.LinkAPIClient, error) {
	publicShareProviders.m.Lock()
	defer publicShareProviders.m.Unlock()

	if c, ok := publicShareProviders.conn[endpoint]; ok {
		return c.(link.LinkAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := link.NewLinkAPIClient(conn)
	publicShareProviders.conn[endpoint] = v
	return v, nil
}

// GetPreferencesClient returns a new PreferencesClient.
func GetPreferencesClient(endpoint string, opts ...Option) (preferences.PreferencesAPIClient, error) {
	preferencesProviders.m.Lock()
	defer preferencesProviders.m.Unlock()

	if c, ok := preferencesProviders.conn[endpoint]; ok {
		return c.(preferences.PreferencesAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := preferences.NewPreferencesAPIClient(conn)
	preferencesProviders.conn[endpoint] = v
	return v, nil
}

// GetPermissionsClient returns a new PermissionsClient.
func GetPermissionsClient(endpoint string, opts ...Option) (permissions.PermissionsAPIClient, error) {
	permissionsProviders.m.Lock()
	defer permissionsProviders.m.Unlock()

	if c, ok := permissionsProviders.conn[endpoint]; ok {
		return c.(permissions.PermissionsAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := permissions.NewPermissionsAPIClient(conn)
	permissionsProviders.conn[endpoint] = v
	return v, nil
}

// GetAppRegistryClient returns a new AppRegistryClient.
func GetAppRegistryClient(endpoint string, opts ...Option) (appregistry.RegistryAPIClient, error) {
	appRegistries.m.Lock()
	defer appRegistries.m.Unlock()

	if c, ok := appRegistries.conn[endpoint]; ok {
		return c.(appregistry.RegistryAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := appregistry.NewRegistryAPIClient(conn)
	appRegistries.conn[endpoint] = v
	return v, nil
}

// GetAppProviderClient returns a new AppRegistryClient.
func GetAppProviderClient(endpoint string, opts ...Option) (appprovider.ProviderAPIClient, error) {
	appProviders.m.Lock()
	defer appProviders.m.Unlock()

	if c, ok := appProviders.conn[endpoint]; ok {
		return c.(appprovider.ProviderAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := appprovider.NewProviderAPIClient(conn)
	appProviders.conn[endpoint] = v
	return v, nil
}

// GetStorageRegistryClient returns a new StorageRegistryClient.
func GetStorageRegistryClient(endpoint string, opts ...Option) (storageregistry.RegistryAPIClient, error) {
	storageRegistries.m.Lock()
	defer storageRegistries.m.Unlock()

	if c, ok := storageRegistries.conn[endpoint]; ok {
		return c.(storageregistry.RegistryAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := storageregistry.NewRegistryAPIClient(conn)
	storageRegistries.conn[endpoint] = v
	return v, nil
}

// GetOCMProviderAuthorizerClient returns a new OCMProviderAuthorizerClient.
func GetOCMProviderAuthorizerClient(endpoint string, opts ...Option) (ocmprovider.ProviderAPIClient, error) {
	ocmProviderAuthorizers.m.Lock()
	defer ocmProviderAuthorizers.m.Unlock()

	if c, ok := ocmProviderAuthorizers.conn[endpoint]; ok {
		return c.(ocmprovider.ProviderAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := ocmprovider.NewProviderAPIClient(conn)
	ocmProviderAuthorizers.conn[endpoint] = v
	return v, nil
}

// GetOCMCoreClient returns a new OCMCoreClient.
func GetOCMCoreClient(endpoint string, opts ...Option) (ocmcore.OcmCoreAPIClient, error) {
	ocmCores.m.Lock()
	defer ocmCores.m.Unlock()

	if c, ok := ocmCores.conn[endpoint]; ok {
		return c.(ocmcore.OcmCoreAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := ocmcore.NewOcmCoreAPIClient(conn)
	ocmCores.conn[endpoint] = v
	return v, nil
}

// GetDataTxClient returns a new DataTxClient.
func GetDataTxClient(endpoint string, opts ...Option) (datatx.TxAPIClient, error) {
	dataTxs.m.Lock()
	defer dataTxs.m.Unlock()

	if c, ok := dataTxs.conn[endpoint]; ok {
		return c.(datatx.TxAPIClient), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return nil, err
	}

	v := datatx.NewTxAPIClient(conn)
	dataTxs.conn[endpoint] = v
	return v, nil
}

// getEndpointByName resolve service names to ip addresses present on the registry.
//	func getEndpointByName(name string) (string, error) {
//		if services, err := utils.GlobalRegistry.GetService(name); err == nil {
//			if len(services) > 0 {
//				for i := range services {
//					for j := range services[i].Nodes() {
//						// return the first one. This MUST be improved upon with selectors.
//						return services[i].Nodes()[j].Address(), nil
//					}
//				}
//			}
//		}
//
//		return "", fmt.Errorf("could not get service by name: %v", name)
//	}
