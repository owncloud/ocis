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
	"math/rand"
	"strings"
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
	ocmInvite "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	ocmProvider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	preferences "github.com/cs3org/go-cs3apis/cs3/preferences/v1beta1"
	sharingCollaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	sharingLink "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	sharingOCM "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageRegistry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	tx "github.com/cs3org/go-cs3apis/cs3/tx/v1beta1"
	"github.com/owncloud/reva/v2/pkg/registry"
	"github.com/pkg/errors"
	"github.com/sercand/kuberesolver/v5"
	mRegistry "go-micro.dev/v4/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func init() {
	// grpc go resolver.Register must only be called during initialization time (i.e. in
	// an init() function), and is not thread-safe.
	kuberesolver.RegisterInCluster()
}

type Selectable[T any] interface {
	Next(opts ...Option) (T, error)
}

var selectors sync.Map

// RemoveSelector removes given id from the selectors map.
func RemoveSelector(id string) {
	selectors.Delete(id)
}

func GetSelector[T any](k string, id string, f func(cc grpc.ClientConnInterface) T, options ...Option) *Selector[T] {
	existingSelector, ok := selectors.Load(k + id)
	if ok {
		return existingSelector.(*Selector[T])
	}

	newSelector := &Selector[T]{
		id:            id,
		clientFactory: f,
		options:       options,
	}

	selectors.Store(k+id, newSelector)

	return newSelector
}

// cachedClient pairs a grpc client with the connection backing it, so Next
// can check the connection's health before handing the client back out.
type cachedClient[T any] struct {
	conn   *grpc.ClientConn
	client T
}

// healthy reports whether the cached connection is safe to hand out again.
// A connection only becomes TransientFailure or Shutdown after a real RPC on
// it has already failed - a connection that was never used for an RPC stays
// Idle, so it is treated as healthy here (there is no way to probe
// reachability without using it).
func (c cachedClient[T]) healthy() bool {
	switch c.conn.GetState() {
	case connectivity.TransientFailure, connectivity.Connecting, connectivity.Shutdown:
		return false
	default:
		return true
	}
}

type Selector[T any] struct {
	id            string
	clientFactory func(cc grpc.ClientConnInterface) T
	clientMap     sync.Map
	options       []Option
}

func (s *Selector[T]) Next(opts ...Option) (T, error) {
	options := ClientOptions{
		registry: registry.GetRegistry(),
	}

	allOpts := append([]Option{}, s.options...)
	allOpts = append(allOpts, opts...)

	for _, opt := range allOpts {
		opt(&options)
	}

	// if the target is given as a recognized gRPC URI, skip registry lookup
	// see https://github.com/grpc/grpc/blob/master/doc/naming.md#name-syntax
	prefix := strings.SplitN(s.id, ":", 2)[0]
	switch {
	case prefix == "dns", prefix == "unix", prefix == "kubernetes":
		// use the id as is and skip registry lookup
		return s.loadOrNewClient(s.id, allOpts)
	case options.registry != nil:
		// use the service registry to look up every known node, and pick
		// the first one that is ready
		services, err := options.registry.GetService(s.id)
		if err != nil {
			return *new(T), fmt.Errorf("%s: %w", s.id, err)
		}

		targets := shuffledTargets(services)
		return s.nextHealthy(targets, allOpts)
	default:
		// if no registry is available, use the id as is
		return s.loadOrNewClient(s.id, allOpts)
	}
}

// nextHealthy walks targets in the given order. It returns the client for
// the first target that is ready. A target with no cache entry is also a
// candidate, since Next has no way to know its state without dialing it. If
// no target is ready, nextHealthy reuses the first one found, so the caller
// gets a real error from a real attempt instead of a fabricated one.
func (s *Selector[T]) nextHealthy(targets []string, opts []Option) (T, error) {
	if len(targets) == 0 {
		return *new(T), fmt.Errorf("%s: no nodes found in registry", s.id)
	}
	var fallback *cachedClient[T]
	for _, target := range targets {
		cached, ok := s.clientMap.Load(target)
		if !ok {
			return s.newClient(target, opts)
		}

		cc := cached.(cachedClient[T])
		if cc.healthy() {
			return cc.client, nil
		}
		if fallback == nil {
			fallback = &cc
		}
	}

	return fallback.client, nil
}

// loadOrNewClient returns the cached client for target, creating a new one if
// none exists yet. It does not check health: target was not looked up
// through the registry, so there is no alternative node to fall back to.
func (s *Selector[T]) loadOrNewClient(target string, opts []Option) (T, error) {
	if cached, ok := s.clientMap.Load(target); ok {
		return cached.(cachedClient[T]).client, nil
	}
	return s.newClient(target, opts)
}

