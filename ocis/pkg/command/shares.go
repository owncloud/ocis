package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/rs/zerolog"
	"github.com/shamaton/msgpack/v2"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/share/manager/jsoncs3"
	"github.com/owncloud/reva/v2/pkg/share/manager/registry"
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
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			// Parse sharing config
			// cfg.Sharing.Commons = cfg.Commons
			// return configlog.ReturnError(sharingparser.ParseConfig(cfg.Sharing))
			// returns error: The jwt_secret has not been set properly in your config for sharing. Make sure your /Users/mk/.ocis/config config contains the proper values (e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in the config/corresponding environment variable).
			return nil
		},
		Subcommands: []*cli.Command{
			cleanupCmd(cfg),
			missingShareBlobsCmd(cfg),
		},
	}
}

func init() {
	register.AddCommand(SharesCommand)
}

func cleanupCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "cleanup",
		Usage: `clean up stale entries in the share manager.`,
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

func missingShareBlobsCmd(cfg *config.Config) *cli.Command {
	// oCIS directory structure for share-manager metadata and user spaces:
	//
	// ocisHome/storage/
	// ...
	// ├── metadata/spaces/js/oncs3-share-manager-metadata/    (rootMetadata - Phase 1,3,4)
	// │   ├── blobs/
	// │   │   ├── d7/02/d7/e1/-37b0-4d41-b8dc-4b90c1d1f907   (Phase 1: read <spaceID>.json blob for Shares data)
	// │   │   │   {"Shares": {"215fee7a-...:480049db-...:84652da9-...": {resource_id: {...}, grantee: {...}, creator: {...}}}}
	// │   │   └── 9c/a3/b2/f5/-42a1-4b8e-9123-456789abcdef   (Phase 4: MISSING received.json blob - reconstructed here)
	// │   │       {"Spaces": {"215fee7a-...:480049db-...": {"States": {"...:...:84652da9-...": {State: 2, MountPoint: {path: "file.txt"}}}}}}
	// │   └── nodes/
	// │       ├── 99/98/b8/bf/-6871-49cc-aca9-dab4984dc1e4.mpk  (Phase 1: <spaceID>.json MPK → points to Shares blob)
	// │       │   {"user.ocis.name": "480049db-...-...-....json", "user.ocis.blobid": "d702d7e1-37b0-4d41-b8dc-4b90c1d1f907"}
	// │       ├── 3a/5f/c2/d8/-1234-5678-abcd-ef0123456789.mpk  (Phase 4: received.json MPK → points to MISSING blob)
	// │       │   {"user.ocis.name": "received.json", "user.ocis.blobid": "9ca3b2f5-42a1-4b8e-9123-456789abcdef", "user.ocis.parentid": "a9a54ce7-..."}
	// │       └── a9/a5/4c/e7/-de30-4d27-94f8-10e4612c66c2.mpk  (Phase 3: parent node for ancestry lookup)
	// │           {"user.ocis.name": "einstein", "user.ocis.id": "a9a54ce7-...", "user.ocis.parentid": "...users-node-id..."}
	// ...
	// │
	// ├──── users/spaces/ (local mode only)                     (rootUsersSpaces - Phase 2a)
	// │     └── 48/0049db-2ca5-4363-a4b3-aec71b9dab4b/nodes/
	// │         └── 49/d1/39/af/-75f2-41bd-b105-0749f59dc98c.mpk  (Phase 2a: user's file MPK with grant info)
	// │             {"user.ocis.name": "file.txt", "user.ocis.id": "49d139af-...", "user.ocis.parentid": "480049db-...",
	// │              "user.ocis.grant.u:4c510ada-...": "t^5^1:c=c39e2f6a-...:..."}
	// ...   OR
	// │
	// └──── Gateway -> Users Storage Service  (--local=false, rootUsersSpaces is in different pod - Phase 2a)
	//         For each unique MountKey{spaceID, opaqueID} from Phase 1 sharesByGrantee:
	//           gateway.Stat(ctx, &provider.StatRequest{
	//             Ref: &provider.Reference{
	//               ResourceId: &provider.ResourceId{
	//                 StorageId: "480049db-...",  // from MountKey.SpaceID
	//                 OpaqueId:  "49d139af-...",  // from MountKey.OpaqueID
	//                 SpaceId:   "480049db-...",
	//               }
	//             }
	//           })
	//           → rspStat.Info.Path = "/einstein/file.txt" → filename = filepath.Base("file.txt")
	//           → OR rspStat.Info.ArbitraryMetadata.Metadata["name"] = "file.txt"
	//           → Store: map[MountKey{spaceID, opaqueID, granteeID, creatorID}] = "file.txt"
	//         Note: Gateway routes request to storage provider pod where user space 480049db-... is mounted
	//
	//
	// Missing share blob reconstruction algorithm:
	//
	// Data structures:
	//   MountKey = {SpaceID, OpaqueID, GranteeID, CreatorID} - uniquely identifies a share mount point
	//   spaceKey = storageID:spaceID (e.g., "215fee7a-...:480049db-...")
	//   shareKey = storageID:spaceID:shareID (e.g., "215fee7a-...:480049db-...:84652da9-...")
	//
	// Phase 1: collectSharesByUser(rootMetadata) → SharesByGranteeSpaceSharekey = map[granteeID]map[spaceKey]map[shareKey]MountKey
	//   Location: rootMetadata = ocisHome/storage/metadata/spaces/js/oncs3-share-manager-metadata/
	//   Scans: rootMetadata/nodes/**/*.mpk where user.ocis.name = "<spaceID>.json"
	//   Reads blob: rootMetadata/blobs/d1/d2/d3/d4/-<suffix> (path computed from user.ocis.blobid)
	//   Blob JSON: {"Shares": {shareID: {resource_id: {storage_id, space_id, opaque_id}, grantee: {Id: {UserId: {opaque_id}}}, creator: {opaque_id}}}}
	//   Extracts:
	//     - resource_id.storage_id → storageID
	//     - resource_id.space_id → spaceID (space where shared resource lives)
	//     - resource_id.opaque_id → resourceOpaqueID (file/folder being shared)
	//     - grantee.Id.UserId.opaque_id → granteeID (user receiving share)
	//     - creator.opaque_id → creatorID (user who created share)
	//   Output: map[granteeID][spaceKey][shareKey] = MountKey{SpaceID, OpaqueID, GranteeID, CreatorID}
	//
	// Phase 2a: collectResourceNamesLocal(rootUsersSpaces) → map[MountKey]filename (local mode, --local=true)
	//   Location: rootUsersSpaces = ocisHome/storage/users/spaces/
	//   Scans: rootUsersSpaces/**/nodes/**/*.mpk (all user space node MPKs)
	//   Extracts from MPK msgpack:
	//     - user.ocis.parentid → spaceID
	//     - user.ocis.id → opaqueID (matches resource_id.opaque_id from Phase 1)
	//     - user.ocis.name → filename (e.g., "missing-blobs.txt")
	//     - Keys "user.ocis.grant.u:<granteeID>" → granteeID
	//     - Value ":c=<creatorID>:" → creatorID
	//   Output: map[MountKey{spaceID, opaqueID, granteeID, creatorID}] = filename
	//
	// Phase 2b: collectResourceNamesViaGateway(ctx, gatewayAddr, sharesByGrantee) → map[MountKey]filename (remote mode, --local=false)
	//   Gateway: gatewayAddr from --gateway-addr flag or cfg.Gateway.GRPC.Addr (default: 127.0.0.1:9142)
	//   For each unique (spaceID, opaqueID) in sharesByGrantee:
	//     - Call gateway.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{ResourceId: {StorageId: spaceID, OpaqueId: opaqueID, SpaceId: spaceID}}})
	//     - Extract filename from rspStat.Info.Path (filepath.Base) or rspStat.Info.ArbitraryMetadata.Metadata["name"]
	//     - Store: map[MountKey{spaceID, opaqueID, granteeID, creatorID}] = filename
	//   Used in production where user storage pod != share-manager metadata pod
	//
	// Phase 3: collectIdToParentId(rootMetadata) → map[nodeID]nodeMeta{ID, Name, ParentID}
	//   Location: rootMetadata/nodes/**/*.mpk
	//   Extracts: user.ocis.id → nodeID, user.ocis.name → Name, user.ocis.parentid → ParentID
	//   Ancestry chain: received.json (parentID=a9a54ce7-...) → parent node (name="einstein") → grandparent (name="users")
	//   Purpose: Fast userID resolution without re-scanning filesystem
	//
	// Phase 4: scanBlobs(idxIdToParentId, rootMetadata) → []BlobInfo
	//   Location: rootMetadata/nodes/**/*.mpk where user.ocis.name = "received.json"
	//   For each received.json MPK:
	//     - Extract user.ocis.blobid, compute blob path: rootMetadata/blobs/d1/d2/d3/d4/-<suffix>
	//     - Check blob existence via os.Stat(), skip if exists
	//     - Resolve userID: resolveUserIDForReceivedMPKFromIndex() using user.ocis.parentid → idxIdToParentId[parentID].Name
	//     - Return BlobInfo{UserID, MPKPath, BlobID, BlobRel, BlobAbs} (Payload filled in next step)
	//
	// Phase 5: buildBlobJSONForUser(userID, sharesByGrantee[userID], resourceNames) → JSON payload
	//   For each shareKey in sharesByGrantee[userID]:
	//     - Look up MountKey in resourceNames to get filename
	//     - Build: {"Spaces": {spaceKey: {"States": {shareKey: {"State": 2, "MountPoint": {"path": filename}, "Hidden": false}}}}}
	//   Write missing blobs to rootMetadata/blobs/d1/d2/d3/d4/-<suffix>
	//
	// ID associations:
	//   granteeID (from share-manager blobs) = userID (from metadata ancestry) → receiving user
	//   MountKey{spaceID, opaqueID, granteeID, creatorID} → filename → mount point path in received.json blob

	return &cli.Command{
		Name:  "fix-missing-share-blobs",
		Usage: `fix missing share blobs in the jsoncs3 share-manager metadata`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ocis-home",
				Value: "~/.ocis",
				Usage: "oCIS home directory",
			},
			&cli.BoolFlag{
				Name:  "local",
				Value: true,
				Usage: "Use local filesystem to collect resource names (false = use gateway service)",
			},
			&cli.StringFlag{
				Name:  "gateway-addr",
				Value: "127.0.0.1:9142",
				Usage: "Gateway address to use for collecting resource names (if not local)",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: false,
				Usage: "Dry run mode enabled",
			},
			&cli.BoolFlag{
				Name:  "debug-dump",
				Value: false,
				Usage: "Debug dump mode enabled",
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
			debugDump := true
			if c.IsSet("debug-dump") {
				debugDump = c.Bool("debug-dump")
			}

			dryRun := true
			if c.IsSet("dry-run") {
				dryRun = c.Bool("dry-run")
			}
			if dryRun {
				fmt.Println("Dry run mode enabled")
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return configlog.ReturnError(err)
			}

			ocisHome := filepath.Join(home, ".ocis")
			if c.IsSet("ocis-home") {
				ocisHome = c.String("ocis-home")
			}

			local := true
			if c.IsSet("local") {
				local = c.Bool("local")
			}

			rootMetadata := filepath.Join(ocisHome, "storage", "metadata", "spaces", "js", "oncs3-share-manager-metadata")
			rootUsersSpaces := filepath.Join(ocisHome, "storage", "users", "spaces")

			sharesByGrantee, err := collectSharesByUser(rootMetadata)
			if err != nil {
				return err
			}

			var resourceNames map[MountKey]string
			var uerr error
			if local {
				resourceNames = collectResourceNamesLocal(rootUsersSpaces)
			} else {
				// Get gateway address from config (cfg.Sharing.Reva.Address would require parsing sharing config)
				// For now, use environment variable or default
				gatewayAddr := cfg.Gateway.GRPC.Addr
				if gatewayAddr == "" {
					gatewayAddr = "127.0.0.1:9142" // default oCIS gateway address
				}
				resourceNames, uerr = collectResourceNamesViaGateway(context.Background(), gatewayAddr, sharesByGrantee)
				if uerr != nil {
					return uerr
				}
			}

			// Build ancestry index once for fast userID resolution without blobs
			idxIdToParentId, err := collectIdToParentId(rootMetadata)
			if err != nil {
				return err
			}

			if debugDump {
				printJSON(sharesByGrantee, "sharesByGrantee")
				printJSON(resourceNames, "resourceNames")
				printJSON(idxIdToParentId, "idxIdToParentId")
			}

			blobs, err := scanBlobs(idxIdToParentId, rootMetadata)
			if err != nil {
				return err
			}

			for i := 0; i < len(blobs); i++ {
				blobInfo := &blobs[i]
				payload, _ := buildBlobJSONForUser(blobInfo.UserID, sharesByGrantee[blobInfo.UserID], resourceNames)
				blobInfo.Payload = payload
			}

			if debugDump {
				for i, rebuild := range blobs {
					printJSON(rebuild, "Rebuild #"+strconv.Itoa(i+1))
				}
			}
			for _, blobInfo := range blobs {
				fmt.Println("    Writing blob at:", blobInfo.BlobAbs)
				fmt.Println("    Payload:", blobInfo.Payload)
				if dryRun {
					continue
				}
				err := os.WriteFile(blobInfo.BlobAbs, []byte(blobInfo.Payload), 0644)
				if err != nil {
					return configlog.ReturnError(err)
				}
			}

			return nil
		},
	}
}

