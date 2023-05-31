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
	"fmt"
	"sync"

	appProvider "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	appRegistry "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	authApplication "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	authProvider "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	authRegistry "github.com/cs3org/go-cs3apis/cs3/auth/registry/v1beta1"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	identityGroup "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	identityUser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmCore "github.com/cs3org/go-cs3apis/cs3/ocm/core/v1beta1"
	ocmProvider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	preferences "github.com/cs3org/go-cs3apis/cs3/preferences/v1beta1"
	sharingCollaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	sharingLink "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	sharingOCM "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageRegistry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	tx "github.com/cs3org/go-cs3apis/cs3/tx/v1beta1"
	"github.com/cs3org/reva/v2/pkg/registry"
	"github.com/pkg/errors"
	mRegistry "go-micro.dev/v4/registry"
	"google.golang.org/grpc"
)

type Selectable[T any] interface {
	Next(opts ...Option) (T, error)
}

type Selector[T any] struct {
	id            string
	clientFactory func(cc *grpc.ClientConn) T
	clientMap     sync.Map
	options       []Option
}

func (s *Selector[T]) Next(opts ...Option) (T, error) {
	options := ClientOptions{}
	// first use selector options
	for _, opt := range s.options {
		opt(&options)
	}
	// then overwrite with supplied
	for _, opt := range opts {
		opt(&options)
	}
	var services []*mRegistry.Service
	if options.registry != nil {
		services, _ = options.registry.GetService(s.id)
	} else {
		services, _ = registry.DiscoverServices(s.id)
	}
	address, err := registry.GetNodeAddress(services)
	if err != nil || address == "" {
		return *new(T), errors.Wrap(err, fmt.Sprintf("could not get node addresses for %s", s.id))
	}

	existingClient, ok := s.clientMap.Load(address)
	if ok {
		return existingClient.(T), nil
	}

	conn, err := NewConn(address, append(s.options, opts...)...)

	if err != nil {
		return *new(T), errors.Wrap(err, fmt.Sprintf("could not create connection for %s to %s", s.id, address))
	}

	newClient := s.clientFactory(conn)
	s.clientMap.Store(address, newClient)

	return newClient, nil
}

// GatewaySelector returns a Selector[gateway.GatewayAPIClient].
func GatewaySelector(id string, options ...Option) (*Selector[gateway.GatewayAPIClient], error) {
	return &Selector[gateway.GatewayAPIClient]{
		id:            id,
		clientFactory: gateway.NewGatewayAPIClient,
		options:       options,
	}, nil
}

// IdentityUserSelector returns a Selector[identityUser.UserAPIClient].
func IdentityUserSelector(id string, options ...Option) (*Selector[identityUser.UserAPIClient], error) {
	return &Selector[identityUser.UserAPIClient]{
		id:            id,
		clientFactory: identityUser.NewUserAPIClient,
		options:       options,
	}, nil
}

// IdentityGroupSelector returns a Selector[identityGroup.GroupAPIClient].
func IdentityGroupSelector(id string, options ...Option) (*Selector[identityGroup.GroupAPIClient], error) {
	return &Selector[identityGroup.GroupAPIClient]{
		id:            id,
		clientFactory: identityGroup.NewGroupAPIClient,
		options:       options,
	}, nil
}

// StorageProviderSelector returns a Selector[storageProvider.ProviderAPIClient].
func StorageProviderSelector(id string, options ...Option) (*Selector[storageProvider.ProviderAPIClient], error) {
	return &Selector[storageProvider.ProviderAPIClient]{
		id:            id,
		clientFactory: storageProvider.NewProviderAPIClient,
		options:       options,
	}, nil
}

// AuthRegistrySelector returns a Selector[authRegistry.RegistryAPIClient].
func AuthRegistrySelector(id string, options ...Option) (*Selector[authRegistry.RegistryAPIClient], error) {
	return &Selector[authRegistry.RegistryAPIClient]{
		id:            id,
		clientFactory: authRegistry.NewRegistryAPIClient,
		options:       options,
	}, nil
}

// AuthProviderSelector returns a Selector[authProvider.RegistryAPIClient].
func AuthProviderSelector(id string, options ...Option) (*Selector[authProvider.ProviderAPIClient], error) {
	return &Selector[authProvider.ProviderAPIClient]{
		id:            id,
		clientFactory: authProvider.NewProviderAPIClient,
		options:       options,
	}, nil
}

