package command

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/reva/v2/pkg/bytesize"
	ocisbs "github.com/owncloud/reva/v2/pkg/storage/fs/ocis/blobstore"
	s3bs "github.com/owncloud/reva/v2/pkg/storage/fs/s3ng/blobstore"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/tree"
	"github.com/urfave/cli/v2"
)

// Blobstore is the entry point for the blobstore command group.
func Blobstore(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "blobstore",
		Usage: "manage the blobstore",
		Subcommands: []*cli.Command{
			BlobstoreCheck(cfg),
			BlobstoreGet(cfg),
		},
	}
}

// BlobstoreCheck uploads random bytes to a random path, downloads and
// verifies them, then deletes them again.  All three steps must succeed.
func BlobstoreCheck(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "check blobstore connectivity via an upload/download/delete round-trip",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "blob-size",
				Usage: "size of the random blob to upload, e.g. 64, 1KB, 1MB, 4MiB",
				Value: "64",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			bs, err := initBlobstore(cfg)
			if err != nil {
				return err
			}
			size, err := bytesize.Parse(c.String("blob-size"))
			if err != nil {
				return fmt.Errorf("invalid --blob-size %q: %w", c.String("blob-size"), err)
			}
			return runBlobstoreRoundTrip(bs, int(size))
		},
	}
}

// BlobstoreGet downloads a specific blob by its ID to verify it is readable.
func BlobstoreGet(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get a blob from the blobstore by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "blobstore path as it appears in log lines (e.g. \"<spaceID>/<pathified_blobID>\" for s3ng or \"…/spaces/<pathified_spaceID>/blobs/<pathified_blobID>\" for ocis); extracts --blob-id and --space-id automatically",
			},
			&cli.StringFlag{
				Name:  "blob-id",
				Usage: "blob ID to download (required when --path is not set)",
			},
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "space ID the blob belongs to (required when --path is not set)",
			},
			&cli.Int64Flag{
				Name: "blob-size",
				Usage: "expected blob size in bytes; only needed for the s3ng driver when the size is known upfront. " +
					"If omitted or wrong, a size mismatch will trigger one automatic retry with the actual size returned by s3ng.",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			bs, err := initBlobstore(cfg)
			if err != nil {
				return err
			}
			blobID, spaceID := c.String("blob-id"), c.String("space-id")
			if c.IsSet("path") {
				spaceID, blobID, err = parseBlobPath(c.String("path"))
				if err != nil {
					return err
				}
			}
			if blobID == "" || spaceID == "" {
				return fmt.Errorf("either --path or both --blob-id and --space-id must be set")
			}
			return downloadBlob(bs, blobID, spaceID, c.Int64("blob-size"))
		},
	}
}

// parseBlobPath extracts spaceID and blobID from a blobstore path as it
// appears in log lines.  Two formats are supported:
//
//   - S3NG:  "<spaceID>/<pathified_blobID>"
//     e.g.  "b19ec764-5398-458a-8ff1-1925bd906999/61/03/ab/c3/-b08a-4556-9937-2bf3065c1202"
//
//   - OCIS:  any path containing "/spaces/<pathified_spaceID>/blobs/<pathified_blobID>"
//     e.g.  "/var/lib/ocis/storage/users/spaces/b1/9ec764-.../blobs/61/03/ab/c3/-b08a-..."
func parseBlobPath(path string) (spaceID, blobID string, err error) {
	// OCIS paths contain both /spaces/ and /blobs/ as structural markers.
	if _, rest, ok := strings.Cut(path, "/spaces/"); ok {
		pathifiedSpaceID, pathifiedBlobID, ok := strings.Cut(rest, "/blobs/")
		if !ok {
			return "", "", fmt.Errorf("cannot parse ocis blob path %q: missing /blobs/ segment", path)
		}
		spaceID = depathify(pathifiedSpaceID, 1)
		blobID = depathify(pathifiedBlobID, 4)
		return
	}

	// S3NG paths: <spaceID>/<pathified_blobID>
	var pathifiedBlobID string
	spaceID, pathifiedBlobID, _ = strings.Cut(path, "/")
	if pathifiedBlobID == "" {
		return "", "", fmt.Errorf("cannot parse s3ng blob path %q: expected <spaceID>/<pathified_blobID>", path)
	}
	blobID = depathify(pathifiedBlobID, 4)
	return
}

// depathify reverses the effect of lookup.Pathify(id, depth, 2):
// it strips the directory separators that were inserted every two characters
// up to the given depth.
// e.g. depathify("61/03/ab/c3/-b08a-4556-9937-2bf3065c1202", 4)
//
//	→ "6103abc3-b08a-4556-9937-2bf3065c1202"
func depathify(path string, depth int) string {
	parts := strings.SplitN(path, "/", depth+1)
	return strings.Join(parts, "")
}

