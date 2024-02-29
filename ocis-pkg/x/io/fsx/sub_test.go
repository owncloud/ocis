package fsx_test

import (
	"io/fs"
	"os/exec"
	"testing"
	"testing/fstest"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/internal/testenv"
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
	cmdTest := testenv.NewCMDTest(t.Name())
	if cmdTest.ShouldRun() {
		_ = fsx.MustSub(fstest.MapFS{}, "./c")
		return
	}

	out, err := cmdTest.Run()

	g := gomega.NewWithT(t)
	g.Expect(err).To(gomega.BeAssignableToTypeOf(&exec.ExitError{}))
	g.Expect(string(out)).To(gomega.ContainSubstring("unable to load subtree fs"))
}
