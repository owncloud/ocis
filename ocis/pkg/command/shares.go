package command

import (
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
	// oCIS Home
	// ├── ...
	// └── storage
	//     └── metadata
	//         ├── ...
	//         └── spaces
	//             ├── ...
	//             └── js
	//                 ├── ...
	//                 └── oncs3-share-manager-metadata
	//                     ├── blobs
	//                     │   ├── ...
	//                     │     └── d7
	//                     │         └── 02
	//                     │             └── d7
	//                     │                 └── e1
	//                     │                     └── -37b0-4d41-b8dc-4b90c1d1f907     <-- oCIS expect this blob to exist, but can be missing
	// 	                   │                        {"Shares":{"215fee7a-1d8f-404d-b563-008109b9258c:480049db-2ca5-4363-a4b3-aec71b9dab4b:84652da9-0af2-4da4-8d17-8b9b13f858f8":{"id":{"opaque_id":"215fee7a-1d8f-404d-b563-008109b9258c:480049db-2ca5-4363-a4b3-aec71b9dab4b:84652da9-0af2-4da4-8d17-8b9b13f858f8"},"resource_id":{"storage_id":"215fee7a-1d8f-404d-b563-008109b9258c","opaque_id":"49d139af-75f2-41bd-b105-0749f59dc98c","space_id":"480049db-2ca5-4363-a4b3-aec71b9dab4b"},"permissions":{"permissions":{"get_path":true,"get_quota":true,"initiate_file_download":true,"list_container":true,"list_recycle":true,"stat":true}},"grantee":{"type":1,"Id":{"UserId":{"idp":"https://localhost:9200","opaque_id":"4c510ada-c86b-4815-8820-42cdf82c3d51","type":1}}},"owner":{"opaque_id":"480049db-2ca5-4363-a4b3-aec71b9dab4b","type":8},"creator":{"idp":"https://localhost:9200","opaque_id":"c39e2f6a-a5e8-42e1-87c6-279bb570a84e","type":1},"ctime":{"seconds":1761035160,"nanos":750780584},"mtime":{"seconds":1761035160,"nanos":750780584}}},"Etag":""}
	//                     └── nodes
	//                         ├── ...
	//                         └── 99
	//                               └── 98
	//                                   └── b8
	//                                       └── bf
	//                                           ├── -6871-49cc-aca9-dab4984dc1e4
	//                                           ├── -6871-49cc-aca9-dab4984dc1e4.mlock
	//                                           └── -6871-49cc-aca9-dab4984dc1e4.mpk <-- has BlobID: d702d7e1-37b0-4d41-b8dc-4b90c1d1f907
	//                                              {
	//                                              	"user.ocis.blobsize": "988",
	//                                              	"user.ocis.id": "85c7657b-32b8-4287-9d19-fa3ba21ba6f9",
	//                                              	"user.ocis.cs.md5": "lؾˡb,s%D5",
	//                                              	"user.ocis.type": "1",
	//                                              	"user.ocis.mtime": "2025-10-21T08:26:00.771942625Z",
	//                                              	"user.ocis.blobid": "d702d7e1-37b0-4d41-b8dc-4b90c1d1f907",
	//                                              	"user.ocis.cs.sha1": "'o\u00128\u0017h%Y\u001e",
	//                                              	"user.ocis.cs.adler32": "&",
	//                                              	"user.ocis.parentid": "b73cb5e9-4842-42eb-a71e-1480969aac11",
	//                                              	"user.ocis.name": "480049db-2ca5-4363-a4b3-aec71b9dab4b.json"
	//                                              }
	//
	// 1. Walk shared nodes (mpk files)
	// 2. Extract BlobID
	// 3. Locate blob file by BlobID
	// 4. If blob file is not found, clean the node

	return &cli.Command{
		Name:  "fix-missing-share-blobs",
		Usage: `fix missing share blobs in the jsoncs3 share-manager metadata`,
		Before: func(c *cli.Context) error {
			// Parse base config to align with other shares subcommands; no config fields are required here
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			root := c.String("root")
			if root == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return configlog.ReturnError(err)
				}
				root = filepath.Join(home, ".ocis", "storage", "metadata", "spaces", "js", "oncs3-share-manager-metadata")
			}

			dryRun := true
			if c.IsSet("dry-run") {
				dryRun = c.Bool("dry-run")
			}
			if dryRun {
				fmt.Println("Dry run mode enabled")
			} else {
				fmt.Println("Dry run mode disabled")
			}

			info, err := os.Stat(root)
			if err != nil {
				return configlog.ReturnError(err)
			}
			if !info.IsDir() {
				return configlog.ReturnError(errors.New("root is not a directory"))
			}

			storageRoots, err := blobSearchRoots()
			if err != nil {
				return configlog.ReturnError(err)
			}

			nodesToClean := []string{}
			// walk mpk files, extract blob id, locate blob file
			err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				nodeToClean, rerr := visitNode(d, path, storageRoots)
				if rerr != nil {
					return rerr
				}
				if nodeToClean != "" {
					nodesToClean = append(nodesToClean, nodeToClean)
				}
				return nil
			})
			if err != nil {
				return configlog.ReturnError(err)
			}

			for _, node := range nodesToClean {
				fmt.Printf("    cleaning node=%s\n", node)
				deleteLeafFiles(node, dryRun)
				pruneEmptyDirsBelowNodes(filepath.Dir(node), dryRun)
			}

			return nil
		},
	}
}

