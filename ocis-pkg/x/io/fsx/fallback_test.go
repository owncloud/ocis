package fsx_test

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

func createFallbackFS(d string) (*fsx.FallbackFS, map[string]string, error) {
	if err := os.Mkdir(d, 0700); err != nil {
		return nil, nil, err
	}

	foo, err := createFile(d, "foo.txt", "foo - fs")
	if err != nil {
		return nil, nil, err
	}

	bar, err := createFile(d, "bar.txt", "bar - fs")
	if err != nil {
		return nil, nil, err
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

	return fsx.NewFallbackFS(fsys, d), m, nil
}

func readFrom(r io.Reader) (string, error) {
	buf := bytes.NewBuffer(nil)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func createFile(d, n, c string) (*os.File, error) {
	f, err := os.CreateTemp(d, n)

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

func TestFallbackFS_Open(t *testing.T) {
	d := filepath.Join(t.TempDir(), "allowed")
	g := gomega.NewWithT(t)

	fsys, m, err := createFallbackFS(d)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	foo, err := fsys.Open(m["foo.txt"])
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer foo.Close()
	content, err := readFrom(foo)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("foo - fs"))

	bar, err := fsys.Open(m["bar.txt"])
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer bar.Close()
	content, err = readFrom(bar)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("bar - fs"))

	baz, err := fsys.Open(m["baz.txt"])
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer baz.Close()
	content, err = readFrom(baz)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("baz - embedded"))

	// parent directory access is here to test the jail join
	forbiddenPath := filepath.Join(d, "../forbidden")
	g.Expect(os.Mkdir(forbiddenPath, 0700)).ToNot(gomega.HaveOccurred())
	forbiddenFile, err := createFile(forbiddenPath, "forbidden.txt", "forbidden - fs")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	forbiddenFile, err = os.Open(forbiddenFile.Name())
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer forbiddenFile.Close()
	content, err = readFrom(forbiddenFile)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("forbidden - fs"))
	_, err = fsys.Open(filepath.Join("../forbidden", filepath.Base(forbiddenFile.Name())))
	g.Expect(err).To(gomega.BeAssignableToTypeOf(&fs.PathError{}))
	g.Expect(err.Error()).To(gomega.ContainSubstring("file does not exist"))
}

func TestFallbackFS_OpenEmbedded(t *testing.T) {
	d := filepath.Join(t.TempDir(), "allowed")
	g := gomega.NewWithT(t)

	fsys, m, err := createFallbackFS(d)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	foo, err := fsys.OpenEmbedded(m["foo.txt"])
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer foo.Close()
	content, err := readFrom(foo)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("foo - embedded"))

	bar, err := fsys.OpenEmbedded(m["bar.txt"])
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer bar.Close()
	content, err = readFrom(bar)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("bar - embedded"))

	baz, err := fsys.OpenEmbedded(m["baz.txt"])
	g.Expect(err).ToNot(gomega.HaveOccurred())
	defer baz.Close()
	content, err = readFrom(baz)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("baz - embedded"))

}

func TestFallbackFS_Create(t *testing.T) {
	d := filepath.Join(t.TempDir(), "allowed")
	g := gomega.NewWithT(t)

	fsys, _, err := createFallbackFS(d)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	baz, err := fsys.Open("baz.txt")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	content, err := readFrom(baz)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("baz - embedded"))
	_ = baz.Close()

	// parent directory access is here to test the jail join
	f, err := fsys.Create("../../../../../../../baz.txt")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(f.Name()).To(gomega.Equal(filepath.Join(d, "baz.txt")))
	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.Write([]byte("baz - fs"))
	g.Expect(err).ToNot(gomega.HaveOccurred())

	baz, _ = fsys.Open("baz.txt")
	content, err = readFrom(baz)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(content).To(gomega.Equal("baz - fs"))
	_ = baz.Close()
}
