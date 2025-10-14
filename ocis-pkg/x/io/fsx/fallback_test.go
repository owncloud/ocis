package fsx_test

import (
	"io"
	"testing"

	"github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

func TestLayeredFS(t *testing.T) {
	g := gomega.NewWithT(t)

	read := func(fs fsx.FS, name string) (string, error) {
		f, err := fs.Open(name)
		if err != nil {
			return "", err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}

	mustRead := func(fs fsx.FS, name string) string {
		s, err := read(fs, name)
		g.Expect(err).ToNot(gomega.HaveOccurred())

		return s
	}

	create := func(fs fsx.FS, name, content string) {
		err := afero.WriteFile(fs, name, []byte(content), 0644)
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}

	primary := fsx.NewMemMapFs()
	create(primary, "both.txt", "primary")
	g.Expect(mustRead(primary, "both.txt")).To(gomega.Equal("primary"))
	create(primary, "primary.txt", "primary")
	g.Expect(mustRead(primary, "primary.txt")).To(gomega.Equal("primary"))

	secondary := fsx.NewMemMapFs()
	create(secondary, "both.txt", "secondary")
	g.Expect(mustRead(secondary, "both.txt")).To(gomega.Equal("secondary"))
	create(secondary, "secondary.txt", "secondary")
	g.Expect(mustRead(secondary, "secondary.txt")).To(gomega.Equal("secondary"))

	fs := fsx.NewFallbackFS(primary, fsx.NewReadOnlyFs(secondary))
	g.Expect(mustRead(fs, "both.txt")).To(gomega.Equal("primary"))
	g.Expect(mustRead(fs, "primary.txt")).To(gomega.Equal("primary"))
	g.Expect(mustRead(fs, "secondary.txt")).To(gomega.Equal("secondary"))

	create(fs, "fallback-fs.txt", "fallback-fs")
	g.Expect(mustRead(fs, "fallback-fs.txt")).To(gomega.Equal("fallback-fs"))
	g.Expect(mustRead(primary, "fallback-fs.txt")).To(gomega.Equal("fallback-fs"))
	g.Expect(mustRead(fs.Primary(), "fallback-fs.txt")).To(gomega.Equal("fallback-fs"))
	_, err := read(secondary, "fallback-fs.txt")
	g.Expect(err).To(gomega.HaveOccurred())
	_, err = read(fs.Secondary(), "fallback-fs.txt")
	g.Expect(err).To(gomega.HaveOccurred())
}

func TestLayeredFS_Primary(t *testing.T) {
	g := gomega.NewWithT(t)
	primary := fsx.NewMemMapFs()
	fs := fsx.NewFallbackFS(primary, fsx.NewMemMapFs())

	g.Expect(primary).To(gomega.BeIdenticalTo(fs.Primary().Fs))
}

func TestLayeredFS_Secondary(t *testing.T) {
	g := gomega.NewWithT(t)
	secondary := fsx.NewMemMapFs()
	fs := fsx.NewFallbackFS(fsx.NewMemMapFs(), secondary)

	g.Expect(secondary).To(gomega.BeIdenticalTo(fs.Secondary().Fs))
}
