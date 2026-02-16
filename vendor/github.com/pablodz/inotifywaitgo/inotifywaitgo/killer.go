package inotifywaitgo

import "os/exec"

func killOthers() error {
	cmd := exec.Command("pkill", "inotifywait")
	return cmd.Run()
}
