package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger simply wraps the zerolog logger.
type Logger struct {
	zerolog.Logger
}

// LoggerFromConfig initializes a service-specific logger instance.
func LoggerFromConfig(name string, cfg *shared.Log) Logger {
	return NewLogger(
		Name(name),
		Level(cfg.Level),
		Pretty(cfg.Pretty),
		Color(cfg.Color),
		File(cfg.File),
	)
}

// NewLogger initializes a new logger instance.
func NewLogger(opts ...Option) Logger {
	options := newOptions(opts...)

	// set GlobalLevel() to the minimum value -1 = TraceLevel, so that only the services' log level matter
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	var logLevel zerolog.Level
	switch strings.ToLower(options.Level) {
	case "panic":
		logLevel = zerolog.PanicLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "trace":
		logLevel = zerolog.TraceLevel
	default:
		logLevel = zerolog.ErrorLevel
	}

	var logger zerolog.Logger

	if options.Pretty {
		logger = log.Output(
			zerolog.NewConsoleWriter(
				func(w *zerolog.ConsoleWriter) {
					w.TimeFormat = time.RFC3339
					w.Out = os.Stderr
					w.NoColor = !options.Color
				},
			),
		)
	} else if options.File != "" {
		f, err := os.OpenFile(options.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			print(fmt.Sprintf("file could not be opened for writing: %s. error: %v", options.File, err))
			os.Exit(1)
		}
		logger = logger.Output(f)
	} else {
		logger = zerolog.New(os.Stderr)
	}

	logger = logger.With().
		Str("service", options.Name).
		Timestamp().
		Logger().Level(logLevel)

	return Logger{
		logger,
	}
}
