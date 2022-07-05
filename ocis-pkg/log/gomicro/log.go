package gomicro

import (
	"os"

	mzlog "github.com/go-micro/plugins/v4/logger/zerolog"
	"github.com/rs/zerolog"
	"go-micro.dev/v4/logger"
)

func init() {
	// this is ugly, but "logger.DefaultLogger" is a global variable and we need to set it _before_ anybody uses it
	setMicroLogger()
}

// for logging reasons we don't want the same logging level on both oCIS and micro. As a framework builder we do not
// want to expose to the end user the internal framework logs unless explicitly specified.
func setMicroLogger() {
	if os.Getenv("MICRO_LOG_LEVEL") == "" {
		_ = os.Setenv("MICRO_LOG_LEVEL", "error")
	}

	lev, err := zerolog.ParseLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if err != nil {
		lev = zerolog.ErrorLevel
	}
	logger.DefaultLogger = mzlog.NewLogger(
		logger.WithLevel(logger.Level(lev)),
		logger.WithFields(map[string]interface{}{
			"system": "go-micro",
		}),
	)
}