// newClient creates a new connection to target, wraps it in a client and caches
// both under target.
func (s *Selector[T]) newClient(target string, opts []Option) (T, error) {
	conn, err := NewConn(target, opts...)
	if err != nil {
		return *new(T), errors.Wrap(err, fmt.Sprintf("could not create connection for %s to %s", s.id, target))
	}

	newClient := cachedClient[T]{conn: conn, client: s.clientFactory(conn)}
	// if another goroutine already created a client for this target, use it instead
	if existing, loaded := s.clientMap.LoadOrStore(target, newClient); loaded {
		_ = conn.Close()
		return existing.(cachedClient[T]).client, nil
	}

	return newClient.client, nil
}

// shuffledTargets flattens every node address out of services and returns
// them in random order, so repeated calls spread load across all known
// nodes instead of always preferring the same one.
func shuffledTargets(services []*mRegistry.Service) []string {
	var targets []string
	for _, svc := range services {
		for _, node := range svc.Nodes {
			targets = append(targets, node.Address)
		}
	}
	rand.Shuffle(len(targets), func(i, j int) {
		targets[i], targets[j] = targets[j], targets[i]
	})
	return targets
}

// GatewaySelector returns a Selector[gateway.GatewayAPIClient].
func GatewaySelector(id string, options ...Option) (*Selector[gateway.GatewayAPIClient], error) {
	return GetSelector[gateway.GatewayAPIClient](
		"GatewaySelector",
		id,
		gateway.NewGatewayAPIClient,
		options...,
	), nil
}

// IdentityUserSelector returns a Selector[identityUser.UserAPIClient].
func IdentityUserSelector(id string, options ...Option) (*Selector[identityUser.UserAPIClient], error) {
	return GetSelector[identityUser.UserAPIClient](
		"IdentityUserSelector",
		id,
		identityUser.NewUserAPIClient,
		options...,
	), nil
}

// IdentityGroupSelector returns a Selector[identityGroup.GroupAPIClient].
func IdentityGroupSelector(id string, options ...Option) (*Selector[identityGroup.GroupAPIClient], error) {
	return GetSelector[identityGroup.GroupAPIClient](
		"IdentityGroupSelector",
		id,
		identityGroup.NewGroupAPIClient,
		options...,
	), nil
}

// StorageProviderSelector returns a Selector[storageProvider.ProviderAPIClient].
func StorageProviderSelector(id string, options ...Option) (*Selector[storageProvider.ProviderAPIClient], error) {
	return GetSelector[storageProvider.ProviderAPIClient](
		"StorageProviderSelector",
		id,
		storageProvider.NewProviderAPIClient,
		options...,
	), nil
}

// SpacesProviderSelector returns a Selector[storageProvider.SpacesAPIClient].
func SpacesProviderSelector(id string, options ...Option) (*Selector[storageProvider.SpacesAPIClient], error) {
	return GetSelector[storageProvider.SpacesAPIClient](
		"SpacesProviderSelector",
		id,
		storageProvider.NewSpacesAPIClient,
		options...,
	), nil
}

// AuthRegistrySelector returns a Selector[authRegistry.RegistryAPIClient].
func AuthRegistrySelector(id string, options ...Option) (*Selector[authRegistry.RegistryAPIClient], error) {
	return GetSelector[authRegistry.RegistryAPIClient](
		"AuthRegistrySelector",
		id,
		authRegistry.NewRegistryAPIClient,
		options...,
	), nil
}

// AuthProviderSelector returns a Selector[authProvider.RegistryAPIClient].
func AuthProviderSelector(id string, options ...Option) (*Selector[authProvider.ProviderAPIClient], error) {
	return GetSelector[authProvider.ProviderAPIClient](
		"AuthProviderSelector",
		id,
		authProvider.NewProviderAPIClient,
		options...,
	), nil
}

// AuthApplicationSelector returns a Selector[authApplication.ApplicationsAPIClient].
func AuthApplicationSelector(id string, options ...Option) (*Selector[authApplication.ApplicationsAPIClient], error) {
	return GetSelector[authApplication.ApplicationsAPIClient](
		"AuthApplicationSelector",
		id,
		authApplication.NewApplicationsAPIClient,
		options...,
	), nil
}

