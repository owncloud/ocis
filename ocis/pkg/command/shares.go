package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageregistry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/rs/zerolog"
	"github.com/shamaton/msgpack/v2"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	sdk "github.com/owncloud/reva/v2/pkg/sdk/common"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3"
	"github.com/owncloud/reva/v2/pkg/share/manager/registry"
	"github.com/owncloud/reva/v2/pkg/storage/utils/walker"
	"github.com/owncloud/reva/v2/pkg/utils"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	mregistry "github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	sharingparser "github.com/owncloud/ocis/v2/services/sharing/pkg/config/parser"
)

// SharesCommand is the entrypoint for the groups command.
func SharesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "shares",
		Usage:    `cli tools to manage entries in the share manager.`,
		Category: "maintenance",
		Subcommands: []*cli.Command{
			cleanupCmd(cfg),
			cleanOrphanedGrantsCmd(cfg),
			moveStuckUploadBlobsCmd(cfg),
		},
	}
}

func init() {
	register.AddCommand(SharesCommand)
}

func cleanupCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "cleanup",
		Usage: `Clean up stale entries in the share manager.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "service-account-id",
				Value:    "",
				Usage:    "Name of the service account to use for the cleanup",
				EnvVars:  []string{"OCIS_SERVICE_ACCOUNT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "service-account-secret",
				Value:    "",
				Usage:    "Secret for the service account",
				EnvVars:  []string{"OCIS_SERVICE_ACCOUNT_SECRET"},
				Required: true,
			},
		},
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			// Parse sharing config
			cfg.Sharing.Commons = cfg.Commons
			return configlog.ReturnError(sharingparser.ParseConfig(cfg.Sharing))
		},
		Action: func(c *cli.Context) error {
			return cleanup(c, cfg)
		},
	}
}

func cleanup(c *cli.Context, cfg *config.Config) error {
	driver := cfg.Sharing.UserSharingDriver
	// cleanup is only implemented for the jsoncs3 share manager
	if driver != "jsoncs3" {
		return configlog.ReturnError(errors.New("cleanup is only implemented for the jsoncs3 share manager"))
	}

	rcfg := revaShareConfig(cfg.Sharing)
	f, ok := registry.NewFuncs[driver]
	if !ok {
		return configlog.ReturnError(errors.New("Unknown share manager type '" + driver + "'"))
	}
	mgr, err := f(rcfg[driver].(map[string]interface{}))
	if err != nil {
		return configlog.ReturnError(err)
	}

	// Initialize registry to make service lookup work
	_ = mregistry.GetRegistry()

	// get an authenticated context
	gatewaySelector, err := pool.GatewaySelector(cfg.Sharing.Reva.Address)
	if err != nil {
		return configlog.ReturnError(err)
	}

	client, err := gatewaySelector.Next()
	if err != nil {
		return configlog.ReturnError(err)
	}

	serviceUserCtx, err := utils.GetServiceUserContext(c.String("service-account-id"), client, c.String("service-account-secret"))
	if err != nil {
		return configlog.ReturnError(err)
	}

	l := logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	serviceUserCtx = l.WithContext(serviceUserCtx)

	mgr.(*jsoncs3.Manager).CleanupStaleShares(serviceUserCtx)

	return nil
}

func cleanOrphanedGrantsCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "clean-orphaned-grants",
		Usage: `Detect and clean orphaned share-manager grants.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "service-account-id",
				Usage:    "Name of the service account to use for the scan",
				EnvVars:  []string{"OCIS_SERVICE_ACCOUNT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "service-account-secret",
				Usage:    "Secret for the service account",
				EnvVars:  []string{"OCIS_SERVICE_ACCOUNT_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "Limit the scan to a specific storage space (opaque ID)",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: true,
				Usage: "Dry run mode enabled",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force removal of suspected orphans even when listing shares fails",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Verbose logging enabled",
			},
		},
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			cfg.Sharing.Commons = cfg.Commons
			return configlog.ReturnError(sharingparser.ParseConfig(cfg.Sharing))
		},
		Action: func(c *cli.Context) error {
			return cleanOrphanedGrants(c, cfg)
		},
	}
}

