package pool

import (
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
	"google.golang.org/grpc"
	"sync"
)

type Selectable[T any] interface {
	Next(opts ...Option) (T, error)
}

type Selector[T any] struct {
	id string
	cf func(cc *grpc.ClientConn) T
	cm sync.Map
}

func (s *Selector[T]) Next(opts ...Option) (T, error) {
	services, _ := registry.DiscoverServices(s.id)
	address, err := registry.GetNodeAddress(services)
	if err != nil || address == "" {
		address = s.id
	}

	existingClient, ok := s.cm.Load(address)
	if ok {
		return existingClient.(T), nil
	}

	conn, err := NewConn(s.id, opts...)
	if err != nil {
		return *new(T), err
	}

	newClient := s.cf(conn)
	s.cm.Store(address, newClient)

	return newClient, nil
}

// GatewaySelector returns a Selector[gateway.GatewayAPIClient].
func GatewaySelector(id string) (Selector[gateway.GatewayAPIClient], error) {
	return Selector[gateway.GatewayAPIClient]{
		id: id,
		cf: gateway.NewGatewayAPIClient,
	}, nil
}

// IdentityUserSelector returns a Selector[identityUser.UserAPIClient].
func IdentityUserSelector(id string) (Selector[identityUser.UserAPIClient], error) {
	return Selector[identityUser.UserAPIClient]{
		id: id,
		cf: identityUser.NewUserAPIClient,
	}, nil
}

// IdentityGroupSelector returns a Selector[identityGroup.GroupAPIClient].
func IdentityGroupSelector(id string) (Selector[identityGroup.GroupAPIClient], error) {
	return Selector[identityGroup.GroupAPIClient]{
		id: id,
		cf: identityGroup.NewGroupAPIClient,
	}, nil
}

// StorageProviderSelector returns a Selector[storageProvider.ProviderAPIClient].
func StorageProviderSelector(id string) (Selector[storageProvider.ProviderAPIClient], error) {
	return Selector[storageProvider.ProviderAPIClient]{
		id: id,
		cf: storageProvider.NewProviderAPIClient,
	}, nil
}

// AuthRegistrySelector returns a Selector[authRegistry.RegistryAPIClient].
func AuthRegistrySelector(id string) (Selector[authRegistry.RegistryAPIClient], error) {
	return Selector[authRegistry.RegistryAPIClient]{
		id: id,
		cf: authRegistry.NewRegistryAPIClient,
	}, nil
}

// AuthProviderSelector returns a Selector[authProvider.RegistryAPIClient].
func AuthProviderSelector(id string) (Selector[authProvider.ProviderAPIClient], error) {
	return Selector[authProvider.ProviderAPIClient]{
		id: id,
		cf: authProvider.NewProviderAPIClient,
	}, nil
}

// AuthApplicationSelector returns a Selector[authApplication.ApplicationsAPIClient].
func AuthApplicationSelector(id string) (Selector[authApplication.ApplicationsAPIClient], error) {
	return Selector[authApplication.ApplicationsAPIClient]{
		id: id,
		cf: authApplication.NewApplicationsAPIClient,
	}, nil
}

// SharingCollaborationSelector returns a Selector[sharingCollaboration.ApplicationsAPIClient].
func SharingCollaborationSelector(id string) (Selector[sharingCollaboration.CollaborationAPIClient], error) {
	return Selector[sharingCollaboration.CollaborationAPIClient]{
		id: id,
		cf: sharingCollaboration.NewCollaborationAPIClient,
	}, nil
}

// SharingOCMSelector returns a Selector[sharingOCM.OcmAPIClient].
func SharingOCMSelector(id string) (Selector[sharingOCM.OcmAPIClient], error) {
	return Selector[sharingOCM.OcmAPIClient]{
		id: id,
		cf: sharingOCM.NewOcmAPIClient,
	}, nil
}

// SharingLinkSelector returns a Selector[sharingLink.LinkAPIClient].
func SharingLinkSelector(id string) (Selector[sharingLink.LinkAPIClient], error) {
	return Selector[sharingLink.LinkAPIClient]{
		id: id,
		cf: sharingLink.NewLinkAPIClient,
	}, nil
}

// PreferencesSelector returns a Selector[preferences.PreferencesAPIClient].
func PreferencesSelector(id string) (Selector[preferences.PreferencesAPIClient], error) {
	return Selector[preferences.PreferencesAPIClient]{
		id: id,
		cf: preferences.NewPreferencesAPIClient,
	}, nil
}

// PermissionsSelector returns a Selector[permissions.PermissionsAPIClient].
func PermissionsSelector(id string) (Selector[permissions.PermissionsAPIClient], error) {
	return Selector[permissions.PermissionsAPIClient]{
		id: id,
		cf: permissions.NewPermissionsAPIClient,
	}, nil
}

// AppRegistrySelector returns a Selector[appRegistry.RegistryAPIClient].
func AppRegistrySelector(id string) (Selector[appRegistry.RegistryAPIClient], error) {
	return Selector[appRegistry.RegistryAPIClient]{
		id: id,
		cf: appRegistry.NewRegistryAPIClient,
	}, nil
}

// AppProviderSelector returns a Selector[appProvider.ProviderAPIClient].
func AppProviderSelector(id string) (Selector[appProvider.ProviderAPIClient], error) {
	return Selector[appProvider.ProviderAPIClient]{
		id: id,
		cf: appProvider.NewProviderAPIClient,
	}, nil
}

// StorageRegistrySelector returns a Selector[storageRegistry.RegistryAPIClient].
func StorageRegistrySelector(id string) (Selector[storageRegistry.RegistryAPIClient], error) {
	return Selector[storageRegistry.RegistryAPIClient]{
		id: id,
		cf: storageRegistry.NewRegistryAPIClient,
	}, nil
}

// OCMProviderSelector returns a Selector[storageRegistry.RegistryAPIClient].
func OCMProviderSelector(id string) (Selector[ocmProvider.ProviderAPIClient], error) {
	return Selector[ocmProvider.ProviderAPIClient]{
		id: id,
		cf: ocmProvider.NewProviderAPIClient,
	}, nil
}

// OCMCoreSelector returns a Selector[ocmCore.OcmCoreAPIClient].
func OCMCoreSelector(id string) (Selector[ocmCore.OcmCoreAPIClient], error) {
	return Selector[ocmCore.OcmCoreAPIClient]{
		id: id,
		cf: ocmCore.NewOcmCoreAPIClient,
	}, nil
}

// TXSelector returns a Selector[tx.TxAPIClient].
func TXSelector(id string) (Selector[tx.TxAPIClient], error) {
	return Selector[tx.TxAPIClient]{
		id: id,
		cf: tx.NewTxAPIClient,
	}, nil
}
