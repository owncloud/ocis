package trash

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
)

const (
	// _nodesGlobPattern is the glob pattern to find all directories with node files
	_nodesGlobPattern = "spaces/*/*/nodes/*/*/*/*"

	_lockSuffix     = ".mlock"
	_metadataSuffix = ".mpk"
	_blobIDAttr     = "user.ocis.blobid"
)

// _revisionRegex matches the name of a revision. The first group is the name of the node.
var _revisionRegex = regexp.MustCompile(`^(.+)\.REV\.[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?Z$`)

// BlobDeleter deletes blobs from the blobstore
type BlobDeleter interface {
	Delete(node *node.Node) error
}

// PurgeOrphansOptions are the options for PurgeOrphanedFiles
type PurgeOrphansOptions struct {
	// Blobstore deletes the blobs of orphaned revisions. If it is nil, orphaned revisions stay on the disk.
	Blobstore BlobDeleter
	// MinAge is the minimum age of a file. Younger files stay on the disk.
	MinAge time.Duration
	// DryRun prints the files, but does not remove them.
	DryRun bool
	// Out receives the report.
	Out io.Writer
}

type orphanPurger struct {
	opts PurgeOrphansOptions

	lockFiles int
	revisions int
	skipped   int
}

// PurgeOrphanedFiles removes the files of nodes that do not exist:
//   - a lock file (.mlock) when its node does not exist
//   - a revision with its metadata file, its lock file and its blob, when its node does not exist and is not in the trash bin
func PurgeOrphanedFiles(basePath string, opts PurgeOrphansOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	dirs, err := filepath.Glob(filepath.Join(basePath, _nodesGlobPattern))
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return errors.New("no nodes found. Double check storage path")
	}

	p := &orphanPurger{opts: opts}
	var errs []error
	for _, dir := range dirs {
		if err := p.purgeDir(basePath, dir); err != nil {
			fmt.Fprintln(opts.Out, err)
			errs = append(errs, err)
		}
	}

	verb := "removed"
	if opts.DryRun {
		verb = "would remove"
	}
	fmt.Fprintf(opts.Out, "%s %d lock files and %d revisions. %d files are younger than %s and stay.\n", verb, p.lockFiles, p.revisions, p.skipped, opts.MinAge)
	if p.opts.Blobstore == nil {
		fmt.Fprintln(opts.Out, "no blobstore is set: orphaned revisions stay.")
	}
	return errors.Join(errs...)
}

func (p *orphanPurger) purgeDir(basePath, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	revisions := map[string]struct{}{}
	var lockFiles []string
	for _, e := range entries {
		name := e.Name()
		base := strings.TrimSuffix(strings.TrimSuffix(name, _lockSuffix), _metadataSuffix)
		if _revisionRegex.MatchString(base) {
			revisions[base] = struct{}{}
			continue
		}
		if strings.HasSuffix(name, _lockSuffix) {
			lockFiles = append(lockFiles, name)
		}
	}

	var errs []error
	for _, name := range lockFiles {
		if err := p.purgeLockFile(dir, name); err != nil {
			errs = append(errs, err)
		}
	}

	names := make([]string, 0, len(revisions))
	for name := range revisions {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := p.purgeRevision(basePath, dir, name); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// purgeLockFile removes a lock file when the file or directory of its node does not exist
func (p *orphanPurger) purgeLockFile(dir, name string) error {
	path := filepath.Join(dir, name)
	if exists(strings.TrimSuffix(path, _lockSuffix)) {
		return nil
	}
	if young, err := p.isYoung(path); err != nil || young {
		return err
	}

	p.lockFiles++
	return p.remove(path)
}

// purgeRevision removes a revision when its node does not exist and is not in the trash bin
func (p *orphanPurger) purgeRevision(basePath, dir, name string) error {
	nodePath := filepath.Join(dir, _revisionRegex.FindStringSubmatch(name)[1])
	if exists(nodePath) {
		return nil
	}
	trashed, err := filepath.Glob(nodePath + node.TrashIDDelimiter + "*")
	if err != nil {
		return err
	}
	for _, t := range trashed {
		if !strings.HasSuffix(t, _lockSuffix) && !strings.HasSuffix(t, _metadataSuffix) {
			return nil
		}
	}

	path := filepath.Join(dir, name)
	files := []string{path, path + _metadataSuffix, path + _lockSuffix}
	for _, f := range files {
		if young, err := p.isYoung(f); err != nil || young {
			return err
		}
	}

	if p.opts.Blobstore == nil {
		fmt.Fprintln(p.opts.Out, "skip revision without blobstore", path)
		return nil
	}

	if err := p.deleteBlob(basePath, path); err != nil {
		// keep the metadata file. It has the only reference to the blob.
		return fmt.Errorf("could not delete the blob of revision %s: %w", path, err)
	}

	p.revisions++
	var errs []error
	for _, f := range files {
		if err := p.remove(f); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (p *orphanPurger) deleteBlob(basePath, revisionPath string) error {
	b, err := os.ReadFile(revisionPath + _metadataSuffix)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	attrs := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &attrs); err != nil {
		return err
	}
	blobID := string(attrs[_blobIDAttr])
	if blobID == "" {
		return nil
	}

	spaceID := spaceIDFromPath(basePath, revisionPath)
	if p.opts.DryRun {
		fmt.Fprintf(p.opts.Out, "would delete blob %s of space %s\n", blobID, spaceID)
		return nil
	}
	if err := p.opts.Blobstore.Delete(&node.Node{SpaceID: spaceID, BlobID: blobID}); err != nil {
		return err
	}
	fmt.Fprintf(p.opts.Out, "deleted blob %s of space %s\n", blobID, spaceID)
	return nil
}

// isYoung returns true when the file is younger than MinAge. A missing file is not young.
func (p *orphanPurger) isYoung(path string) (bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if time.Since(info.ModTime()) < p.opts.MinAge {
		p.skipped++
		return true, nil
	}
	return false, nil
}

func (p *orphanPurger) remove(path string) error {
	if p.opts.DryRun {
		if exists(path) {
			fmt.Fprintln(p.opts.Out, "would remove", path)
		}
		return nil
	}
	err := os.Remove(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(p.opts.Out, "removed", path)
	return nil
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// spaceIDFromPath returns the space ID from a path like <basePath>/spaces/ab/cdef/nodes/...
func spaceIDFromPath(basePath, path string) string {
	rel, err := filepath.Rel(filepath.Join(basePath, "spaces"), path)
	if err != nil {
		return ""
	}
	sep := string(filepath.Separator)
	first, rest, _ := strings.Cut(rel, sep)
	second, _, found := strings.Cut(rest, sep)
	if !found {
		return ""
	}
	return first + second
}