func cleanOrphanedGrants(c *cli.Context, cfg *config.Config) error {
	dryRun := c.Bool("dry-run")
	verbose := c.Bool("verbose")
	spaceID := c.String("space-id")
	force := c.Bool("force")

	revaOpts := append(cfg.Sharing.Reva.GetRevaOptions(), pool.WithRegistry(mregistry.GetRegistry()))

	gatewaySelector, err := pool.GatewaySelector(cfg.Sharing.Reva.Address, revaOpts...)
	if err != nil {
		return configlog.ReturnError(err)
	}

	client, err := gatewaySelector.Next()
	if err != nil {
		return configlog.ReturnError(err)
	}

	serviceUserCtx, err := utils.GetServiceUserContext(c.String("service-account-id"), client, c.String("service-account-secret"))
	if err != nil {
		return configlog.ReturnError(err)
	}

	// Mandatory stale cleanup: reconcile jsoncs3 caches before scanning
	if cfg.Sharing.UserSharingDriver == "jsoncs3" {
		rcfg := revaShareConfig(cfg.Sharing)
		if f, ok := registry.NewFuncs["jsoncs3"]; ok {
			if mgr, mErr := f(rcfg["jsoncs3"].(map[string]interface{})); mErr == nil {
				l := logger()
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
				cleanCtx := l.WithContext(serviceUserCtx)
				mgr.(*jsoncs3.Manager).CleanupStaleShares(cleanCtx)
			}
		}
	}

	storageRegistrySelector, err := pool.StorageRegistrySelector(cfg.Sharing.Reva.Address, revaOpts...)
	if err != nil {
		return configlog.ReturnError(err)
	}

	spaces, err := resolveTargetSpaces(serviceUserCtx, gatewaySelector, spaceID)
	if err != nil {
		return configlog.ReturnError(err)
	}
	if len(spaces) == 0 {
		fmt.Println("No storage spaces found")
		if dryRun {
			fmt.Println("Dry run mode enabled: no changes executed")
		}
		return nil
	}

	walker := walker.NewWalker(gatewaySelector)
	totalNodes := 0
	totalNodesWithGrants := 0
	totalGrants := 0
	orphanCandidates := 0
	orphanDeleted := 0
	orphanDeleteErrors := 0

	if dryRun {
		fmt.Println("Dry run mode enabled")
	} else {
		fmt.Println("Dry run disabled: later stages will alter the grants")
	}

	for idx, space := range spaces {
		// Copy loop variable to avoid Go for-range capture: otherwise the callback would
		// reference the final value of 'space' for all iterations; pin this iteration's space/root.
		currentSpace := space
		rootResourceID := currentSpace.GetRoot()
		identifier := describeSpace(space)
		if space.GetRoot() == nil {
			fmt.Printf("[%d/%d] Skipping space %s: missing root resource\n", idx+1, len(spaces), identifier)
			continue
		}

		nodeCount := 0
		nodesWithGrants := 0
		grantCount := 0
		// Cache storage providers per space during this walk to avoid repeated
		// ListStorageProviders calls for every node in the same space.
		// Gateway already has a provider cache, this reduces roundtrips further.
		spaceProvidersCache := make(map[string][]*storageregistry.ProviderInfo)
		err = walker.Walk(serviceUserCtx, space.GetRoot(), func(wd string, info *provider.ResourceInfo, walkErr error) error {
			pathLabel := renderWalkerPath(wd, info)
			if walkErr != nil {
				fmt.Printf("  Warning: failed to access %s: %v\n", pathLabel, walkErr)
				return nil
			}

			if info == nil || info.GetId() == nil {
				fmt.Printf("  Warning: missing resource identifier for %s\n", pathLabel)
				return nil
			}

			nodeCount++

			// Pass A (per-node): Query storage grants. If none, do not touch share-manager.
			grants, grantErr := listResourceGrants(serviceUserCtx, storageRegistrySelector, info.GetId(), revaOpts, spaceProvidersCache)
			if grantErr != nil {
				fmt.Printf("  Warning: failed to list grants for %s: %v\n", pathLabel, grantErr)
				return nil
			}

			if len(grants) == 0 {
				if verbose {
					fmt.Printf("  %s\n", pathLabel)
				}
				return nil
			}

			// Pass B (same node, only when grants>0): Query share-manager for this resource and detect orphans.
			shares, shareErr := listSharesForResource(serviceUserCtx, gatewaySelector, info.GetId())
			if shareErr != nil && !dryRun && !force {
				fmt.Printf("  Warning: failed to list shares for %s: %v\n", pathLabel, shareErr)
				fmt.Printf("  Info: %s is too damaged to safely verify; use '--force' to remove suspected orphans regardless of share-manager errors\n", pathLabel)
			}

			nodesWithGrants++
			grantCount += len(grants)
			fmt.Printf("  %s -> %d grant(s)\n", pathLabel, len(grants))
			if verbose {
				for _, g := range grants {
					fmt.Printf("    - %s\n", describeGrant(g))
				}
			}

			// compare grantees, report suspected orphans; optionally delete when not a dry-run and shares were listed
			for _, g := range grants {
				// Skip non-user/group grantees (aligns with spec: do not touch public links)
				if gr := g.GetGrantee(); gr != nil {
					if gr.GetType() != provider.GranteeType_GRANTEE_TYPE_USER && gr.GetType() != provider.GranteeType_GRANTEE_TYPE_GROUP {
						if verbose {
							fmt.Printf("    SKIP(non-user/group): %s\n", describeGrant(g))
						}
						continue
					}
				}
				// Ignore ANY grants on the space root: built-in ownership/manager grants are not backed by share-manager
				if resourceIDEqual(info.GetId(), rootResourceID) {
					if verbose {
						fmt.Printf("    SKIP(root): %s\n", describeGrant(g))
					}
					continue
				}
				if !hasMatchingShareForGrant(shares, g) {
					orphanCandidates++
					if !dryRun && (shareErr == nil || force) {
						if err := removeGrantForResource(serviceUserCtx, storageRegistrySelector, info.GetId(), g, revaOpts); err != nil {
							fmt.Printf("    REMOVE-FAILED: %s error=%v\n", describeGrant(g), err)
							orphanDeleteErrors++
						} else {
							fmt.Printf("    REMOVED: %s\n", describeGrant(g))
							orphanDeleted++
						}
					} else {
						if shareErr != nil && !dryRun && !force {
							fmt.Printf("    ORPHAN: %s (not removed due to share list error; rerun with --force to remove)\n", describeGrant(g))
						} else {
							fmt.Printf("    ORPHAN: %s\n", describeGrant(g))
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("[%d/%d] Failed to walk space %s: %v\n", idx+1, len(spaces), identifier, err)
			continue
		}

		totalNodes += nodeCount
		totalNodesWithGrants += nodesWithGrants
		totalGrants += grantCount
		fmt.Printf("[%d/%d] Space %s: %d nodes visited, %d nodes with grants, %d grants total\n", idx+1, len(spaces), identifier, nodeCount, nodesWithGrants, grantCount)
	}

	fmt.Printf("Summary: %d space(s), %d node(s) visited, %d nodes with grants, %d grants total\n", len(spaces), totalNodes, totalNodesWithGrants, totalGrants)
	if dryRun {
		fmt.Println("Dry run mode: no grants were modified")
	} else {
		fmt.Printf("Orphans: %d candidate(s), %d removed, %d errors\n", orphanCandidates, orphanDeleted, orphanDeleteErrors)
		// Post-cleanup reconciliation: run stale cleanup again to ensure caches reflect removals
		if cfg.Sharing.UserSharingDriver == "jsoncs3" {
			rcfg := revaShareConfig(cfg.Sharing)
			if f, ok := registry.NewFuncs["jsoncs3"]; ok {
				if mgr, mErr := f(rcfg["jsoncs3"].(map[string]interface{})); mErr == nil {
					l := logger()
					zerolog.SetGlobalLevel(zerolog.InfoLevel)
					cleanCtx := l.WithContext(serviceUserCtx)
					mgr.(*jsoncs3.Manager).CleanupStaleShares(cleanCtx)
				}
			}
		}
	}

	return nil
}

// resourceIDEqual compares two provider.ResourceId values by fields.
func resourceIDEqual(a, b *provider.ResourceId) bool {
	if a == nil || b == nil {
		return false
	}
	return a.GetStorageId() == b.GetStorageId() &&
		a.GetSpaceId() == b.GetSpaceId() &&
		a.GetOpaqueId() == b.GetOpaqueId()
}

// listResourceGrants fetches grants for a resource, reusing storage provider
// resolution per space via the provided cache (keyed by space_id).
// Test locally with: simple-ocis-blobs/test_clean_orphaned_grants.sh
func listResourceGrants(ctx context.Context, registrySelector pool.Selectable[storageregistry.RegistryAPIClient], resourceID *provider.ResourceId, revaOpts []pool.Option, providerCache map[string][]*storageregistry.ProviderInfo) ([]*provider.Grant, error) {
	if resourceID == nil {
		return nil, errors.New("missing resource id")
	}

	spaceID := resourceID.GetSpaceId()
	providerInfos, ok := providerCache[spaceID]
	if !ok {
		registryClient, err := registrySelector.Next()
		if err != nil {
			return nil, err
		}

		filters := map[string]string{
			"storage_id": resourceID.GetStorageId(),
			"space_id":   spaceID,
			"opaque_id":  resourceID.GetOpaqueId(),
		}
		listReq := &storageregistry.ListStorageProvidersRequest{Opaque: &types.Opaque{}}
		sdk.EncodeOpaqueMap(listReq.Opaque, filters)

		providersResp, err := registryClient.ListStorageProviders(ctx, listReq)
		if err != nil {
			return nil, err
		}
		if providersResp.GetStatus().GetCode() != rpc.Code_CODE_OK {
			return nil, fmt.Errorf("list storage providers failed: %s", providersResp.GetStatus().GetMessage())
		}
		if len(providersResp.GetProviders()) == 0 {
			// Cache empty to avoid repeated lookups for this space in this run
			providerCache[spaceID] = nil
			return nil, nil
		}
		providerInfos = providersResp.GetProviders()
		providerCache[spaceID] = providerInfos
	}

	ref := &provider.Reference{ResourceId: resourceID}
	grants := make([]*provider.Grant, 0)
	for _, info := range providerInfos {
		if info.GetAddress() == "" {
			continue
		}
		provClient, err := pool.GetStorageProviderServiceClient(info.GetAddress(), revaOpts...)
		if err != nil {
			return nil, err
		}

		resp, err := provClient.ListGrants(ctx, &provider.ListGrantsRequest{Ref: ref})
		if err != nil {
			return nil, err
		}
		if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
			return nil, fmt.Errorf("list grants failed: %s", resp.GetStatus().GetMessage())
		}
		grants = append(grants, resp.GetGrants()...)
	}

	return grants, nil
}

func listSharesForResource(ctx context.Context, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], resourceID *provider.ResourceId) ([]*collaboration.Share, error) {
	// ListShares filtered by resource id. This reads the share-manager indexes
	// to determine if a share exists that matches a storage grant.
	client, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	req := &collaboration.ListSharesRequest{
		Filters: []*collaboration.Filter{
			{
				Type: collaboration.Filter_TYPE_RESOURCE_ID,
				Term: &collaboration.Filter_ResourceId{ResourceId: resourceID},
			},
		},
	}
	res, err := client.ListShares(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("list shares failed: %s", res.GetStatus().GetMessage())
	}
	return res.GetShares(), nil
}

func removeGrantForResource(ctx context.Context, registrySelector pool.Selectable[storageregistry.RegistryAPIClient], resourceID *provider.ResourceId, grant *provider.Grant, revaOpts []pool.Option) error {
	if resourceID == nil || grant == nil || grant.Grantee == nil {
		return errors.New("missing resource or grant")
	}
	registryClient, err := registrySelector.Next()
	if err != nil {
		return err
	}
	listReq := &storageregistry.ListStorageProvidersRequest{Opaque: &types.Opaque{}}
	sdk.EncodeOpaqueMap(listReq.Opaque, map[string]string{
		"storage_id": resourceID.GetStorageId(),
		"space_id":   resourceID.GetSpaceId(),
		"opaque_id":  resourceID.GetOpaqueId(),
	})
	providersResp, err := registryClient.ListStorageProviders(ctx, listReq)
	if err != nil {
		return err
	}
	if providersResp.GetStatus().GetCode() != rpc.Code_CODE_OK || len(providersResp.GetProviders()) == 0 {
		return fmt.Errorf("resolve provider failed: %s", providersResp.GetStatus().GetMessage())
	}
	ref := &provider.Reference{ResourceId: resourceID}
	for _, info := range providersResp.GetProviders() {
		if info.GetAddress() == "" {
			continue
		}
		provClient, err := pool.GetStorageProviderServiceClient(info.GetAddress(), revaOpts...)
		if err != nil {
			return err
		}
		_, err = provClient.RemoveGrant(ctx, &provider.RemoveGrantRequest{Ref: ref, Grant: grant})
		if err != nil {
			return err
		}
		// We consider success on first provider answering OK
		return nil
	}
	return errors.New("no provider with address")
}

func hasMatchingShareForGrant(shares []*collaboration.Share, g *provider.Grant) bool {
	// Compare grantee identity between storage grant and share-manager share.
	// ResourceId already matched by the ListShares filter.
	for _, s := range shares {
		// resource already filtered; compare grantee
		sg := s.GetGrantee()
		if sg == nil || g == nil || g.Grantee == nil {
			continue
		}
		if sg.GetType() != g.Grantee.GetType() {
			continue
		}
		switch sg.GetType() {
		case provider.GranteeType_GRANTEE_TYPE_USER:
			if sg.GetUserId() != nil && g.Grantee.GetUserId() != nil && utils.UserIDEqual(sg.GetUserId(), g.Grantee.GetUserId()) {
				return true
			}
		case provider.GranteeType_GRANTEE_TYPE_GROUP:
			if sg.GetGroupId() != nil && g.Grantee.GetGroupId() != nil && sg.GetGroupId().GetOpaqueId() == g.Grantee.GetGroupId().GetOpaqueId() {
				return true
			}
		}
	}
	return false
}

func resolveTargetSpaces(ctx context.Context, selector pool.Selectable[gateway.GatewayAPIClient], spaceID string) ([]*provider.StorageSpace, error) {
	client, err := selector.Next()
	if err != nil {
		return nil, err
	}

	req := &provider.ListStorageSpacesRequest{}
	if spaceID != "" {
		req.Filters = []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &provider.ListStorageSpacesRequest_Filter_Id{Id: &provider.StorageSpaceId{OpaqueId: spaceID}},
			},
		}
	} else {
		req.Opaque = &types.Opaque{Map: map[string]*types.OpaqueEntry{
			"unrestricted": {
				Decoder: "plain",
				Value:   []byte(strconv.FormatBool(true)),
			},
		}}
	}

	resp, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("list storage spaces failed: %s", resp.GetStatus().GetMessage())
	}
	if spaceID != "" && len(resp.GetStorageSpaces()) == 0 {
		return nil, fmt.Errorf("space not found: %s", spaceID)
	}

	return resp.GetStorageSpaces(), nil
}

