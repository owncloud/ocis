package command

import (
	"context"
	"errors"
	"fmt"

	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence"
	pcs3 "github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence/cs3"
	pfile "github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence/file"
	metadatacs3 "github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	mregistry "github.com/owncloud/ocis/v2/ocis-pkg/registry"
	sharingparser "github.com/owncloud/ocis/v2/services/sharing/pkg/config/parser"
)

// cleanCorruptPublicSharesCmd removes corrupt public-share entries (those with a
// missing/nil resource_id) from the public-share manager persistence.
//
// Background: a single public-share row with a nil resource_id makes the json
// public-share manager's ListPublicShares panic with a nil-pointer dereference,
// because it builds a cache key from resource_id without a nil-check. As the
// manager reads all rows and filters in memory, that one bad row poisons the
// endpoint for the whole tenant: every Members/permissions panel load and every
// link-with-password creation fails. This command removes the offending rows so
// affected deployments can be unblocked until the code fix is rolled out.
//
// It reads the raw persistence (it does NOT call ListPublicShares, so it never
// hits the panic) and writes back through the same metadata storage path the
// manager uses, so blob size, mtime and etag are recomputed automatically and no
// on-disk metadata surgery is required.
func cleanCorruptPublicSharesCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "clean-corrupt-public-shares",
		Usage: `Remove corrupt public-share entries (nil resource_id) that crash ListPublicShares.`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: true,
				Usage: "Only report corrupt entries, do not modify anything",
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
			return cleanCorruptPublicShares(c, cfg)
		},
	}
}

func cleanCorruptPublicShares(c *cli.Context, cfg *config.Config) error {
	dryRun := c.Bool("dry-run")
	verbose := c.Bool("verbose")

	// Initialize the service registry so the jsoncs3 driver can resolve the
	// storage-system provider via service discovery.
	_ = mregistry.GetRegistry()

	driver := cfg.Sharing.PublicSharingDriver
	p, err := publicSharePersistence(driver, cfg)
	if err != nil {
		return configlog.ReturnError(err)
	}

	ctx := context.Background()
	if err := p.Init(ctx); err != nil {
		return configlog.ReturnError(fmt.Errorf("failed to initialize public-share persistence: %w", err))
	}

	db, err := p.Read(ctx)
	if err != nil {
		return configlog.ReturnError(fmt.Errorf("failed to read public shares: %w", err))
	}

	if dryRun {
		fmt.Print("Dry run mode enabled: no entries will be removed\n\n")
	}
	fmt.Printf("Scanning %d public-share entries (driver=%s)\n", len(db), driver)

	corrupt := findCorruptPublicShares(db)
	for _, f := range corrupt {
		fmt.Printf("  CORRUPT id=%s token=%s: %s\n", f.ID, f.Token, f.Reason)
	}
	if verbose {
		for id, raw := range db {
			share, ok := decodePublicShare(raw)
			if ok && share.GetResourceId() != nil && share.GetResourceId().GetStorageId() != "" {
				fmt.Printf("  ok id=%s token=%s resource_id=%s\n", id, share.GetToken(), share.GetResourceId().GetOpaqueId())
			}
		}
	}

	if len(corrupt) == 0 {
		fmt.Println("\nNo corrupt public-share entries found. Nothing to do.")
		return nil
	}

	if dryRun {
		fmt.Printf("\nDry run: found %d corrupt entr(y/ies). Re-run with --dry-run=false to remove them.\n", len(corrupt))
		return nil
	}

	for _, f := range corrupt {
		delete(db, f.ID)
	}
	if err := p.Write(ctx, db); err != nil {
		return configlog.ReturnError(fmt.Errorf("failed to write public shares: %w", err))
	}

	fmt.Printf("\nRemoved %d corrupt public-share entr(y/ies). %d entr(y/ies) remaining.\n", len(corrupt), len(db))
	return nil
}

// corruptShare identifies a public-share persistence entry that would crash the
// json public-share manager (or that the manager cannot decode at all).
type corruptShare struct {
	ID     string
	Token  string
	Reason string
}

// findCorruptPublicShares returns the entries that must be removed: those whose
// stored value is not a decodable public-share record, and those with a nil or
// empty resource_id (the condition that makes ListPublicShares panic).
func findCorruptPublicShares(db persistence.PublicShares) []corruptShare {
	var out []corruptShare
	for id, raw := range db {
		share, ok := decodePublicShare(raw)
		if !ok {
			out = append(out, corruptShare{ID: id, Reason: "entry is not a valid public-share record"})
			continue
		}
		if share.GetResourceId() == nil || share.GetResourceId().GetStorageId() == "" {
			out = append(out, corruptShare{ID: id, Token: share.GetToken(), Reason: "nil/empty resource_id"})
		}
	}
	return out
}

// decodePublicShare extracts the cs3 PublicShare from a raw persistence entry.
// The persistence stores each entry as map{"share": <json-encoded-PublicShare>, "password": <string>}.
// It returns ok=false when the entry does not match that shape.
func decodePublicShare(raw interface{}) (*link.PublicShare, bool) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, false
	}
	s, ok := m["share"].(string)
	if !ok {
		return nil, false
	}
	var ps link.PublicShare
	if err := utils.UnmarshalJSONToProtoV1([]byte(s), &ps); err != nil {
		return nil, false
	}
	return &ps, true
}

// publicSharePersistence builds the persistence backend for the configured public-share
// driver, mirroring how reva's json public-share manager constructs it. Only the file-backed
// ("json") and cs3-backed ("jsoncs3") drivers store data in a way this command can repair.
func publicSharePersistence(driver string, cfg *config.Config) (persistence.Persistence, error) {
	switch driver {
	case "jsoncs3":
		d := cfg.Sharing.PublicSharingDrivers.JSONCS3
		s, err := metadatacs3.NewCS3Storage(d.ProviderAddr, d.ProviderAddr, d.SystemUserID, d.SystemUserIDP, d.SystemUserAPIKey)
		if err != nil {
			return nil, err
		}
		return pcs3.New(s), nil
	case "json":
		file := cfg.Sharing.PublicSharingDrivers.JSON.File
		if file == "" {
			file = "/var/tmp/reva/publicshares"
		}
		return pfile.New(file), nil
	default:
		return nil, errors.New("clean-corrupt-public-shares is only implemented for the 'jsoncs3' and 'json' public-share drivers, got: " + driver)
	}
}
