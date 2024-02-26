package fsx_test

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

func TestMustSub(t *testing.T) {
	tfs := fstest.MapFS{
		"a/foo.txt": &fstest.MapFile{},
		"b/bar.txt": &fstest.MapFile{},
		"b/baz": &fstest.MapFile{
			Mode: fs.ModeDir,
		},
	}

	sfs := fsx.MustSub(tfs, "b")
	entries, err := fs.ReadDir(sfs, ".")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestMustSub_Exit(t *testing.T) {
	if os.Getenv("WITH_EXIT") == "1" {
		_ = fsx.MustSub(fstest.MapFS{}, "./c")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMustSub_Exit")
	cmd.Env = append(os.Environ(), "WITH_EXIT=1")
	err := cmd.Run()

	var e *exec.ExitError
	assert.Equal(t, true, errors.As(err, &e))
	assert.Equal(t, "exit status 1", e.Error())
}
