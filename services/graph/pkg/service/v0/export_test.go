package svc

import "context"

var (
	CS3ReceivedShareToLibreGraphPermissions = cs3ReceivedShareToLibreGraphPermissions
	Cs3ReceivedSharesToDriveItems           = cs3ReceivedSharesToDriveItems
)

// AssignGuestVaultRole exposes assignGuestVaultRole to the tests
func (g Graph) AssignGuestVaultRole(ctx context.Context, userID string) error {
	return g.assignGuestVaultRole(ctx, userID)
}
