package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		Subcommands: []*cli.Command{
			cleanupCmd(cfg),
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
