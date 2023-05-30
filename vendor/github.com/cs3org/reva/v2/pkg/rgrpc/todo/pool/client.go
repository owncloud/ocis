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
	"google.golang.org/grpc"
)

// GetGatewayServiceClient returns a GatewayServiceClient.
func GetGatewayServiceClient(id string, opts ...Option) (gateway.GatewayAPIClient, error) {
	return getClient[gateway.GatewayAPIClient](id, &gatewayProviders, func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
		return gateway.NewGatewayAPIClient(cc)
	}, opts...)
}

// GetUserProviderServiceClient returns a UserProviderServiceClient.
func GetUserProviderServiceClient(id string, opts ...Option) (user.UserAPIClient, error) {
	return getClient[user.UserAPIClient](id, &userProviders, func(cc *grpc.ClientConn) user.UserAPIClient {
		return user.NewUserAPIClient(cc)
	}, opts...)
}

// GetGroupProviderServiceClient returns a GroupProviderServiceClient.
func GetGroupProviderServiceClient(id string, opts ...Option) (group.GroupAPIClient, error) {
	return getClient[group.GroupAPIClient](id, &groupProviders, func(cc *grpc.ClientConn) group.GroupAPIClient {
		return group.NewGroupAPIClient(cc)
	}, opts...)
}

