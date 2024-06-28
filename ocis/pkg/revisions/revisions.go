// Package revisions allows manipulating revisions in a storage provider.
package revisions

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
)

var (
	// _nodesGlobPattern is the glob pattern to find all nodes
	_nodesGlobPattern = "spaces/*/*/*/*/*/*/*/*"
	// regex to determine if a node versioned. Examples:
	// 9113a718-8285-4b32-9042-f930f1a58ac2.REV.2024-05-22T07:32:53.89969726Z
	// 9113a718-8285-4b32-9042-f930f1a58ac2.REV.2024-05-22T07:32:53.89969726Z.mpk
	// 9113a718-8285-4b32-9042-f930f1a58ac2.REV.2024-05-22T07:32:53.89969726Z.mlock
	_versionRegex = regexp.MustCompile(`\.REV\.[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]+Z*`)
)

// DelBlobstore is the interface for a blobstore that can delete blobs.
type DelBlobstore interface {
	Delete(node *node.Node) error
}

// PurgeRevisions removes all revisions from a storage provider.
func PurgeRevisions(p string, bs DelBlobstore, dryRun bool, verbose bool) error {
	pattern := filepath.Join(p, _nodesGlobPattern)
	if verbose {
		fmt.Println("Looking for nodes in", pattern)
	}

	nodes, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return errors.New("no nodes found, double check storage path")
	}

	countFiles := 0
	countBlobs := 0
	countRevisions := 0
	for _, d := range nodes {
		if !_versionRegex.MatchString(d) {
			continue
		}

		var blobID string
		e := filepath.Ext(d)
		switch e {
		case ".mpk":
			blobID, err = getBlobID(d)
			if err != nil {
				fmt.Printf("error getting blobID from %s: %v\n", d, err)
				continue
			}

			countBlobs++
		case ".mlock":
			// no extra action on .mlock files
		default:
			countRevisions++
		}

		if !dryRun {
			if blobID != "" {
				//  TODO: needs spaceID for s3ng
				if err := bs.Delete(&node.Node{BlobID: blobID}); err != nil {
					fmt.Printf("error deleting blob %s: %v\n", blobID, err)
					continue
				}
			}

			if err := os.Remove(d); err != nil {
				fmt.Printf("error removing %s: %v\n", d, err)
				continue
			}

		}

		countFiles++

		if verbose {
			if dryRun {
				fmt.Println("Would delete", d)
				if blobID != "" {
					fmt.Println("Would delete blob", blobID)
				}
			} else {
				fmt.Println("Deleted", d)
				if blobID != "" {
					fmt.Println("Deleted blob", blobID)
				}
			}
		}
	}

	switch {
	case countFiles == 0 && countRevisions == 0 && countBlobs == 0:
		fmt.Println("❎ No revisions found. Storage provider is clean.")
	case !dryRun:
		fmt.Printf("✅ Deleted %d revisions (%d files / %d blobs)\n", countRevisions, countFiles, countBlobs)
	default:
		fmt.Printf("👉 Would delete %d revisions (%d files / %d blobs)\n", countRevisions, countFiles, countBlobs)
	}
	return nil
}

func getBlobID(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return "", err
	}

	if bid := m["user.ocis.blobid"]; string(bid) != "" {
		return string(bid), nil
	}

	return "", nil
}