func printJSON(v any, label string) {
	jbs, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(label, string(jbs))
}

type BlobInfo struct {
	UserID  string
	MPKPath string
	BlobID  string
	Payload string
	BlobRel string
	BlobAbs string
}

// MountKey uniquely identifies a target: (space_id, resource_id.opaque_id, grantee_id (user), creator_id)
type MountKey struct {
	SpaceID   string
	OpaqueID  string
	GranteeID string
	CreatorID string
}

type SharesByGranteeSpaceSharekey map[string]map[string]map[string]MountKey

// Blob rebuilding pipeline: scan received.json MPKs, detect missing blob, produce <userID, payload>
func scanBlobs(idxIdToParentId map[string]nodeMeta, rootMetadata string) ([]BlobInfo, error) {
	var blobs []BlobInfo
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
		if mpk["user.ocis.name"] != "received.json" {
			return nil
		}
		blobID := mpk["user.ocis.blobid"]
		blobPathRel, ok := computeBlobPathRelative(blobID)
		if !ok {
			return nil
		}
		blobPathAbs := filepath.Join(rootMetadata, blobPathRel)
		_, statErr := os.Stat(blobPathAbs)
		blobExists := statErr == nil
		if blobExists {
			return nil
		}

		userID, uerr := resolveUserIDForReceivedMPKFromIndex(rootMetadata, path, idxIdToParentId)
		if uerr != nil {
			return nil
		}

		blobs = append(blobs, BlobInfo{
			UserID:  userID,
			MPKPath: path,
			BlobID:  blobID,
			BlobRel: blobPathRel,
			BlobAbs: blobPathAbs,
		})
		return nil
	})
	return blobs, nil
}

