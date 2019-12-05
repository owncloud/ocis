package log

import (
	"os"
	"strings"
	"time"

	mlog "github.com/micro/go-micro/util/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger simply wraps the zerolog logger.
type Logger struct {
	zerolog.Logger
}

// NewLogger initializes a new logger instance.
func NewLogger(opts ...Option) Logger {
	options := newOptions(opts...)

	switch strings.ToLower(options.Level) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
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
	} else {
		logger = zerolog.New(os.Stderr)
	}

	logger = logger.With().
		Str("service", options.Name).
		Timestamp().
		Logger()

	mlog.SetLogger(
		logWrapper{
			logger,
		},
	)

	return Logger{
		logger,
	}
}

// logWrapper implements the required interface for the go-micro logger.
type logWrapper struct {
	logger zerolog.Logger
}

// Log makes use of github.com/go-log/log.Log
func (w logWrapper) Log(v ...interface{}) {
	tmp := make([]string, len(v))

	for _, row := range v {
		tmp = append(tmp, row.(string))
	}

	w.logger.Info().Msg(strings.Join(tmp, " "))
}

// Logf makes use of github.com/go-log/log.Logf
func (w logWrapper) Logf(format string, v ...interface{}) {
	w.logger.Info().Msgf(strings.TrimRight(format, "\n"), v...)
}