func describeSpace(space *provider.StorageSpace) string {
	if space == nil {
		return "<nil>"
	}
	spaceID := ""
	if space.GetId() != nil {
		spaceID = space.GetId().GetOpaqueId()
	}
	spaceName := space.GetName()
	if spaceName == "" {
		if spaceID == "" {
			return "<unknown>"
		}
		return spaceID
	}
	if spaceID == "" {
		return spaceName
	}
	return spaceName + " (" + spaceID + ")"
}

func describeGrant(grant *provider.Grant) string {
	if grant == nil {
		return "<nil>"
	}
	grantee := describeGrantee(grant.GetGrantee())
	perms := summarizePermissions(grant.GetPermissions())
	return grantee + " permissions=[" + perms + "]"
}

func describeGrantee(grantee *provider.Grantee) string {
	if grantee == nil {
		return "grantee:<nil>"
	}
	switch grantee.GetType() {
	case provider.GranteeType_GRANTEE_TYPE_USER:
		if grantee.GetUserId() == nil {
			return "user:<unknown>"
		}
		return "user:" + grantee.GetUserId().GetOpaqueId()
	case provider.GranteeType_GRANTEE_TYPE_GROUP:
		if grantee.GetGroupId() == nil {
			return "group:<unknown>"
		}
		return "group:" + grantee.GetGroupId().GetOpaqueId()
	default:
		return strings.ToLower(grantee.GetType().String())
	}
}

