package checks

import (
	"context"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"net"
	"strings"
	"time"
)

// NewTCPCheck returns a check that connects to a given tcp endpoint.
func NewTCPCheck(address string) func(context.Context) error {
	return func(_ context.Context) error {
		if strings.Contains(address, "0.0.0.0") || strings.Contains(address, "::") {
			outboundIp, err := handlers.GetOutBoundIP()
			if err != nil {
				return err
			}
			address = strings.Replace(address, "0.0.0.0", outboundIp, 1)
			address = strings.Replace(address, "::", outboundIp, 1)
			address = strings.Replace(address, "[::]", "["+outboundIp+"]", 1)
		}

		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			return err
		}

		err = conn.Close()
		if err != nil {
			return err
		}

		return nil
	}
}
