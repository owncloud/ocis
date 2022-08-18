package configlog

import (
	"fmt"
	"os"
)

// LogError logs the error
func LogError(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

// LogError logs the error and returns it unchanged
func LogReturnError(err error) error {
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return err
}

// LogReturnFatal logs the error and calls os.Exit(1) and returns nil if no error is passed
func LogReturnFatal(err error) error {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return nil
}