func summarizePermissions(p *provider.ResourcePermissions) string {
	if p == nil {
		return "none"
	}
	perms := make([]string, 0, 8)
	if p.GetAddGrant() {
		perms = append(perms, "add-grant")
	}
	if p.GetInitiateFileDownload() {
		perms = append(perms, "download")
	}
	if p.GetInitiateFileUpload() {
		perms = append(perms, "upload")
	}
	if p.GetListGrants() {
		perms = append(perms, "list-grants")
	}
	if p.GetListContainer() {
		perms = append(perms, "list")
	}
	if p.GetStat() {
		perms = append(perms, "stat")
	}
	if p.GetRemoveGrant() {
		perms = append(perms, "remove-grant")
	}
	if p.GetUpdateGrant() {
		perms = append(perms, "update-grant")
	}
	if len(perms) == 0 {
		return "none"
	}
	return strings.Join(perms, ",")
}

func renderWalkerPath(wd string, info *provider.ResourceInfo) string {
	base := strings.TrimPrefix(filepath.ToSlash(wd), "/")
	segment := ""
	if info != nil {
		segment = strings.TrimPrefix(filepath.ToSlash(info.GetPath()), "/")
		if segment == "." {
			segment = ""
		}
	}
	combined := path.Join(base, segment)
	combined = strings.TrimPrefix(combined, "/")
	if combined == "" {
		return "/"
	}
	return "/" + combined
}

