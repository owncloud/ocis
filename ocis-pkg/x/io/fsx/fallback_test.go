package fsx_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

func testFallbackFSENV() (*fsx.FallbackFS, map[string]string, func(), error) {
	nf := func(n, c string) (*os.File, error) {
		f, err := os.CreateTemp("", n)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		_, err = f.Write([]byte(c))
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	foo, err := nf("foo.txt", "foo - fs")
	if err != nil {
		return nil, nil, nil, err
	}

	bar, err := nf("bar.txt", "bar - fs")
	if err != nil {
		return nil, nil, nil, err
	}

	m := map[string]string{
		"foo.txt": filepath.Base(foo.Name()),
		"bar.txt": filepath.Base(bar.Name()),
		"baz.txt": "baz.txt",
	}

	fsys := fstest.MapFS{
		m["foo.txt"]: &fstest.MapFile{
			Data: []byte("foo - embedded"),
		},
		m["bar.txt"]: &fstest.MapFile{
			Data: []byte("bar - embedded"),
		},
		m["baz.txt"]: &fstest.MapFile{
			Data: []byte("baz - embedded"),
		},
	}

	fbsys := fsx.NewFallbackFS(fsys, os.TempDir())

	return fbsys, m, func() {
		_ = os.Remove(foo.Name())
		_ = os.Remove(bar.Name())
	}, nil
}

func testFallbackFSContent(t *testing.T, f fs.File, c string) {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(f); err != nil {
		t.Fatalf("expected to read from file, but got %v", err)
	}

	if buf.String() != c {
		t.Fatalf("expected to read from fs, but got %s", buf.String())
	}
}

func TestFallbackFS_Open(t *testing.T) {
	fsys, m, cleanup, err := testFallbackFSENV()
	if err != nil {
		t.Fatalf("expected to create test environment, but got %v", err)
	}
	defer cleanup()

	foo, _ := fsys.Open(m["foo.txt"])
	testFallbackFSContent(t, foo, "foo - fs")

	bar, _ := fsys.Open(m["bar.txt"])
	testFallbackFSContent(t, bar, "bar - fs")

	baz, _ := fsys.Open(m["baz.txt"])
	testFallbackFSContent(t, baz, "baz - embedded")
}

func TestFallbackFS_OpenEmbedded(t *testing.T) {
	fsys, m, cleanup, err := testFallbackFSENV()
	if err != nil {
		t.Fatalf("expected to create test environment, but got %v", err)
	}
	defer cleanup()

	foo, _ := fsys.OpenEmbedded(m["foo.txt"])
	testFallbackFSContent(t, foo, "foo - embedded")

	bar, _ := fsys.OpenEmbedded(m["bar.txt"])
	testFallbackFSContent(t, bar, "bar - embedded")

	baz, _ := fsys.OpenEmbedded(m["baz.txt"])
	testFallbackFSContent(t, baz, "baz - embedded")
}

func TestFallbackFS_Create(t *testing.T) {
	fsys, _, cleanup, err := testFallbackFSENV()
	if err != nil {
		t.Fatalf("expected to create test environment, but got %v", err)
	}
	defer cleanup()

	baz, _ := fsys.Open("baz.txt")
	testFallbackFSContent(t, baz, "baz - embedded")

	// parent directory access is here to test the jail join
	f, err := fsys.Create("../././././././baz.txt")
	if err != nil {
		t.Fatalf("expected to create file, but got %v", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.Write([]byte("baz - fs"))
	if err != nil {
		t.Fatalf("expected to write to file, but got %v", err)
	}

	baz, _ = fsys.Open("baz.txt")
	testFallbackFSContent(t, baz, "baz - fs")
}