// 1) From storages blobs: build shareKey → MountKey, grouped by user (grantee)
func collectSharesByUser(rootMetadata string) (SharesByGranteeSpaceSharekey, error) {
	// granteeID -> spaceKey -> shareKey -> MountKey
	idxShare := SharesByGranteeSpaceSharekey{}
	nodesRoot := filepath.Join(rootMetadata, "nodes")

	// Walk storages/*.json mpk → read their blob JSON with Shares{}
	_ = filepath.WalkDir(nodesRoot, func(path string, dir os.DirEntry, err error) error {
		if err != nil || dir.IsDir() || filepath.Ext(path) != ".mpk" {
			return nil
		}
		mpkBin, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		mpk := unmarshalMPK(mpkBin)
		name := mpk["user.ocis.name"]
		blobId := mpk["user.ocis.blobid"]
		if blobId == "" || !strings.HasSuffix(name, ".json") || name == "received.json" {
			return nil
		}
		rel, ok := computeBlobPathRelative(blobId)
		if !ok {
			return nil
		}
		blobBin, jerr := os.ReadFile(filepath.Join(rootMetadata, rel))
		if jerr != nil {
			return nil
		}
		var blobSharesModel struct {
			Shares map[string]struct {
				ResourceID struct {
					StorageID string `json:"storage_id"`
					SpaceID   string `json:"space_id"`
					OpaqueID  string `json:"opaque_id"`
				} `json:"resource_id"`
				Grantee struct {
					ID struct {
						UserID struct {
							OpaqueID string `json:"opaque_id"`
						} `json:"UserId"`
					} `json:"Id"`
				} `json:"grantee"`
				Creator struct {
					OpaqueID string `json:"opaque_id"`
				} `json:"creator"`
			} `json:"Shares"`
		}
		if json.Unmarshal(blobBin, &blobSharesModel) != nil {
			return nil
		}
		for shareID, v := range blobSharesModel.Shares {
			storageID := v.ResourceID.StorageID
			spaceID := v.ResourceID.SpaceID
			resourceOpaque := v.ResourceID.OpaqueID
			granteeID := v.Grantee.ID.UserID.OpaqueID
			creatorID := v.Creator.OpaqueID
			if storageID == "" || spaceID == "" || resourceOpaque == "" || granteeID == "" || creatorID == "" {
				continue
			}
			spaceKey := storageID + ":" + spaceID
			shareKey := shareID
			if !strings.HasPrefix(shareKey, spaceKey+":") {
				shareKey = spaceKey + ":" + shareKey
			}
			mountKey := MountKey{SpaceID: spaceID, OpaqueID: resourceOpaque, GranteeID: granteeID, CreatorID: creatorID}
			if _, ok := idxShare[granteeID]; !ok {
				idxShare[granteeID] = map[string]map[string]MountKey{}
			}
			if _, ok := idxShare[granteeID][spaceKey]; !ok {
				idxShare[granteeID][spaceKey] = map[string]MountKey{}
			}
			idxShare[granteeID][spaceKey][shareKey] = mountKey
		}
		return nil
	})
	return idxShare, nil
}

