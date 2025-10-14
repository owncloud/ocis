package configlog

import (
	"fmt"
	"os"
)

// Error logs the error
func Error(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

// ReturnError logs the error and returns it unchanged
func ReturnError(err error) error {
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return err
}

// ReturnFatal logs the error and calls os.Exit(1) and returns nil if no error is passed
func ReturnFatal(err error) error {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return nil
}
