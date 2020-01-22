// +build darwin

package client

import (
	"os/exec"
)

func open(input string) *exec.Cmd {
	return exec.Command("open", input)
}
