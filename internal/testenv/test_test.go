package testenv

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"
)

func TestNewSubTest(t *testing.T) {
	testString := "this is a sub-test"
	cmdTest := NewCMDTest(t.Name())
	if cmdTest.ShouldRun() {
		fmt.Println(testString)
		return
	}

	out, err := cmdTest.Run()

	g := gomega.NewWithT(t)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(string(out)).To(gomega.ContainSubstring(testString))
}
