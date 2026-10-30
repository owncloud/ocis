package checks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	revastore "github.com/owncloud/reva/v2/pkg/store"
)

// NewNatsCheck checks the reachability of a nats server.
func NewNatsCheck(natsCluster string, options ...nats.Option) func(context.Context) error {
	return func(_ context.Context) error {
		n, err := nats.Connect(natsCluster, options...)
		if err != nil {
			return fmt.Errorf("could not connect to nats server: %v", err)
		}
		defer n.Close()
		if n.Status() != nats.CONNECTED {
			return fmt.Errorf("nats server not connected")
		}
		return nil
	}
}

// natsCheckTimeout limits a probe when the probe context has no deadline.
const natsCheckTimeout = 2 * time.Second

// NewNatsCheckWithConfig checks the reachability of the nats server at address.
// The address can contain more than one server, separated by commas.
// The check authenticates and uses TLS like the service connections configured by cfg.
// The check uses address and ignores the endpoint of cfg. An empty address means
// that the service does not use NATS. The function then returns nil.
//
// The check keeps one connection and uses it again for each probe. A new connection
// for each probe reads the .creds file each time. A probe can then read a rotated file
// before the write is complete. Each connection also counts against the connection
// limit of the NATS account. The nats client reconnects the connection when it loses
// the server. The connection reads the .creds file again on each reconnect.
//
// The connection status stays CONNECTED until the client notices a lost server.
// Thus each probe sends a PING to the server and waits for the PONG.
func NewNatsCheckWithConfig(address string, cfg stream.NatsConfig) func(context.Context) error {
	if address == "" {
		return nil
	}
	c := &natsCheck{address: address, cfg: cfg, lock: make(chan struct{}, 1)}
	return c.check
}

// NewNatsStoreCheck checks the reachability of the nodes of a store.
// The store connects to any node that it can reach. Thus the check succeeds when
// one node is reachable. A check of each node fails when one node of a healthy
// cluster is down.
// The function returns nil when the store type is not a NATS type or when the store has no nodes.
func NewNatsStoreCheck(storeType string, nodes []string, cfg stream.NatsConfig) func(context.Context) error {
	if !revastore.UsesNats(storeType) || len(nodes) == 0 {
		return nil
	}
	return NewNatsCheckWithConfig(strings.Join(nodes, ","), cfg)
}

type natsCheck struct {
	address string
	cfg     stream.NatsConfig

	// lock holds one token. A probe holds the token while it creates the connection.
	// A probe that waits for the token stops when its context ends.
	lock chan struct{}
	conn *nats.Conn
}

func (c *natsCheck) check(ctx context.Context) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, natsCheckTimeout)
		defer cancel()
	}

	conn, err := c.connection(ctx)
	if err != nil {
		return err
	}
	if !conn.IsConnected() {
		return fmt.Errorf("nats server not connected: status %v, last error %v", conn.Status(), conn.LastError())
	}
	if err := conn.FlushWithContext(ctx); err != nil {
		return fmt.Errorf("nats server not reachable: %v", err)
	}
	return nil
}

// connection returns the connection. It creates the connection on the first call.
// The nats client then connects again by itself, so the connection does not close.
// A failed configuration does not create a connection. The next probe tries again.
func (c *natsCheck) connection(ctx context.Context) (*nats.Conn, error) {
	select {
	case c.lock <- struct{}{}:
	case <-ctx.Done():
		return nil, fmt.Errorf("could not connect to nats server: %v", ctx.Err())
	}
	defer func() { <-c.lock }()

	if c.conn != nil {
		return c.conn, nil
	}

	options, err := stream.NatsOptions(c.cfg)
	if err != nil {
		return nil, fmt.Errorf("could not configure the nats connection: %v", err)
	}
	// The first connect blocks the probe. The timeout keeps the connect in the probe deadline.
	deadline, _ := ctx.Deadline()
	timeout := min(time.Until(deadline), natsCheckTimeout)
	if timeout <= 0 {
		return nil, fmt.Errorf("could not connect to nats server: %v", context.DeadlineExceeded)
	}
	options = append(options,
		nats.Timeout(timeout),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
	)
	conn, err := nats.Connect(c.address, options...)
	if err != nil {
		return nil, fmt.Errorf("could not connect to nats server: %v", err)
	}
	c.conn = conn
	return conn, nil
}