// 2) From users' MPKs: MountKey → filename (used for MountPoint.path)
func collectResourceNamesLocal(rootUsersSpaces string) map[MountKey]string {
	idxMounts := map[MountKey]string{}
	_ = filepath.WalkDir(rootUsersSpaces, func(path string, dir os.DirEntry, err error) error {
		if err != nil || dir.IsDir() || filepath.Ext(path) != ".mpk" {
			return nil
		}
		mpkBin, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		mpk := unmarshalMPK(mpkBin)
		spaceID := mpk["user.ocis.parentid"]
		opaqueID := mpk["user.ocis.id"]
		name := mpk["user.ocis.name"]
		if spaceID == "" || opaqueID == "" || name == "" {
			return nil
		}
		for key, val := range mpk {
			if !strings.HasPrefix(key, "user.ocis.grant.u:") {
				continue
			}
			granteeID := strings.TrimPrefix(key, "user.ocis.grant.u:")
			if granteeID == "" {
				continue
			}
			s := string(val)
			idx := strings.Index(s, ":c=")
			if idx < 0 {
				continue
			}
			rest := s[idx+3:]
			end := strings.IndexByte(rest, ':')
			creatorID := rest
			if end > 0 {
				creatorID = rest[:end]
			}
			if creatorID == "" {
				continue
			}
			idxMounts[MountKey{SpaceID: spaceID, OpaqueID: opaqueID, GranteeID: granteeID, CreatorID: creatorID}] = name
		}
		return nil
	})
	return idxMounts
}

