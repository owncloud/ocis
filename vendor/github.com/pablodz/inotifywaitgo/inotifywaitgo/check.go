package inotifywaitgo

import (
	"os/exec"
)

// CheckDependencies verifies if inotifywait is installed.
func checkDependencies() (bool, error) {
	path, err := exec.LookPath("inotifywait")
	if err != nil {
		return false, err
	}
	if path != "" {
		return true, nil
	}
	return false, nil
}
