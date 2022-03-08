package svc

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/audit/pkg/config"
	"github.com/owncloud/ocis/audit/pkg/types"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Log is used to log to different outputs
type Log func([]byte)

// Marshaller is used to marshal events
type Marshaller func(interface{}) ([]byte, error)

// AuditLoggerFromConfig will start a new AuditLogger generated from the config
func AuditLoggerFromConfig(cfg config.Auditlog, ch <-chan interface{}, log log.Logger) {
	var logs []Log

	if cfg.LogToConsole {
		logs = append(logs, WriteToStdout())
	}

	if cfg.LogToFile {
		logs = append(logs, WriteToFile(cfg.FilePath, log))
	}

	StartAuditLogger(ch, log, Marshal(cfg.Format, log), logs...)

}

// StartAuditLogger will block. run in seperate go routine
func StartAuditLogger(ch <-chan interface{}, log log.Logger, marshaller Marshaller, logto ...Log) {
	for {
		i := <-ch

		var auditEvent interface{}
		switch ev := i.(type) {
		case events.ShareCreated:
			auditEvent = types.ShareCreated(ev)
		default:
			log.Error().Interface("event", ev).Msg(fmt.Sprintf("can't handle event of type '%T'", ev))
			continue

		}

		b, err := marshaller(auditEvent)
		if err != nil {
			log.Error().Err(err).Msg("error marshaling the event")
			continue
		}

		for _, l := range logto {
			l(b)
		}
	}

}

// WriteToFile returns a Log function writing to a file
func WriteToFile(path string, log log.Logger) Log {
	return func(content []byte) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Error().Err(err).Msgf("error opening file '%s'", path)
			return
		}
		defer file.Close()
		if _, err := fmt.Fprintln(file, string(content)); err != nil {
			log.Error().Err(err).Msgf("error writing to file '%s'", path)
		}
	}
}

// WriteToStdout return a Log function writing to Stdout
func WriteToStdout() Log {
	return func(content []byte) {
		fmt.Println(string(content))
	}
}

// Marshal returns a Marshaller from the `format` string
func Marshal(format string, log log.Logger) Marshaller {
	switch format {
	default:
		log.Error().Msgf("unknown format '%s'", format)
		return nil
	case "json":
		return json.Marshal
	}
}
