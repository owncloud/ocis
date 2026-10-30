package svc

import "context"

var (
	CS3ReceivedShareToLibreGraphPermissions = cs3ReceivedShareToLibreGraphPermissions
	Cs3ReceivedSharesToDriveItems           = cs3ReceivedSharesToDriveItems
)

// WithGuestVaultRole exposes withGuestVaultRole to the tests
func (g Graph) WithGuestVaultRole(ctx context.Context, userID string, createShare func() error) error {
	return g.withGuestVaultRole(ctx, userID, createShare)
}
