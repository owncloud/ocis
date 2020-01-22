// +build windows

package client

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	cmd      = "url.dll,FileProtocolHandler"
	runDll32 = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
)

func open(input string) *exec.Cmd {
	cmd := exec.Command(runDll32, cmd, input)
	return cmd
}
