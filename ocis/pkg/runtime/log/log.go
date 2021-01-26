package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	// Level sets a project wide log level
	Level zerolog.Level = zerolog.InfoLevel
)

// NewLogger configures a logger.
func NewLogger(options ...Option) zerolog.Logger {
	zerolog.SetGlobalLevel(Level)

	o := NewOptions()
	for _, f := range options {
		f(o)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if o.Pretty {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	return logger
}
