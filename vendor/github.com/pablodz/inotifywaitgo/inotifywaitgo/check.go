package inotifywaitgo

import (
	"bufio"
	"os/exec"
)

// Function to checkDependencies if inotifywait is installed
func checkDependencies() (bool, error) {
	cmd := exec.Command("bash", "-c", "which inotifywait")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	if err := cmd.Start(); err != nil {
		return false, err
	}

	// Read the output of inotifywait and split it into lines
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			return true, nil
		}
	}
	return false, nil
}