// oCIS directory structure for share-manager metadata and user spaces:
//
// ocisHome/storage/
// ...
// └── metadata/
//     ├── spaces/js/oncs3-share-manager-metadata/    (rootMetadata - Phase 1,3,4)
//     │   ├── blobs/
//     │   │   ├── 9c/a3/b2/f5/-42a1-4b8e-9123-456789abcdef   (Phase 4: MISSING received.json blob - reconstructed here)
//     │   │   │   {"Spaces": {"215fee7a-...:480049db-...": {"States": {"...:...:84652da9-...": {State: 2, MountPoint: {path: "file.txt"}}}}}}
//     │   │   └── d7/02/d7/e1/-37b0-4d41-b8dc-4b90c1d1f907   (Phase 1: read <spaceID>.json blob for Shares data)
//     │   │       {"Shares": {"215fee7a-...:480049db-...:84652da9-...": {resource_id: {...}, grantee: {...}, creator: {...}}}}
//     │   └── nodes/
//     │       ├── 3a/5f/c2/d8/-1234-5678-abcd-ef0123456789.mpk  (Phase 4: received.json MPK → points to MISSING blob)
//     │       │   {"user.ocis.name": "received.json", "user.ocis.blobid": "9ca3b2f5-42a1-4b8e-9123-456789abcdef", "user.ocis.parentid": "a9a54ce7-..."}
//     │       ├── 99/98/b8/bf/-6871-49cc-aca9-dab4984dc1e4.mpk  (Phase 1: <spaceID>.json MPK → points to Shares blob)
//     │       │   {"user.ocis.name": "480049db-...-...-....json", "user.ocis.blobid": "d702d7e1-37b0-4d41-b8dc-4b90c1d1f907"}
//     │       └── a9/a5/4c/e7/-de30-4d27-94f8-10e4612c66c2.mpk  (Phase 3: parent node for ancestry lookup)
//     │           {"user.ocis.name": "einstein", "user.ocis.id": "a9a54ce7-...", "user.ocis.parentid": "...users-node-id..."}
//     └── uploads/                                   (rootMetadataUploads)
//         ├── d702d7e1-37b0-4d41-b8dc-4b90c1d1f907   (Phase 1: read <spaceID>.json blob for Shares data; blobUploadsPath = filepath.Join(rootMetadataUploads, blobID))
//         │       {"Shares": {"215fee7a-...:480049db-...:84652da9-...": {resource_id: {...}, grantee: {...}, creator: {...}}}}
//         └── 1c93b82b-d22d-41e0-8038-5a706e9b409e.info
//                 {"MetaData": {"dir": "/users/4c510ada-c86b-4815-8820-42cdf82c3d51", "filename": "received.json", ...}, "Storage": {"NodeName": "received.json", "SpaceRoot": "jsoncs3-share-manager-metadata", ...}}

func moveStuckUploadBlobsCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "move-stuck-upload-blobs",
		Usage: `Move stuck upload blobs to the jsoncs3 share-manager metadata`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: false,
				Usage: "Dry run mode enabled",
			},
			&cli.StringFlag{
				Name:     "basepath",
				Aliases:  []string{"p"},
				Usage:    "the basepath of the decomposedfs (e.g. /var/tmp/ocis/storage/users)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "filename",
				Value: "received.json",
				Usage: "File to move from uploads/ to share manager metadata blobs/",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "Verbose logging enabled",
			},
		},
		Before: func(c *cli.Context) error {
			// Parse base config to align with other shares subcommands; no config fields are required here
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			filename := c.String("filename")
			verbose := c.Bool("verbose")

			dryRun := true
			if c.IsSet("dry-run") {
				dryRun = c.Bool("dry-run")
			}
			if dryRun {
				fmt.Print("Dry run mode enabled\n\n")
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return configlog.ReturnError(err)
			}

			ocisHome := filepath.Join(home, ".ocis")
			if c.IsSet("basepath") {
				ocisHome = c.String("basepath")
			}

			rootMetadata := filepath.Join(ocisHome, "storage", "metadata")
			rootMetadataBlobs := filepath.Join(rootMetadata, "spaces", "js", "oncs3-share-manager-metadata")

			fmt.Printf("Scanning for missing blobs in: %s \n\n", rootMetadataBlobs)
			missingBlobs, err := scanMissingBlobs(rootMetadataBlobs, filename)
			if err != nil {
				return err
			}
			if verbose {
				printJSON(missingBlobs, "missingBlobs")
			}

			if len(missingBlobs) == 0 {
				fmt.Println("No missing blobs found")
				return nil
			}

			rootMetadataUploads := filepath.Join(rootMetadata, "uploads")
			fmt.Printf("Found %d missing blobs. Restoring from %s\n", len(missingBlobs), rootMetadataUploads)
			remainingBlobIDs := restoreFromUploads(rootMetadataUploads, missingBlobs, dryRun)

			if verbose {
				printJSON(remainingBlobIDs, "remainingBlobIDs")
			}

			return nil
		},
	}
}

