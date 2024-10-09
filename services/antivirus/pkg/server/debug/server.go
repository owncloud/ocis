package debug

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/scanners"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	checkHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger).
			WithCheck("nats reachability", func(ctx context.Context) error {
				_, err := stream.NatsFromConfig("healthcheckfornats", false, stream.NatsConfig(options.Config.Events))
				if err != nil {
					return fmt.Errorf("could not connect to nats server: %v", err)
				}
				return nil
			}).
			WithCheck("antivirus reachability", func(ctx context.Context) error {
				cfg := options.Config
				switch cfg.Scanner.Type {
				default:
					// there is not av configured, return no error here
					return nil
				case "clamav":
					_, err := net.Dial("tcp", cfg.Scanner.ClamAV.Socket)
					if err != nil {
						return fmt.Errorf("could not connect to clamav server: %v", err)
					}
				case "icap":
					_, err := scanners.NewICAP(cfg.Scanner.ICAP.URL, cfg.Scanner.ICAP.Service, cfg.Scanner.ICAP.Timeout)
					if err != nil {
						return fmt.Errorf("could not connect to icap server: %v", err)
					}
				}
				return nil
			}),
	)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(checkHandler),
		debug.Ready(checkHandler),
	), nil
}
