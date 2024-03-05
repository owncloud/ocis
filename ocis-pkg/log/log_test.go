package log_test

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/internal/testenv"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

func TestDeprecation(t *testing.T) {
	cmdTest := testenv.NewCMDTest(t.Name())
	if cmdTest.ShouldRun() {
		log.Deprecation("this is a deprecation")
		return
	}

	out, err := cmdTest.Run()

	g := gomega.NewWithT(t)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(string(out)).To(gomega.HavePrefix("\033[1;31mDEPRECATION: this is a deprecation"))
}
