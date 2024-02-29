package log_test

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/internal/testenv"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

func TestDefault(t *testing.T) {
	cmdTest := testenv.NewCMDTest(t.Name())
	if cmdTest.ShouldRun() {
		log.Default().Info().Msg("this is a test")
		return
	}

	g := gomega.NewWithT(t)

	tests := []struct {
		name     string
		env      []string
		validate func(result string)
	}{
		{
			name: "default",
			env:  []string{},
			validate: func(result string) {
				g.Expect(result).To(gomega.ContainSubstring("info"))
				g.Expect(result).To(gomega.ContainSubstring("this is a test"))
			},
		},
		{
			name: "error level",
			env:  []string{"OCIS_LOG_LEVEL=error"},
			validate: func(result string) {
				g.Expect(result).ToNot(gomega.ContainSubstring("info"))
				g.Expect(result).ToNot(gomega.ContainSubstring("this is a test"))
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out, err := cmdTest.Run(tt.env...)
			g.Expect(err).ToNot(gomega.HaveOccurred())

			tt.validate(string(out))
		})
	}
}

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