// AuthApplicationSelector returns a Selector[authApplication.ApplicationsAPIClient].
func AuthApplicationSelector(id string, options ...Option) (*Selector[authApplication.ApplicationsAPIClient], error) {
	return &Selector[authApplication.ApplicationsAPIClient]{
		id:            id,
		clientFactory: authApplication.NewApplicationsAPIClient,
		options:       options,
	}, nil
}

// SharingCollaborationSelector returns a Selector[sharingCollaboration.ApplicationsAPIClient].
func SharingCollaborationSelector(id string, options ...Option) (*Selector[sharingCollaboration.CollaborationAPIClient], error) {
	return &Selector[sharingCollaboration.CollaborationAPIClient]{
		id:            id,
		clientFactory: sharingCollaboration.NewCollaborationAPIClient,
		options:       options,
	}, nil
}

// SharingOCMSelector returns a Selector[sharingOCM.OcmAPIClient].
func SharingOCMSelector(id string, options ...Option) (*Selector[sharingOCM.OcmAPIClient], error) {
	return &Selector[sharingOCM.OcmAPIClient]{
		id:            id,
		clientFactory: sharingOCM.NewOcmAPIClient,
		options:       options,
	}, nil
}

// SharingLinkSelector returns a Selector[sharingLink.LinkAPIClient].
func SharingLinkSelector(id string, options ...Option) (*Selector[sharingLink.LinkAPIClient], error) {
	return &Selector[sharingLink.LinkAPIClient]{
		id:            id,
		clientFactory: sharingLink.NewLinkAPIClient,
		options:       options,
	}, nil
}

// PreferencesSelector returns a Selector[preferences.PreferencesAPIClient].
func PreferencesSelector(id string, options ...Option) (*Selector[preferences.PreferencesAPIClient], error) {
	return &Selector[preferences.PreferencesAPIClient]{
		id:            id,
		clientFactory: preferences.NewPreferencesAPIClient,
		options:       options,
	}, nil
}

// PermissionsSelector returns a Selector[permissions.PermissionsAPIClient].
func PermissionsSelector(id string, options ...Option) (*Selector[permissions.PermissionsAPIClient], error) {
	return &Selector[permissions.PermissionsAPIClient]{
		id:            id,
		clientFactory: permissions.NewPermissionsAPIClient,
		options:       options,
	}, nil
}

// AppRegistrySelector returns a Selector[appRegistry.RegistryAPIClient].
func AppRegistrySelector(id string, options ...Option) (*Selector[appRegistry.RegistryAPIClient], error) {
	return &Selector[appRegistry.RegistryAPIClient]{
		id:            id,
		clientFactory: appRegistry.NewRegistryAPIClient,
		options:       options,
	}, nil
}

// AppProviderSelector returns a Selector[appProvider.ProviderAPIClient].
func AppProviderSelector(id string, options ...Option) (*Selector[appProvider.ProviderAPIClient], error) {
	return &Selector[appProvider.ProviderAPIClient]{
		id:            id,
		clientFactory: appProvider.NewProviderAPIClient,
		options:       options,
	}, nil
}

// StorageRegistrySelector returns a Selector[storageRegistry.RegistryAPIClient].
func StorageRegistrySelector(id string, options ...Option) (*Selector[storageRegistry.RegistryAPIClient], error) {
	return &Selector[storageRegistry.RegistryAPIClient]{
		id:            id,
		clientFactory: storageRegistry.NewRegistryAPIClient,
		options:       options,
	}, nil
}

// OCMProviderSelector returns a Selector[storageRegistry.RegistryAPIClient].
func OCMProviderSelector(id string, options ...Option) (*Selector[ocmProvider.ProviderAPIClient], error) {
	return &Selector[ocmProvider.ProviderAPIClient]{
		id:            id,
		clientFactory: ocmProvider.NewProviderAPIClient,
		options:       options,
	}, nil
}

// OCMCoreSelector returns a Selector[ocmCore.OcmCoreAPIClient].
func OCMCoreSelector(id string, options ...Option) (*Selector[ocmCore.OcmCoreAPIClient], error) {
	return &Selector[ocmCore.OcmCoreAPIClient]{
		id:            id,
		clientFactory: ocmCore.NewOcmCoreAPIClient,
		options:       options,
	}, nil
}

// TXSelector returns a Selector[tx.TxAPIClient].
func TXSelector(id string, options ...Option) (*Selector[tx.TxAPIClient], error) {
	return &Selector[tx.TxAPIClient]{
		id:            id,
		clientFactory: tx.NewTxAPIClient,
		options:       options,
	}, nil
}