// initBlobstore builds a tree.Blobstore from the service configuration.
// Only the "ocis" and "s3ng" drivers are supported.
func initBlobstore(cfg *config.Config) (tree.Blobstore, error) {
	switch cfg.Driver {
	case "ocis":
		return ocisbs.New(cfg.Drivers.OCIS.Root)
	case "s3ng":
		return s3bs.New(
			cfg.Drivers.S3NG.Endpoint,
			cfg.Drivers.S3NG.Region,
			cfg.Drivers.S3NG.Bucket,
			cfg.Drivers.S3NG.AccessKey,
			cfg.Drivers.S3NG.SecretKey,
			s3bs.Options{
				DisableContentSha256:  cfg.Drivers.S3NG.DisableContentSha256,
				DisableMultipart:      cfg.Drivers.S3NG.DisableMultipart,
				SendContentMd5:        cfg.Drivers.S3NG.SendContentMd5,
				ConcurrentStreamParts: cfg.Drivers.S3NG.ConcurrentStreamParts,
				NumThreads:            cfg.Drivers.S3NG.NumThreads,
				PartSize:              cfg.Drivers.S3NG.PartSize,
			},
		)
	default:
		return nil, fmt.Errorf("blobstore operations are not supported for driver '%s'", cfg.Driver)
	}
}

// blobSizeMismatchRe matches the error produced by the s3ng blobstore when the
// retrieved object size does not match node.Blobsize, and captures the actual size.
var blobSizeMismatchRe = regexp.MustCompile(`blob has unexpected size\. \d+ bytes expected, got (\d+) bytes`)

// downloadBlob downloads a single blob identified by blobID, drains the reader
// to surface any streaming errors, then closes it.
// If the s3ng blobstore rejects the download due to a size mismatch, the actual
// size is extracted from the error and the download is retried with the correct value.
func downloadBlob(bs tree.Blobstore, blobID, spaceID string, blobSize int64) error {
	n := &node.Node{
		BlobID:   blobID,
		SpaceID:  spaceID,
		Blobsize: blobSize,
	}
	rc, err := bs.Download(n)
	if err != nil {
		if m := blobSizeMismatchRe.FindStringSubmatch(err.Error()); m != nil {
			fmt.Printf("blob size mismatch, retrying with actual size %s bytes\n", m[1])
			if _, err := fmt.Sscan(m[1], &n.Blobsize); err != nil {
				return fmt.Errorf("download failed: could not parse actual blob size %q: %w", m[1], err)
			}
			rc, err = bs.Download(n)
		}
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
	}
	defer func() {
		if cerr := rc.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close blob reader for blob %s in space %s: %v\n", blobID, spaceID, cerr)
		}
	}()
	if _, err := io.Copy(io.Discard, rc); err != nil {
		return fmt.Errorf("download failed while reading: %w", err)
	}
	fmt.Println("Download: OK")
	return nil
}

// runBlobstoreRoundTrip uploads random data, verifies it can be retrieved, and
// then removes it again.  All three steps must succeed for the check to pass.
func runBlobstoreRoundTrip(bs tree.Blobstore, blobSize int) error {
	// 1. Generate random bytes that serve as the test payload.
	data := make([]byte, blobSize)
	if _, err := rand.Read(data); err != nil {
		return fmt.Errorf("failed to generate random data: %w", err)
	}

	// Write the payload to a temporary file.  The OCIS blobstore may rename
	// (move) this file into the blobstore, so os.Remove in the defer is a
	// no-op for that driver – which is fine.
	tmpFile, err := os.CreateTemp("", "ocis-blobstore-check-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // best-effort cleanup

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// 2. Build a test node with random identifiers.
	testNode := &node.Node{
		SpaceID:  uuid.New().String(),
		BlobID:   uuid.New().String(),
		Blobsize: int64(len(data)),
	}
	fmt.Printf("Uploading test blob: spaceID=%s blobID=%s\n", testNode.SpaceID, testNode.BlobID)

	// 3. Upload.
	if err := bs.Upload(testNode, tmpPath); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	fmt.Println("Upload: OK")

	// 4. Download and verify.
	rc, err := bs.Download(testNode)
	if err != nil {
		_ = bs.Delete(testNode)
		return fmt.Errorf("download failed: %w", err)
	}
	downloaded, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		_ = bs.Delete(testNode)
		return fmt.Errorf("failed to read downloaded data: %w", err)
	}
	if !bytes.Equal(data, downloaded) {
		_ = bs.Delete(testNode)
		return fmt.Errorf("data integrity check failed: downloaded content does not match uploaded content")
	}
	fmt.Println("Download and verify: OK")

	// 5. Delete.
	if err := bs.Delete(testNode); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	fmt.Println("Delete: OK")
	fmt.Println("Blobstore check successful.")
	return nil
}
