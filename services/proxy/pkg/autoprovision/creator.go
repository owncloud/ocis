package autoprovision

import (
	"context"
	"encoding/json"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	"github.com/cs3org/reva/v2/pkg/token"
	settingsService "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
)

// Creator provides an interface to get a user or reva token with admin privileges
type Creator interface {
	// GetAutoProvisionAdmin returns a user with the Admin role assigned
	GetAutoProvisionAdmin() (*cs3.User, error)
	// GetAutoProvisionAdminToken returns a reva token with admin privileges
	GetAutoProvisionAdminToken(ctx context.Context) (string, error)
}

// Options defines the available options for this package.
type Options struct {
	tokenManager token.Manager
}

// Option defines a single option function.
type Option func(o *Options)

// WithTokenManager sets the reva token manager
func WithTokenManager(t token.Manager) Option {
	return func(o *Options) {
		o.tokenManager = t
	}
}

type creator struct {
	Options
}

// NewCreator returns a new Creator instance
func NewCreator(opts ...Option) creator {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}

	return creator{
		Options: opt,
	}
}

// This returns an hardcoded internal User, that is privileged to create new User via
// the Graph API. This user is needed for autoprovisioning of users from incoming OIDC
// claims.
func (c creator) GetAutoProvisionAdmin() (*cs3.User, error) {
	roleIDsJSON, err := json.Marshal([]string{settingsService.BundleUUIDRoleAdmin})
	if err != nil {
		return nil, err
	}

	autoProvisionUserCreator := &cs3.User{
		DisplayName: "Autoprovision User",
		Username:    "autoprovisioner",
		Id: &cs3.UserId{
			Idp:      "internal",
			OpaqueId: "autoprov-user-id00-0000-000000000000",
		},
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"roles": {
					Decoder: "json",
					Value:   roleIDsJSON,
				},
			},
		},
	}
	return autoProvisionUserCreator, nil
}

func (c creator) GetAutoProvisionAdminToken(ctx context.Context) (string, error) {
	userCreator, err := c.GetAutoProvisionAdmin()
	if err != nil {
		return "", err
	}

	s, err := scope.AddOwnerScope(nil)
	if err != nil {
		return "", err
	}

	token, err := c.tokenManager.MintToken(ctx, userCreator, s)
	if err != nil {
		return "", err
	}
	return token, nil
}