func printJSON(v any, label string) {
	jbs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}
	fmt.Println(label, string(jbs))
}

// Scan for missing received.json blobs
func scanMissingBlobs(rootMetadata, filename string) (map[string]string, error) {
	missingBlobs := make(map[string]string) // blobID -> blobPathAbs
	nodesRoot := filepath.Join(rootMetadata, "nodes")

	_ = filepath.WalkDir(nodesRoot, func(path string, dir os.DirEntry, err error) error {
		if err != nil || dir.IsDir() || filepath.Ext(path) != ".mpk" {
			return nil
		}
		mpkBin, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		mpk := unmarshalMPK(mpkBin)
		if mpk["user.ocis.name"] != filename {
			return nil
		}
		blobID := mpk["user.ocis.blobid"]
		blobPathRel, ok := computeBlobPathRelative(blobID)
		if !ok {
			return nil
		}
		blobPathAbs := filepath.Join(rootMetadata, blobPathRel)
		if _, statErr := os.Stat(blobPathAbs); statErr == nil {
			return nil
		}
		missingBlobs[blobID] = blobPathAbs
		return nil
	})

	return missingBlobs, nil
}

// Attempt fast path restoration from uploads/ folder
func restoreFromUploads(rootMetadataUploads string, missing map[string]string, dryRun bool) map[string]bool {
	remainingBlobIDs := make(map[string]bool)

	for blobID, blobPathAbs := range missing {
		remainingBlobIDs[blobID] = true

		blobUploadsPath := filepath.Join(rootMetadataUploads, blobID)
		if dryRun {
			fmt.Printf("    DRY RUN: move %s to %s\n", blobUploadsPath, blobPathAbs)
			continue
		}

		// Check if the blob exists in the uploads folder and move it to the share manager metadata blobs/ folder
		if _, err := os.Stat(blobUploadsPath); err != nil {
			fmt.Printf("    Blob %s: not found in %s\n", blobID, blobUploadsPath)
			continue
		}
		fmt.Printf("    Move %s to %s\n", blobUploadsPath, blobPathAbs)
		if err := os.MkdirAll(filepath.Dir(blobPathAbs), 0755); err != nil {
			fmt.Printf("    Warning: Failed to create dir: %v\n", err)
			continue
		}
		if err := os.Rename(blobUploadsPath, blobPathAbs); err != nil {
			fmt.Printf("    Warning: Failed to move blob: %v\n", err)
			continue
		}

		// Remove the info file after the blob is moved
		infoPath := blobUploadsPath + ".info"
		if _, err := os.Stat(infoPath); err != nil {
			fmt.Printf("    Info file %s: not found\n", infoPath)
			continue
		}
		if err := os.Remove(infoPath); err != nil {
			fmt.Printf("    Warning: Failed to remove info file: %v\n", err)
			continue
		}

		remainingBlobIDs[blobID] = false
	}

	return remainingBlobIDs
}

func computeBlobPathRelative(bid string) (string, bool) {
	hyphen := strings.Index(bid, "-")
	if hyphen < 0 || hyphen < 8 {
		return "", false
	}
	prefix8 := bid[:hyphen]
	if len(prefix8) < 8 {
		return "", false
	}
	d1, d2, d3, d4 := prefix8[0:2], prefix8[2:4], prefix8[4:6], prefix8[6:8]
	suffix := bid[hyphen:]
	return filepath.Join("blobs", d1, d2, d3, d4, suffix), true
}

func unmarshalMPK(bin []byte) map[string]string {
	keyValue := map[string][]byte{}
	_ = msgpack.Unmarshal(bin, &keyValue)
	out := make(map[string]string, len(keyValue))
	for k, v := range keyValue {
		out[k] = string(v)
	}
	return out
}