// SharingCollaborationSelector returns a Selector[sharingCollaboration.ApplicationsAPIClient].
func SharingCollaborationSelector(id string, options ...Option) (*Selector[sharingCollaboration.CollaborationAPIClient], error) {
	return GetSelector[sharingCollaboration.CollaborationAPIClient](
		"SharingCollaborationSelector",
		id,
		sharingCollaboration.NewCollaborationAPIClient,
		options...,
	), nil
}

// SharingOCMSelector returns a Selector[sharingOCM.OcmAPIClient].
func SharingOCMSelector(id string, options ...Option) (*Selector[sharingOCM.OcmAPIClient], error) {
	return GetSelector[sharingOCM.OcmAPIClient](
		"SharingOCMSelector",
		id,
		sharingOCM.NewOcmAPIClient,
		options...,
	), nil
}

// SharingLinkSelector returns a Selector[sharingLink.LinkAPIClient].
func SharingLinkSelector(id string, options ...Option) (*Selector[sharingLink.LinkAPIClient], error) {
	return GetSelector[sharingLink.LinkAPIClient](
		"SharingLinkSelector",
		id,
		sharingLink.NewLinkAPIClient,
		options...,
	), nil
}

// PreferencesSelector returns a Selector[preferences.PreferencesAPIClient].
func PreferencesSelector(id string, options ...Option) (*Selector[preferences.PreferencesAPIClient], error) {
	return GetSelector[preferences.PreferencesAPIClient](
		"PreferencesSelector",
		id,
		preferences.NewPreferencesAPIClient,
		options...,
	), nil
}

// PermissionsSelector returns a Selector[permissions.PermissionsAPIClient].
func PermissionsSelector(id string, options ...Option) (*Selector[permissions.PermissionsAPIClient], error) {
	return GetSelector[permissions.PermissionsAPIClient](
		"PermissionsSelector",
		id,
		permissions.NewPermissionsAPIClient,
		options...,
	), nil
}

// AppRegistrySelector returns a Selector[appRegistry.RegistryAPIClient].
func AppRegistrySelector(id string, options ...Option) (*Selector[appRegistry.RegistryAPIClient], error) {
	return GetSelector[appRegistry.RegistryAPIClient](
		"AppRegistrySelector",
		id,
		appRegistry.NewRegistryAPIClient,
		options...,
	), nil
}

// AppProviderSelector returns a Selector[appProvider.ProviderAPIClient].
func AppProviderSelector(id string, options ...Option) (*Selector[appProvider.ProviderAPIClient], error) {
	return GetSelector[appProvider.ProviderAPIClient](
		"AppProviderSelector",
		id,
		appProvider.NewProviderAPIClient,
		options...,
	), nil
}

// StorageRegistrySelector returns a Selector[storageRegistry.RegistryAPIClient].
func StorageRegistrySelector(id string, options ...Option) (*Selector[storageRegistry.RegistryAPIClient], error) {
	return GetSelector[storageRegistry.RegistryAPIClient](
		"StorageRegistrySelector",
		id,
		storageRegistry.NewRegistryAPIClient,
		options...,
	), nil
}

// OCMProviderSelector returns a Selector[storageRegistry.RegistryAPIClient].
func OCMProviderSelector(id string, options ...Option) (*Selector[ocmProvider.ProviderAPIClient], error) {
	return GetSelector[ocmProvider.ProviderAPIClient](
		"OCMProviderSelector",
		id,
		ocmProvider.NewProviderAPIClient,
		options...,
	), nil
}

// OCMCoreSelector returns a Selector[ocmCore.OcmCoreAPIClient].
func OCMCoreSelector(id string, options ...Option) (*Selector[ocmCore.OcmCoreAPIClient], error) {
	return GetSelector[ocmCore.OcmCoreAPIClient](
		"OCMCoreSelector",
		id,
		ocmCore.NewOcmCoreAPIClient,
		options...,
	), nil
}

// OCMInviteSelector returns a Selector[ocmInvite.InviteAPIClient].
func OCMInviteSelector(id string, options ...Option) (*Selector[ocmInvite.InviteAPIClient], error) {
	return GetSelector[ocmInvite.InviteAPIClient](
		"OCMInviteSelector",
		id,
		ocmInvite.NewInviteAPIClient,
		options...,
	), nil
}

// TXSelector returns a Selector[tx.TxAPIClient].
func TXSelector(id string, options ...Option) (*Selector[tx.TxAPIClient], error) {
	return GetSelector[tx.TxAPIClient](
		"TXSelector",
		id,
		tx.NewTxAPIClient,
		options...,
	), nil
}
