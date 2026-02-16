package testenv

import (
	"fmt"
	"os"
	"os/exec"
)

// CMDTest spawns a new independent test environment
type CMDTest struct {
	n string
}

// NewCMDTest creates a new CMDTest instance
func NewCMDTest(name string) CMDTest {
	return CMDTest{
		n: name,
	}
}

// Run runs the cmd subtest
func (t CMDTest) Run(envs ...string) ([]byte, error) {
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", t.n))
	cmd.Env = append(os.Environ(), "RUN_CMD_TEST=1")
	cmd.Env = append(cmd.Env, envs...)

	return cmd.CombinedOutput()
}

// ShouldRun checks if the cmd subtest should run
func (CMDTest) ShouldRun() bool {
	return os.Getenv("RUN_CMD_TEST") == "1"
}
