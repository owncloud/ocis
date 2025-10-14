package command

import (
	"errors"
	"fmt"
	"path/filepath"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/ocis/pkg/revisions"
	ocisbs "github.com/owncloud/reva/v2/pkg/storage/fs/ocis/blobstore"
	"github.com/owncloud/reva/v2/pkg/storage/fs/posix/lookup"
	s3bs "github.com/owncloud/reva/v2/pkg/storage/fs/s3ng/blobstore"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/urfave/cli/v2"
)

var (
	// _nodesGlobPattern is the glob pattern to find all nodes
	_nodesGlobPattern = "spaces/*/*/nodes/"
)

// RevisionsCommand is the entrypoint for the revisions command.
func RevisionsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "revisions",
		Usage: "ocis revisions functionality",
		Subcommands: []*cli.Command{
			PurgeRevisionsCommand(cfg),
		},
		Before: func(_ *cli.Context) error {
			return configlog.ReturnError(parser.ParseConfig(cfg, true))
		},
		Action: func(_ *cli.Context) error {
			fmt.Println("Read the docs")
			return nil
		},
	}
}

// PurgeRevisionsCommand allows removing all revisions from a storage provider.
func PurgeRevisionsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge",
		Usage: "purge revisions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "basepath",
				Aliases:  []string{"p"},
				Usage:    "the basepath of the decomposedfs (e.g. /var/tmp/ocis/storage/metadata)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "blobstore",
				Aliases: []string{"b"},
				Usage:   "the blobstore type. Can be (none, ocis, s3ng). Default ocis. Note: When using s3ng this needs same configuration as the storage-users service",
				Value:   "ocis",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "do not delete anything, just print what would be deleted",
				Value: true,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "print verbose output",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "resource-id",
				Aliases: []string{"r"},
				Usage:   "purge all revisions of this file/space. If not set, all revisions will be purged",
			},
			&cli.StringFlag{
				Name:  "glob-mechanism",
				Usage: "the glob mechanism to find all nodes. Can be 'glob', 'list' or 'workers'. 'glob' uses globbing with a single worker. 'workers' spawns multiple go routines, accelatering the command drastically but causing high cpu and ram usage. 'list' looks for references by listing directories with multiple workers. Default is 'glob'",
				Value: "glob",
			},
		},
		Action: func(c *cli.Context) error {
			basePath := c.String("basepath")
			if basePath == "" {
				fmt.Println("basepath is required")
				return cli.ShowCommandHelp(c, "revisions")
			}

			var (
				bs  revisions.DelBlobstore
				err error
			)
			switch c.String("blobstore") {
			case "s3ng":
				bs, err = s3bs.New(
					cfg.StorageUsers.Drivers.S3NG.Endpoint,
					cfg.StorageUsers.Drivers.S3NG.Region,
					cfg.StorageUsers.Drivers.S3NG.Bucket,
					cfg.StorageUsers.Drivers.S3NG.AccessKey,
					cfg.StorageUsers.Drivers.S3NG.SecretKey,
					s3bs.Options{},
				)
			case "ocis":
				bs, err = ocisbs.New(basePath)
			case "none":
				bs = nil
			default:
				err = errors.New("blobstore type not supported")
			}
			if err != nil {
				fmt.Println(err)
				return err
			}

			var rid *provider.ResourceId
			resid, err := storagespace.ParseID(c.String("resource-id"))
			if err == nil {
				rid = &resid
			}

			mechanism := c.String("glob-mechanism")
			if rid.GetOpaqueId() != "" {
				mechanism = "glob"
			}

			var ch <-chan string
			switch mechanism {
			default:
				fallthrough
			case "glob":
				p := generatePath(basePath, rid)
				if rid.GetOpaqueId() == "" {
					p = filepath.Join(p, "*/*/*/*/*")
				}
				ch = revisions.Glob(p)
			case "workers":
				p := generatePath(basePath, rid)
				ch = revisions.GlobWorkers(p, "/*", "/*/*/*/*")
			case "list":
				p := filepath.Join(basePath, "spaces")
				if rid != nil {
					p = generatePath(basePath, rid)
				}
				ch = revisions.List(p, 10)
			}

			files, blobs, revisions := revisions.PurgeRevisions(ch, bs, c.Bool("dry-run"), c.Bool("verbose"))
			printResults(files, blobs, revisions, c.Bool("dry-run"))
			return nil
		},
	}
}

func printResults(countFiles, countBlobs, countRevisions int, dryRun bool) {
	switch {
	case countFiles == 0 && countRevisions == 0 && countBlobs == 0:
		fmt.Println("❎ No revisions found. Storage provider is clean.")
	case !dryRun:
		fmt.Printf("✅ Deleted %d revisions (%d files / %d blobs)\n", countRevisions, countFiles, countBlobs)
	default:
		fmt.Printf("👉 Would delete %d revisions (%d files / %d blobs)\n", countRevisions, countFiles, countBlobs)
	}
}

func generatePath(basePath string, rid *provider.ResourceId) string {
	if rid == nil {
		return filepath.Join(basePath, _nodesGlobPattern)
	}

	sid := lookup.Pathify(rid.GetSpaceId(), 1, 2)
	if sid == "" {
		return ""
	}

	nid := lookup.Pathify(rid.GetOpaqueId(), 4, 2)
	if nid == "" {
		return filepath.Join(basePath, "spaces", sid, "nodes")
	}

	return filepath.Join(basePath, "spaces", sid, "nodes", nid+"*")
}

func init() {
	register.AddCommand(RevisionsCommand)
}
