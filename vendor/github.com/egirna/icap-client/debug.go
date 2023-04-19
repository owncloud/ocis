package icapclient

import (
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// the debug mode determiner & the writer to the write the debug output to
var (
	DEBUG       = false
	debugWriter io.Writer
	logger      *log.Logger
)

const (
	debugPrefix = "icap-client says: "
)

// SetDebugMode sets the debug mode for the entire package depending on the bool
func SetDebugMode(debug bool) {
	DEBUG = debug

	if DEBUG { // setting os.Stdout as the default debug writer if debug mode is enabled & also the debug prefix
		debugWriter = os.Stdout
		logger = log.New(debugWriter, debugPrefix, log.LstdFlags)

	}
}

// SetDebugOutput sets writer to write the debug outputs (default: os.Stdout)
func SetDebugOutput(w io.Writer) {
	debugWriter = w
	logger.SetOutput(debugWriter)
}

func logDebug(a ...interface{}) {
	if DEBUG {
		logger.Println(a...)
	}
}

func logfDebug(s string, a ...interface{}) {
	if DEBUG {
		logger.Printf(s, a...)
	}
}

func dumpDebug(a ...interface{}) {
	if DEBUG {
		spew.Fdump(debugWriter, a...)
	}
}
