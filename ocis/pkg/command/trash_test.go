package command_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis/pkg/command"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// orphanedLockFile creates a lock file without a node and returns its path
func orphanedLockFile(t *testing.T, base string) string {
	t.Helper()
	path := filepath.Join(base, "spaces", "ab", "cdef", "nodes", "11", "22", "33", "44", "-node.mlock")
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0700))
	require.NoError(t, os.WriteFile(path, nil, 0600))
	old := time.Now().Add(-48 * time.Hour)
	require.NoError(t, os.Chtimes(path, old, old))
	return path
}

func runPurgeOrphanedFiles(t *testing.T, args ...string) error {
	t.Helper()
	app := &cli.App{
		Commands: []*cli.Command{command.TrashPurgeOrphanedFilesCommand(config.DefaultConfig())},
	}
	return app.Run(append([]string{"ocis", "purge-orphaned-files"}, args...))
}

func TestTrashPurgeOrphanedFilesIsDryRunByDefault(t *testing.T) {
	base := t.TempDir()
	lockFile := orphanedLockFile(t, base)

	require.NoError(t, runPurgeOrphanedFiles(t, "--basepath", base))

	require.FileExists(t, lockFile)
}

func TestTrashPurgeOrphanedFilesRemovesOrphanedFiles(t *testing.T) {
	base := t.TempDir()
	lockFile := orphanedLockFile(t, base)

	require.NoError(t, runPurgeOrphanedFiles(t, "--basepath", base, "--dry-run=false"))

	require.NoFileExists(t, lockFile)
}

func TestTrashPurgeOrphanedFilesKeepsFilesYoungerThanMinAge(t *testing.T) {
	base := t.TempDir()
	lockFile := orphanedLockFile(t, base)

	require.NoError(t, runPurgeOrphanedFiles(t, "--basepath", base, "--dry-run=false", "--min-age", "72h"))

	require.FileExists(t, lockFile)
}

func TestTrashPurgeOrphanedFilesRejectsUnknownBlobstore(t *testing.T) {
	base := t.TempDir()
	orphanedLockFile(t, base)

	require.Error(t, runPurgeOrphanedFiles(t, "--basepath", base, "--blobstore", "unknown"))
}
