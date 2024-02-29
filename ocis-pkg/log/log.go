package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	mzlog "github.com/go-micro/plugins/v4/logger/zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-micro.dev/v4/logger"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

var (
	RequestIDString = "request-id"
	defaultLogger   atomic.Value
)

func init() {
	// this is ugly, but "logger.DefaultLogger" is a global variable and we need to set it _before_ anybody uses it
	setMicroLogger()

	// set the default logger for ocis
	setDefaultLogger()
}

// For logging reasons, we don't want the same logging level on both oCIS and micro.
// As a framework builder,
// we do not want to expose to the end user the internal framework logs unless explicitly specified.
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

// setDefaultLogger creates an environment-specific logger and sets it as the default logger.
func setDefaultLogger() {
	var opts []Option

	opts = append(opts,
		Pretty(os.Getenv("OCIS_LOG_PRETTY") == "true"),
		Color(os.Getenv("OCIS_LOG_COLOR") == "true"),
	)

	if level := os.Getenv("OCIS_LOG_LEVEL"); level != "" {
		opts = append(opts, Level(level))
	}

	if file := os.Getenv("OCIS_LOG_FILE"); file != "" {
		opts = append(opts, File(file))
	}

	l := NewLogger(opts...)
	defaultLogger.Store(&l)
}

// Default returns the default pre-configured logger.
func Default() *Logger {
	return defaultLogger.Load().(*Logger)
}

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

// NopLogger initializes a no-operation logger.
func NopLogger() Logger {
	return Logger{zerolog.Nop()}
}

type LineInfoHook struct{}

// Run is a hook to add line info to log messages.
// I found the zerolog example for this here:
// https://github.com/rs/zerolog/issues/22#issuecomment-1127295489
func (h LineInfoHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		e.Str("line", fmt.Sprintf("%s:%d", file, line))
	}
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

	if logLevel <= zerolog.InfoLevel {
		var lineInfoHook LineInfoHook
		logger = logger.Hook(lineInfoHook)
	}

	return Logger{
		logger,
	}
}

// SubloggerWithRequestID returns a sublogger with the x-request-id added to all events
func (l Logger) SubloggerWithRequestID(c context.Context) Logger {
	return Logger{
		l.With().Str(RequestIDString, chimiddleware.GetReqID(c)).Logger(),
	}
}

// Deprecation logs a deprecation message,
// it is used to inform the user that a certain feature is deprecated and will be removed in the future.
// Do not use a logger here because the message MUST be visible independent of the log level.
func Deprecation(a ...any) {
	fmt.Printf("\033[1;31mDEPRECATION: %s\033[0m\n", a...)
}
