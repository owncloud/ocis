package svc

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/asim/go-micro/plugins/events/nats/v4"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/owncloud/ocis/audit/pkg/config"
	"github.com/owncloud/ocis/audit/pkg/types"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

func startConsumer(c config.Eventstream, log log.Logger) (<-chan interface{}, error) {
	s, err := server.NewNatsStream(nats.Address(c.Address), nats.ClusterID(c.ClusterID))
	if err != nil {
		return nil, err
	}

	return events.Consume(s, "audit", events.ShareCreated{})
}

func startAuditLogger(c config.Auditlog, ch <-chan interface{}, log log.Logger) {
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

		b, err := marshal(auditEvent, c.Format)
		if err != nil {
			log.Error().Err(err).Msg("error marshaling the event")
			continue
		}

		if c.LogToConsole {
			log.Info().Msg(string(b))
		}

		if c.LogToFile {
			err := writeToFile(c.FilePath, b)
			if err != nil {
				log.Error().Err(err).Msg("error writing audit log file")
			}
		}

	}

}

func writeToFile(path string, ev []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := fmt.Fprintln(file, string(ev)); err != nil {
		return err
	}
	return nil
}

func marshal(ev interface{}, format string) ([]byte, error) {
	switch format {
	default:
		return nil, fmt.Errorf("unsupported format '%s'", format)
	case "json":
		return json.Marshal(ev)
	}
}
