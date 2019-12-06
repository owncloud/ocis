package log

import (
	"fmt"
	"os"
	"runtime"
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
		mlog.SetLevel(mlog.LevelFatal)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		mlog.SetLevel(mlog.LevelFatal)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		mlog.SetLevel(mlog.LevelError)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		mlog.SetLevel(mlog.LevelWarn)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		mlog.SetLevel(mlog.LevelInfo)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		mlog.SetLevel(mlog.LevelDebug)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		mlog.SetLevel(mlog.LevelTrace)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		mlog.SetLevel(mlog.LevelInfo)
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
		microZerolog{
			logger: logger,
		},
	)

	return Logger{
		logger,
	}
}

// microZerolog implements the required interface for the go-micro logger.
type microZerolog struct {
	logger zerolog.Logger
}

// Log makes use of github.com/go-log/log.Log.
func (mz microZerolog) Log(v ...interface{}) {
	pc := parentCaller()
	msg := fmt.Sprint(v...)

	switch {
	case strings.HasSuffix(pc, "Fatal"):
		mz.logger.Fatal().Msg(msg)
	case strings.HasSuffix(pc, "Error"):
		mz.logger.Error().Msg(msg)
	case strings.HasSuffix(pc, "Info"):
		mz.logger.Info().Msg(msg)
	case strings.HasSuffix(pc, "Warn"):
		mz.logger.Warn().Msg(msg)
	case strings.HasSuffix(pc, "Debug"):
		mz.logger.Debug().Msg(msg)
	case strings.HasSuffix(pc, "Trace"):
		mz.logger.Debug().Msg(msg)
	default:
		mz.logger.Info().Msg(msg)
	}
}

// Logf makes use of github.com/go-log/log.Logf.
func (mz microZerolog) Logf(format string, v ...interface{}) {
	pc := parentCaller()
	msg := fmt.Sprintf(strings.TrimRight(format, "\n"), v...)

	switch {
	case strings.HasSuffix(pc, "Fatalf"):
		mz.logger.Fatal().Msg(msg)
	case strings.HasSuffix(pc, "Errorf"):
		mz.logger.Error().Msg(msg)
	case strings.HasSuffix(pc, "Infof"):
		mz.logger.Info().Msg(msg)
	case strings.HasSuffix(pc, "Warnf"):
		mz.logger.Warn().Msg(msg)
	case strings.HasSuffix(pc, "Debugf"):
		mz.logger.Debug().Msg(msg)
	case strings.HasSuffix(pc, "Tracef"):
		mz.logger.Debug().Msg(msg)
	default:
		mz.logger.Info().Msg(msg)
	}
}

// parentCaller tries to detect which log method had been invoked.
func parentCaller() string {
	pc, _, _, ok := runtime.Caller(4)
	fn := runtime.FuncForPC(pc)

	if ok && fn != nil {
		return fn.Name()
	}

	return ""
}
