// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.CallerSkipFrameCount = 2
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

// Mode changes the logging format.
type Mode string

const (
	// JSONMode outputs JSON.
	JSONMode Mode = "json"
	// ConsoleMode outputs human-readable logs.
	ConsoleMode Mode = "console"
)

// Option is the option to use to configure the logger.
type Option func(l *zerolog.Logger)

// New creates a new logger.
func New(opts ...Option) *zerolog.Logger {
	// create a default logger
	zl := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	for _, opt := range opts {
		opt(&zl)
	}
	return &zl
}

// WithLevel is an option to configure the logging level.
func WithLevel(lvl string) Option {
	return func(l *zerolog.Logger) {
		zlvl := parseLevel(lvl)
		*l = l.Level(zlvl)
	}
}

// WithWriter is an option to configure the logging output.
func WithWriter(w io.Writer, m Mode) Option {
	return func(l *zerolog.Logger) {
		if m == ConsoleMode {
			*l = l.Output(zerolog.ConsoleWriter{Out: w, TimeFormat: "2006-01-02 15:04:05.999"})
		} else {
			*l = l.Output(w)
		}
	}
}

func parseLevel(v string) zerolog.Level {
	if v == "" {
		return zerolog.InfoLevel
	}

	lvl, err := zerolog.ParseLevel(v)
	if err != nil {
		return zerolog.InfoLevel
	}

	return lvl
}
