package svc

import (
	"context"
	"math"
	"time"

	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

var (
	CS3ReceivedShareToLibreGraphPermissions = cs3ReceivedShareToLibreGraphPermissions
	Cs3ReceivedSharesToDriveItems           = cs3ReceivedSharesToDriveItems
)

// WithGuestVaultRole exposes withGuestVaultRole to the tests
func (g Graph) WithGuestVaultRole(ctx context.Context, userID string, createShare func() error) error {
	return g.withGuestVaultRole(ctx, userID, createShare)
}

// SeedNonFiniteSpecialDriveItem caches a special drive item whose photo facet holds +Inf for the given space
func (g Graph) SeedNonFiniteSpecialDriveItem(root *storageprovider.ResourceId, imageNode, readmeNode string, mtime *types.Timestamp) {
	photo := libregraph.NewPhoto()
	photo.SetExposureDenominator(math.Inf(1))
	item := libregraph.NewDriveItem()
	item.SetId("pro-1$sameID!image")
	item.SetPhoto(*photo)
	g.specialDriveItemsCache.Set(
		spaceRootStatKey(root, imageNode, readmeNode),
		specialDriveItemEntry{specialDriveItems: []libregraph.DriveItem{*item}, rootMtime: mtime},
		time.Hour,
	)
}
