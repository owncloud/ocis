package fsx_test

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

func TestLayeredFS_Primary(t *testing.T) {
	g := gomega.NewWithT(t)
	primary := fsx.NewMemMapFs()
	fs := fsx.NewFallbackFS(primary, fsx.NewMemMapFs())

	g.Expect(primary).Should(gomega.BeIdenticalTo(fs.Primary().Fs))
}

func TestLayeredFS_Secondary(t *testing.T) {
	g := gomega.NewWithT(t)
	secondary := fsx.NewMemMapFs()
	fs := fsx.NewFallbackFS(fsx.NewMemMapFs(), secondary)

	g.Expect(secondary).Should(gomega.BeIdenticalTo(fs.Secondary().Fs))
}