// GetStorageProviderServiceClient returns a StorageProviderServiceClient.
func GetStorageProviderServiceClient(id string, opts ...Option) (storageprovider.ProviderAPIClient, error) {
	return getClient[storageprovider.ProviderAPIClient](id, &storageProviders, func(cc *grpc.ClientConn) storageprovider.ProviderAPIClient {
		return storageprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetAuthRegistryServiceClient returns a new AuthRegistryServiceClient.
func GetAuthRegistryServiceClient(id string, opts ...Option) (authregistry.RegistryAPIClient, error) {
	return getClient[authregistry.RegistryAPIClient](id, &authRegistries, func(cc *grpc.ClientConn) authregistry.RegistryAPIClient {
		return authregistry.NewRegistryAPIClient(cc)
	}, opts...)
}

// GetAuthProviderServiceClient returns a new AuthProviderServiceClient.
func GetAuthProviderServiceClient(id string, opts ...Option) (authprovider.ProviderAPIClient, error) {
	return getClient[authprovider.ProviderAPIClient](id, &authProviders, func(cc *grpc.ClientConn) authprovider.ProviderAPIClient {
		return authprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetAppAuthProviderServiceClient returns a new AppAuthProviderServiceClient.
func GetAppAuthProviderServiceClient(id string, opts ...Option) (applicationauth.ApplicationsAPIClient, error) {
	return getClient[applicationauth.ApplicationsAPIClient](id, &appAuthProviders, func(cc *grpc.ClientConn) applicationauth.ApplicationsAPIClient {
		return applicationauth.NewApplicationsAPIClient(cc)
	}, opts...)
}

// GetUserShareProviderClient returns a new UserShareProviderClient.
func GetUserShareProviderClient(id string, opts ...Option) (collaboration.CollaborationAPIClient, error) {
	return getClient[collaboration.CollaborationAPIClient](id, &userShareProviders, func(cc *grpc.ClientConn) collaboration.CollaborationAPIClient {
		return collaboration.NewCollaborationAPIClient(cc)
	}, opts...)
}

// GetOCMShareProviderClient returns a new OCMShareProviderClient.
func GetOCMShareProviderClient(id string, opts ...Option) (ocm.OcmAPIClient, error) {
	return getClient[ocm.OcmAPIClient](id, &ocmShareProviders, func(cc *grpc.ClientConn) ocm.OcmAPIClient {
		return ocm.NewOcmAPIClient(cc)
	}, opts...)
}

// GetOCMInviteManagerClient returns a new OCMInviteManagerClient.
func GetOCMInviteManagerClient(id string, opts ...Option) (invitepb.InviteAPIClient, error) {
	return getClient[invitepb.InviteAPIClient](id, &ocmInviteManagers, func(cc *grpc.ClientConn) invitepb.InviteAPIClient {
		return invitepb.NewInviteAPIClient(cc)
	}, opts...)
}

// GetPublicShareProviderClient returns a new PublicShareProviderClient.
func GetPublicShareProviderClient(id string, opts ...Option) (link.LinkAPIClient, error) {
	return getClient[link.LinkAPIClient](id, &publicShareProviders, func(cc *grpc.ClientConn) link.LinkAPIClient {
		return link.NewLinkAPIClient(cc)
	}, opts...)
}

// GetPreferencesClient returns a new PreferencesClient.
func GetPreferencesClient(id string, opts ...Option) (preferences.PreferencesAPIClient, error) {
	return getClient[preferences.PreferencesAPIClient](id, &preferencesProviders, func(cc *grpc.ClientConn) preferences.PreferencesAPIClient {
		return preferences.NewPreferencesAPIClient(cc)
	}, opts...)
}

// GetPermissionsClient returns a new PermissionsClient.
func GetPermissionsClient(id string, opts ...Option) (permissions.PermissionsAPIClient, error) {
	return getClient[permissions.PermissionsAPIClient](id, &permissionsProviders, func(cc *grpc.ClientConn) permissions.PermissionsAPIClient {
		return permissions.NewPermissionsAPIClient(cc)
	}, opts...)
}

// GetAppRegistryClient returns a new AppRegistryClient.
func GetAppRegistryClient(id string, opts ...Option) (appregistry.RegistryAPIClient, error) {
	return getClient[appregistry.RegistryAPIClient](id, &appRegistries, func(cc *grpc.ClientConn) appregistry.RegistryAPIClient {
		return appregistry.NewRegistryAPIClient(cc)
	}, opts...)
}

// GetAppProviderClient returns a new AppRegistryClient.
func GetAppProviderClient(id string, opts ...Option) (appprovider.ProviderAPIClient, error) {
	return getClient[appprovider.ProviderAPIClient](id, &appProviders, func(cc *grpc.ClientConn) appprovider.ProviderAPIClient {
		return appprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetStorageRegistryClient returns a new StorageRegistryClient.
func GetStorageRegistryClient(id string, opts ...Option) (storageregistry.RegistryAPIClient, error) {
	return getClient[storageregistry.RegistryAPIClient](id, &storageRegistries, func(cc *grpc.ClientConn) storageregistry.RegistryAPIClient {
		return storageregistry.NewRegistryAPIClient(cc)
	}, opts...)
}

// GetOCMProviderAuthorizerClient returns a new OCMProviderAuthorizerClient.
func GetOCMProviderAuthorizerClient(id string, opts ...Option) (ocmprovider.ProviderAPIClient, error) {
	return getClient[ocmprovider.ProviderAPIClient](id, &ocmProviderAuthorizers, func(cc *grpc.ClientConn) ocmprovider.ProviderAPIClient {
		return ocmprovider.NewProviderAPIClient(cc)
	}, opts...)
}

// GetOCMCoreClient returns a new OCMCoreClient.
func GetOCMCoreClient(id string, opts ...Option) (ocmcore.OcmCoreAPIClient, error) {
	return getClient[ocmcore.OcmCoreAPIClient](id, &ocmCores, func(cc *grpc.ClientConn) ocmcore.OcmCoreAPIClient {
		return ocmcore.NewOcmCoreAPIClient(cc)
	}, opts...)
}

// GetDataTxClient returns a new DataTxClient.
func GetDataTxClient(id string, opts ...Option) (datatx.TxAPIClient, error) {
	return getClient[datatx.TxAPIClient](id, &dataTxs, func(cc *grpc.ClientConn) datatx.TxAPIClient {
		return datatx.NewTxAPIClient(cc)
	}, opts...)
}

func getClient[T any](id string, p *provider, cf func(cc *grpc.ClientConn) T, opts ...Option) (T, error) {
	services, _ := registry.DiscoverServices(id)
	address, err := registry.GetNodeAddress(services)
	if err == nil && id != "" {
		id = address
	}

	p.m.Lock()
	defer p.m.Unlock()

	if c, ok := p.conn[id]; ok {
		return c.(T), nil
	}

	conn, err := NewConn(id, opts...)
	if err != nil {
		return *new(T), err
	}

	v := cf(conn)
	p.conn[id] = v
	return v, nil
}