// collectResourceNamesViaGateway collects resource names from user spaces via gateway service calls
// Used in production where user storage is in a different pod
func collectResourceNamesViaGateway(ctx context.Context, gatewayAddr string, sharesByGrantee SharesByGranteeSpaceSharekey) (map[MountKey]string, error) {
	idxMounts := map[MountKey]string{}

	gatewaySelector, err := pool.GatewaySelector(gatewayAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway selector: %w", err)
	}

	client, err := gatewaySelector.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway client: %w", err)
	}

	// For each share, we need to stat the resource to get its MPK attributes
	// Group by spaceID+opaqueID to avoid duplicate stats
	type resourceKey struct {
		SpaceID  string
		OpaqueID string
	}
	visited := map[resourceKey]bool{}

	for granteeID, spaceShares := range sharesByGrantee {
		for _, shares := range spaceShares {
			for _, mountKey := range shares {
				resKey := resourceKey{SpaceID: mountKey.SpaceID, OpaqueID: mountKey.OpaqueID}
				if visited[resKey] {
					continue
				}
				visited[resKey] = true

				// Stat the resource to get metadata
				ref := &provider.Reference{
					ResourceId: &provider.ResourceId{
						StorageId: mountKey.SpaceID, // In oCIS, space_id is used as storage_id for resources
						OpaqueId:  mountKey.OpaqueID,
						SpaceId:   mountKey.SpaceID,
					},
				}

				rspStat, err := client.Stat(ctx, &provider.StatRequest{Ref: ref})
				if err != nil {
					// Log error but continue processing other resources
					fmt.Printf("Warning: failed to stat resource %s:%s: %v\n", mountKey.SpaceID, mountKey.OpaqueID, err)
					continue
				}

				if rspStat.Status.Code != rpc.Code_CODE_OK {
					fmt.Printf("Warning: stat returned non-OK status for resource %s:%s: %v\n", mountKey.SpaceID, mountKey.OpaqueID, rspStat.Status.Message)
					continue
				}

				info := rspStat.Info
				if info == nil {
					continue
				}

				// Extract filename from Path or ArbitraryMetadata
				filename := filepath.Base(info.Path)
				if filename == "" || filename == "." {
					// Try to get name from arbitrary metadata
					if info.ArbitraryMetadata != nil && info.ArbitraryMetadata.Metadata != nil {
						if nameVal, ok := info.ArbitraryMetadata.Metadata["name"]; ok {
							filename = nameVal
						}
					}
				}

				if filename == "" {
					fmt.Printf("Warning: no filename found for resource %s:%s\n", mountKey.SpaceID, mountKey.OpaqueID)
					continue
				}

				// Store the filename for this MountKey
				// Note: We're storing for the specific grantee/creator combination
				idxMounts[MountKey{
					SpaceID:   mountKey.SpaceID,
					OpaqueID:  mountKey.OpaqueID,
					GranteeID: granteeID,
					CreatorID: mountKey.CreatorID,
				}] = filename
			}
		}
	}

	return idxMounts, nil
}

