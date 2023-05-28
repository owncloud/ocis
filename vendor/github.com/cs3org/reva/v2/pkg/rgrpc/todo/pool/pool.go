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
	"github.com/cs3org/reva/v2/pkg/registry"
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
	return getClient[gateway.GatewayAPIClient](endpoint, &gatewayProviders, func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
		return gateway.NewGatewayAPIClient(cc)
	}, opts...)
}

// GetUserProviderServiceClient returns a UserProviderServiceClient.
func GetUserProviderServiceClient(endpoint string, opts ...Option) (user.UserAPIClient, error) {
	return getClient[user.UserAPIClient](endpoint, &userProviders, func(cc *grpc.ClientConn) user.UserAPIClient {
		return user.NewUserAPIClient(cc)
	}, opts...)
}

// GetGroupProviderServiceClient returns a GroupProviderServiceClient.
func GetGroupProviderServiceClient(endpoint string, opts ...Option) (group.GroupAPIClient, error) {
	return getClient[group.GroupAPIClient](endpoint, &groupProviders, func(cc *grpc.ClientConn) group.GroupAPIClient {
		return group.NewGroupAPIClient(cc)
	}, opts...)
}

// GetStorageProviderServiceClient returns a StorageProviderServiceClient.
func GetStorageProviderServiceClient(endpoint string, opts ...Option) (storageprovider.ProviderAPIClient, error) {
	return getClient[storageprovider.ProviderAPIClient](endpoint, &storageProviders, func(cc *grpc.ClientConn) storageprovider.ProviderAPIClient {
		return storageprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetAuthRegistryServiceClient returns a new AuthRegistryServiceClient.
func GetAuthRegistryServiceClient(endpoint string, opts ...Option) (authregistry.RegistryAPIClient, error) {
	return getClient[authregistry.RegistryAPIClient](endpoint, &authRegistries, func(cc *grpc.ClientConn) authregistry.RegistryAPIClient {
		return authregistry.NewRegistryAPIClient(cc)
	}, opts...)
}

// GetAuthProviderServiceClient returns a new AuthProviderServiceClient.
func GetAuthProviderServiceClient(endpoint string, opts ...Option) (authprovider.ProviderAPIClient, error) {
	return getClient[authprovider.ProviderAPIClient](endpoint, &authProviders, func(cc *grpc.ClientConn) authprovider.ProviderAPIClient {
		return authprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetAppAuthProviderServiceClient returns a new AppAuthProviderServiceClient.
func GetAppAuthProviderServiceClient(endpoint string, opts ...Option) (applicationauth.ApplicationsAPIClient, error) {
	return getClient[applicationauth.ApplicationsAPIClient](endpoint, &appAuthProviders, func(cc *grpc.ClientConn) applicationauth.ApplicationsAPIClient {
		return applicationauth.NewApplicationsAPIClient(cc)
	}, opts...)
}

// GetUserShareProviderClient returns a new UserShareProviderClient.
func GetUserShareProviderClient(endpoint string, opts ...Option) (collaboration.CollaborationAPIClient, error) {
	return getClient[collaboration.CollaborationAPIClient](endpoint, &userShareProviders, func(cc *grpc.ClientConn) collaboration.CollaborationAPIClient {
		return collaboration.NewCollaborationAPIClient(cc)
	}, opts...)
}

// GetOCMShareProviderClient returns a new OCMShareProviderClient.
func GetOCMShareProviderClient(endpoint string, opts ...Option) (ocm.OcmAPIClient, error) {
	return getClient[ocm.OcmAPIClient](endpoint, &ocmShareProviders, func(cc *grpc.ClientConn) ocm.OcmAPIClient {
		return ocm.NewOcmAPIClient(cc)
	}, opts...)
}

// GetOCMInviteManagerClient returns a new OCMInviteManagerClient.
func GetOCMInviteManagerClient(endpoint string, opts ...Option) (invitepb.InviteAPIClient, error) {
	return getClient[invitepb.InviteAPIClient](endpoint, &ocmInviteManagers, func(cc *grpc.ClientConn) invitepb.InviteAPIClient {
		return invitepb.NewInviteAPIClient(cc)
	}, opts...)
}

// GetPublicShareProviderClient returns a new PublicShareProviderClient.
func GetPublicShareProviderClient(endpoint string, opts ...Option) (link.LinkAPIClient, error) {
	return getClient[link.LinkAPIClient](endpoint, &publicShareProviders, func(cc *grpc.ClientConn) link.LinkAPIClient {
		return link.NewLinkAPIClient(cc)
	}, opts...)
}

// GetPreferencesClient returns a new PreferencesClient.
func GetPreferencesClient(endpoint string, opts ...Option) (preferences.PreferencesAPIClient, error) {
	return getClient[preferences.PreferencesAPIClient](endpoint, &preferencesProviders, func(cc *grpc.ClientConn) preferences.PreferencesAPIClient {
		return preferences.NewPreferencesAPIClient(cc)
	}, opts...)
}

// GetPermissionsClient returns a new PermissionsClient.
func GetPermissionsClient(endpoint string, opts ...Option) (permissions.PermissionsAPIClient, error) {
	return getClient[permissions.PermissionsAPIClient](endpoint, &permissionsProviders, func(cc *grpc.ClientConn) permissions.PermissionsAPIClient {
		return permissions.NewPermissionsAPIClient(cc)
	}, opts...)
}

// GetAppRegistryClient returns a new AppRegistryClient.
func GetAppRegistryClient(endpoint string, opts ...Option) (appregistry.RegistryAPIClient, error) {
	return getClient[appregistry.RegistryAPIClient](endpoint, &appRegistries, func(cc *grpc.ClientConn) appregistry.RegistryAPIClient {
		return appregistry.NewRegistryAPIClient(cc)
	}, opts...)
}

// GetAppProviderClient returns a new AppRegistryClient.
func GetAppProviderClient(endpoint string, opts ...Option) (appprovider.ProviderAPIClient, error) {
	return getClient[appprovider.ProviderAPIClient](endpoint, &appProviders, func(cc *grpc.ClientConn) appprovider.ProviderAPIClient {
		return appprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetStorageRegistryClient returns a new StorageRegistryClient.
func GetStorageRegistryClient(endpoint string, opts ...Option) (storageregistry.RegistryAPIClient, error) {
	return getClient[storageregistry.RegistryAPIClient](endpoint, &storageRegistries, func(cc *grpc.ClientConn) storageregistry.RegistryAPIClient {
		return storageregistry.NewRegistryAPIClient(cc)
	}, opts...)
}

// GetOCMProviderAuthorizerClient returns a new OCMProviderAuthorizerClient.
func GetOCMProviderAuthorizerClient(endpoint string, opts ...Option) (ocmprovider.ProviderAPIClient, error) {
	return getClient[ocmprovider.ProviderAPIClient](endpoint, &ocmProviderAuthorizers, func(cc *grpc.ClientConn) ocmprovider.ProviderAPIClient {
		return ocmprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetOCMCoreClient returns a new OCMCoreClient.
func GetOCMCoreClient(endpoint string, opts ...Option) (ocmcore.OcmCoreAPIClient, error) {
	return getClient[ocmcore.OcmCoreAPIClient](endpoint, &ocmCores, func(cc *grpc.ClientConn) ocmcore.OcmCoreAPIClient {
		return ocmcore.NewOcmCoreAPIClient(cc)
	}, opts...)
}

// GetDataTxClient returns a new DataTxClient.
func GetDataTxClient(endpoint string, opts ...Option) (datatx.TxAPIClient, error) {
	return getClient[datatx.TxAPIClient](endpoint, &dataTxs, func(cc *grpc.ClientConn) datatx.TxAPIClient {
		return datatx.NewTxAPIClient(cc)
	}, opts...)
}

func getClient[T any](endpoint string, p *provider, cf func(cc *grpc.ClientConn) T, opts ...Option) (T, error) {
	services, _ := registry.GetServiceByAddress(endpoint)
	address, err := registry.GetNodeAddress(services)
	if err == nil && endpoint != "" {
		endpoint = address
	}

	p.m.Lock()
	defer p.m.Unlock()

	if c, ok := p.conn[endpoint]; ok {
		return c.(T), nil
	}

	conn, err := NewConn(endpoint, opts...)
	if err != nil {
		return *new(T), err
	}

	v := cf(conn)
	p.conn[endpoint] = v
	return v, nil
}