func visitNode(d os.DirEntry, path string, storageRoots []string) (string, error) {
	if d.IsDir() {
		return "", nil
	}
	if filepath.Ext(path) != ".mpk" {
		return "", nil
	}

	b, rerr := os.ReadFile(path)
	if rerr != nil {
		return "", rerr
	}
	fmt.Printf("\n  mpk=%s\n", path)

	bid, rerr := decodeBlobID(b)
	if rerr != nil {
		return "", rerr
	}
	fmt.Printf("    blobID=%s\n", bid)
	if bid == "" {
		return "", nil
	}

	foundPath := findBlobPath(storageRoots, bid)
	if foundPath != "" {
		fmt.Printf("    exists blobPath=%s\n", foundPath)
	} else {
		fmt.Printf("    exists blobPath=(not found)\n")
		return path, nil
	}
	return "", nil
}

func blobSearchRoots() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{
		filepath.Join(home, ".ocis", "storage", "metadata", "spaces"),
		filepath.Join(home, ".ocis", "storage", "users", "spaces"),
	}, nil
}

func decodeBlobID(b []byte) (string, error) {
	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return string(m["user.ocis.blobid"]), nil
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

func findBlobPath(storageRoots []string, bid string) string {
	blobPathRelative, ok := computeBlobPathRelative(bid)
	if !ok {
		return ""
	}
	stopWalk := errors.New("stopWalk")
	for _, r := range storageRoots {
		foundPath := ""
		err := filepath.WalkDir(r, func(sp string, de os.DirEntry, e error) error {
			if e != nil || de.IsDir() {
				return nil
			}
			if strings.HasSuffix(sp, blobPathRelative) {
				foundPath = sp
				return stopWalk
			}
			return nil
		})
		if foundPath != "" {
			return foundPath
		}
		if err != nil && !errors.Is(err, stopWalk) {
			continue
		}
	}
	return ""
}

func deleteLeafFiles(node string, dryRun bool) {
	dir := filepath.Dir(node)
	base := strings.TrimSuffix(filepath.Base(node), filepath.Ext(node))
	leafs := []string{
		filepath.Join(dir, base),
		filepath.Join(dir, base+".mlock"),
		filepath.Join(dir, base+".mpk"),
	}
	for _, f := range leafs {
		if dryRun {
			fmt.Printf("      DRY RUN, SKIPPED: %s\n", f)
			continue
		}
		// if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
		// 	_ = err // return via caller if wiring deletion
		// }
	}
}

func pruneEmptyDirsBelowNodes(dir string, dryRun bool) {
	sep := string(os.PathSeparator)
	marker := sep + "nodes" + sep
	idx := strings.Index(dir, marker)
	if idx < 0 {
		return
	}
	nodesBase := dir[:idx+len("nodes")] // path ending with .../nodes
	after := dir[idx+len(marker):]
	parts := []string{}
	if after != "" {
		parts = strings.Split(after, sep)
	}
	for i := len(parts); i >= 1; i-- {
		cand := filepath.Join(nodesBase, filepath.Join(parts[:i]...))
		if dryRun {
			fmt.Printf("      DRY RUN, SKIPPED: %s\n", cand)
			continue
		}
		// entries, err := os.ReadDir(cand)
		// if err == nil && len(entries) == 0 {
		// 	_ = os.Remove(cand)
		// }
	}
}