// 3) Resolve userID for a received.json mpk without reading its blob: parent lookup
func resolveUserIDForReceivedMPK(rootMetadata, receivedMPKPath string) (string, error) {
	mpkBin, err := os.ReadFile(receivedMPKPath)
	if err != nil {
		return "", err
	}
	mpk := unmarshalMPK(mpkBin)
	if mpk["user.ocis.name"] != "received.json" {
		return "", errors.New("not a received.json mpk")
	}
	parentID := mpk["user.ocis.parentid"]
	if parentID == "" {
		return "", errors.New("missing parent id")
	}
	// Find the parent node’s mpk by user.ocis.id == parentID to get its name (userID)
	nodesRoot := filepath.Join(rootMetadata, "nodes")
	userID := ""
	_ = filepath.WalkDir(nodesRoot, func(path string, dir os.DirEntry, err error) error {
		if userID != "" || err != nil || dir.IsDir() || filepath.Ext(path) != ".mpk" {
			return nil
		}
		mpkBin, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		mpk := unmarshalMPK(mpkBin)
		if mpk["user.ocis.id"] == parentID && mpk["user.ocis.name"] != "" {
			userID = mpk["user.ocis.name"]
		}
		return nil
	})
	if userID == "" {
		return "", errors.New("userID not found for parent")
	}
	return userID, nil
}

// Build an index of nodeID -> {name,parentID} for fast ancestry lookups
type nodeMeta struct {
	ID       string
	Name     string
	ParentID string
}

func collectIdToParentId(rootMetadata string) (map[string]nodeMeta, error) {
	idxIdToParentId := map[string]nodeMeta{}
	nodesRoot := filepath.Join(rootMetadata, "nodes")
	err := filepath.WalkDir(nodesRoot, func(path string, dir os.DirEntry, err error) error {
		if err != nil || dir.IsDir() || filepath.Ext(path) != ".mpk" {
			return nil
		}
		mpkBin, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		mpk := unmarshalMPK(mpkBin)
		if mpk["user.ocis.id"] == "" {
			return nil
		}
		idxIdToParentId[mpk["user.ocis.id"]] = nodeMeta{ID: mpk["user.ocis.id"], Name: mpk["user.ocis.name"], ParentID: mpk["user.ocis.parentid"]}
		return nil
	})
	return idxIdToParentId, err
}

