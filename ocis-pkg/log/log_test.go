package log_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

func TestDefault(t *testing.T) {
	if os.Getenv("TestDefault") == "true" {
		log.Default().Info().Msg("this is a test")
		return
	}

	g := gomega.NewWithT(t)

	tests := []struct {
		name     string
		args     []string
		validate func(result string)
	}{
		{
			name: "default",
			args: []string{"OCIS_LOG_PRETTY=false", "OCIS_LOG_COLOR=false"},
			validate: func(result string) {
				g.Expect(result).To(gomega.ContainSubstring("info"))
				g.Expect(result).To(gomega.ContainSubstring("this is a test"))
			},
		},
		{
			name: "default",
			args: []string{"OCIS_LOG_PRETTY=false", "OCIS_LOG_COLOR=false", "OCIS_LOG_LEVEL=error"},
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

			cmd := exec.Command(os.Args[0], "-test.run=TestDefault")
			cmd.Env = append(os.Environ(), "TestDefault=true")
			cmd.Env = append(cmd.Env, tt.args...)

			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			tt.validate(string(out))
		})
	}
}
