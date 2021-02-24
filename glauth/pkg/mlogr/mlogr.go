package mlogr

import (
	"errors"

	"github.com/go-logr/logr"
	plog "github.com/owncloud/ocis/ocis-pkg/log"

	"github.com/rs/zerolog"
)

const debugVerbosity = 6
const traceVerbosity = 8

// New returns a logr.Logger which is implemented by the log.
func New(l *plog.Logger) logr.Logger {
	return logger{
		l:         l,
		verbosity: 0,
		prefix:    "glauth",
		values:    nil,
	}
}

// logger is a logr.Logger that uses the ocis-pkg log.
type logger struct {
	l         *plog.Logger
	verbosity int
	prefix    string
	values    []interface{}
}

func (l logger) clone() logger {
	out := l
	out.values = copySlice(l.values)
	return out
}

func copySlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

// add converts a bunch of arbitrary key-value pairs into zerolog fields.
func add(e *zerolog.Event, keysAndVals []interface{}) {

	// make sure we got an even number of arguments
	if len(keysAndVals)%2 != 0 {
		e.Interface("args", keysAndVals).
			AnErr("zerologr-err", errors.New("odd number of arguments passed as key-value pairs for logging")).
			Stack()
		return
	}

	for i := 0; i < len(keysAndVals); {
		// process a key-value pair,
		// ensuring that the key is a string
		key, val := keysAndVals[i], keysAndVals[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, log additional error
			e.Interface("invalid key", key).
				AnErr("zerologr-err", errors.New("non-string key argument passed to logging, ignoring all later arguments")).
				Stack()
			return
		}
		e.Interface(keyStr, val)

		i += 2
	}
}

func (l logger) Info(msg string, keysAndVals ...interface{}) {
	if l.Enabled() {
		var e *zerolog.Event
		if l.verbosity < debugVerbosity {
			e = l.l.Info()
		} else if l.verbosity < traceVerbosity {
			e = l.l.Debug()
		} else {
			e = l.l.Trace()
		}
		e.Int("verbosity", l.verbosity)
		if l.prefix != "" {
			e.Str("name", l.prefix)
		}
		add(e, l.values)
		add(e, keysAndVals)
		e.Msg(msg)
	}
}

func (l logger) Enabled() bool {
	return true
}

func (l logger) Error(err error, msg string, keysAndVals ...interface{}) {
	e := l.l.Error().Err(err)
	if l.prefix != "" {
		e.Str("name", l.prefix)
	}
	add(e, l.values)
	add(e, keysAndVals)
	e.Msg(msg)
}

func (l logger) V(verbosity int) logr.InfoLogger {
	//new := l.clone()
	//new.level = level
	//return new
	l.verbosity = verbosity
	return l
}

// WithName returns a new logr.Logger with the specified name appended. zerologr
// uses '/' characters to separate name elements.  Callers should not pass '/'
// in the provided name string, but this library does not actually enforce that.
func (l logger) WithName(name string) logr.Logger {
	nl := l.clone()
	if len(l.prefix) > 0 {
		nl.prefix = l.prefix + "/"
	}
	nl.prefix += name
	return nl
}
func (l logger) WithValues(kvList ...interface{}) logr.Logger {
	nl := l.clone()
	nl.values = append(nl.values, kvList...)
	return nl
}

var _ logr.Logger = logger{}
var _ logr.InfoLogger = logger{}
