package inotifywaitgo

import "os/exec"

func killOthers() error {
	cmd := exec.Command("bash", "-c", "pkill inotifywait").Run()
	return cmd
}
