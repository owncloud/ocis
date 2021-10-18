package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mdlog "go-micro.dev/v4/debug/log"
	mlog "go-micro.dev/v4/util/log"
	"go-micro.dev/v4/util/ring"
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
		Logger()

	mlog.SetLogger(
		microZerolog{
			logger: logger,
			buffer: ring.New(mdlog.DefaultSize),
		},
	)

	return Logger{
		logger,
	}
}

// microZerolog implements the required interface for the go-micro logger.
type microZerolog struct {
	logger zerolog.Logger
	buffer *ring.Buffer
}

func (mz microZerolog) Read(opts ...mdlog.ReadOption) ([]mdlog.Record, error) {
	options := mdlog.ReadOptions{}
	for _, o := range opts {
		o(&options)
	}

	var entries []*ring.Entry

	if !options.Since.IsZero() {
		entries = mz.buffer.Since(options.Since)
	}

	if options.Count > 0 {
		switch len(entries) > 0 {
		case true:
			if options.Count > len(entries) {
				entries = entries[0:options.Count]
			}
		default:
			entries = mz.buffer.Get(options.Count)
		}
	}

	records := make([]mdlog.Record, 0, len(entries))
	for _, entry := range entries {
		record := mdlog.Record{
			Timestamp: entry.Timestamp,
			Message:   entry.Value,
		}
		records = append(records, record)
	}

	return records, nil
}

func (mz microZerolog) Write(record mdlog.Record) error {
	level := record.Metadata["level"]
	mz.log(level, fmt.Sprint(record.Message))
	mz.buffer.Put(record.Message)
	return nil
}

func (mz microZerolog) Stream() (mdlog.Stream, error) {
	stream, stop := mz.buffer.Stream()
	records := make(chan mdlog.Record, 128)
	last10 := mz.buffer.Get(10)

	go func() {
		for _, entry := range last10 {
			records <- mdlog.Record{
				Timestamp: entry.Timestamp,
				Message:   entry.Value,
				Metadata:  make(map[string]string),
			}
		}
		for entry := range stream {
			records <- mdlog.Record{
				Timestamp: entry.Timestamp,
				Message:   entry.Value,
				Metadata:  make(map[string]string),
			}
		}
	}()
	return &logStream{
		stream: records,
		stop:   stop,
	}, nil
}

func (mz microZerolog) log(level string, msg string) {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		l = zerolog.InfoLevel
	}

	mz.logger.WithLevel(l).Msg(msg)

	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if l == zerolog.FatalLevel {
		os.Exit(1)
	}
}
