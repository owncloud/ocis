package inotifywaitgo

import (
	"errors"
	"fmt"
	"strings"
)

func GenerateShellCommands(s *Settings) ([]string, error) {
	if s.Options == nil {
		return nil, errors.New(OPT_NIL)
	}

	if s.Dir == "" {
		return nil, errors.New(DIR_EMPTY)
	}

	baseCmd := []string{
		"inotifywait",
		"-c", // switch to CSV output
	}

	if s.Options.Monitor {
		baseCmd = append(baseCmd, "-m")
	}

	if s.Options.Recursive {
		baseCmd = append(baseCmd, "-r")
	}

	if len(s.Options.Events) > 0 {
		for _, event := range s.Options.Events {
			// if event not in VALID_EVENTS
			if !Contains(VALID_EVENTS, int(event)) {
				return nil, errors.New(INVALID_EVENT)
			}
			baseCmd = append(baseCmd, "-e", EVENT_MAP[int(event)])
		}
	}

	baseCmd = append(baseCmd, s.Dir)

	// Trim spaces on all elements
	var outCmd []string
	for _, v := range baseCmd {
		outCmd = append(outCmd, strings.TrimSpace(v))
	}

	if s.Verbose {
		fmt.Println("Generated command:", strings.Join(outCmd, " "))
	}

	return outCmd, nil
}

// Contains checks if a slice contains an item
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