// Resolve userID for a received.json mpk using ancestry: the node whose parent is "users" has Name == userID
func resolveUserIDForReceivedMPKFromIndex(rootMetadata, receivedMPKPath string, idxIdToParentId map[string]nodeMeta) (string, error) {
	mpkBin, err := os.ReadFile(receivedMPKPath)
	if err != nil {
		return "", err
	}
	mpk := unmarshalMPK(mpkBin)
	if mpk["user.ocis.name"] != "received.json" {
		return "", errors.New("not a received.json mpk")
	}
	parentID := mpk["user.ocis.parentid"]
	if parentID == "" {
		return "", errors.New("missing parent id")
	}
	// Fast path: the immediate parent directory of received.json is the userID directory
	if node, ok := idxIdToParentId[parentID]; ok && node.Name != "" {
		return node.Name, nil
	}
	// Direct parent mpk lookup if not in index
	if mpkPathRel, ok := computeNodeMPKPathRelative(parentID); ok {
		parentMPK := filepath.Join(rootMetadata, mpkPathRel)
		if mpkBinParent, err := os.ReadFile(parentMPK); err == nil {
			mpkParent := unmarshalMPK(mpkBinParent)
			if parentName := mpkParent["user.ocis.name"]; parentName != "" {
				return parentName, nil
			}
		}
	}
	// On-demand scan fallback if index doesn't contain the parent node
	if userID, scanErr := resolveUserIDForReceivedMPK(rootMetadata, receivedMPKPath); scanErr == nil && userID != "" {
		fmt.Println("resolveUserIDForReceivedMPK fallback: uid", userID)
		return userID, nil
	}
	// Fallback: Walk up until we find a node named "users"; the child just below it is the userID
	lastParentName := ""
	curr := parentID
	for i := 0; i < 1024 && curr != ""; i++ { // safety bound
		node, ok := idxIdToParentId[curr]
		if !ok {
			break
		}
		if node.Name == "users" {
			if lastParentName == "" {
				return "", errors.New("users ancestor found but child name empty")
			}
			return lastParentName, nil
		}
		lastParentName = node.Name
		curr = node.ParentID
	}
	if lastParentName != "" { // best-effort
		return lastParentName, nil
	}
	return "", errors.New("userID not found via ancestry")
}

// 4) For each received.json mpk whose blob is missing, build payload for its user (userID == MountKey.GranteeID)
func buildBlobJSONForUser(granteeID string, sharesByUser map[string]map[string]MountKey, resourceNames map[MountKey]string) (string, error) {
	spaces := map[string]any{}
	for spaceKey, shares := range sharesByUser {
		states := map[string]any{}
		for shareKey, mountKey := range shares {
			// Ensure we only emit shares for this grantee
			if mountKey.GranteeID != granteeID {
				continue
			}
			mountPath := resourceNames[mountKey] // empty if not found; still valid JSON
			states[shareKey] = map[string]any{
				"State":      2,
				"MountPoint": map[string]string{"path": mountPath},
				"Hidden":     false,
			}
		}
		if len(states) > 0 {
			spaces[spaceKey] = map[string]any{"States": states}
		}
	}
	blob := map[string]any{"Spaces": spaces}
	blobString, _ := json.MarshalIndent(blob, "", "  ")
	return string(blobString), nil
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

// computeNodeMPKPathRelative converts a node UUID to nodes/<d1>/<d2>/<d3>/<d4>/-<suffix>.mpk
func computeNodeMPKPathRelative(nodeID string) (string, bool) {
	hyphen := strings.Index(nodeID, "-")
	if hyphen < 0 || hyphen < 8 {
		return "", false
	}
	prefix8 := nodeID[:hyphen]
	if len(prefix8) < 8 {
		return "", false
	}
	d1, d2, d3, d4 := prefix8[0:2], prefix8[2:4], prefix8[4:6], prefix8[6:8]
	suffix := nodeID[hyphen:]
	return filepath.Join("nodes", d1, d2, d3, d4, suffix+".mpk"), true
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
