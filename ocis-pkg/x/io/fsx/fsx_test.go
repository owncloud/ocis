package fsx_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

func TestBase(t *testing.T) {
	g := gomega.NewWithT(t)
	wrapped := fsx.NewMemMapFs()
	fs := fsx.BaseFS{Fs: wrapped.Fs}

	g.Expect(wrapped.Fs).Should(gomega.BeIdenticalTo(fs.Fs))
}

func TestBase_IOFS(t *testing.T) {
	g := gomega.NewWithT(t)
	fs := fsx.BaseFS{Fs: fsx.NewMemMapFs()}
	g.Expect(reflect.TypeOf(fs.IOFS()).Name()).Should(gomega.Equal("IOFS"))
}

func TestFromIOFS(t *testing.T) {
	g := gomega.NewWithT(t)
	fs := fsx.FromIOFS(os.DirFS("."))

	g.Expect(reflect.TypeOf(fs.Fs).String()).Should(gomega.Equal("*afero.FromIOFS"))
}

func TestNewOsFs(t *testing.T) {
	g := gomega.NewWithT(t)
	fs := fsx.NewOsFs()

	g.Expect(reflect.TypeOf(fs.Fs).String()).Should(gomega.Equal("*afero.OsFs"))
}

func TestNewBasePathFs(t *testing.T) {
	g := gomega.NewWithT(t)
	base := fsx.NewMemMapFs()

	err := base.MkdirAll("first/foo/bar", 0755)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	err = base.MkdirAll("second/foo/baz", 0755)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	fs := fsx.NewBasePathFs(base, "first")

	info, err := fs.Stat("/foo/bar")
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(info.IsDir()).To(gomega.BeTrue())

	info, err = fs.Stat("/foo/baz")
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(info).To(gomega.BeNil())

	info, err = fs.Stat("../second/foo/baz")
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(info).To(gomega.BeNil())
}

func TestNewReadOnlyFs(t *testing.T) {
	g := gomega.NewWithT(t)
	base := fsx.NewMemMapFs()

	err := base.MkdirAll("first/foo/bar", 0755)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	fs := fsx.NewReadOnlyFs(base)
	err = fs.MkdirAll("first/foo/bay", 0755)
	g.Expect(err).To(gomega.HaveOccurred())
}
