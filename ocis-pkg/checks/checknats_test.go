package checks_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/natstest"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/stretchr/testify/require"
)

func TestNewNatsCheckWithConfig(t *testing.T) {
	s, creds := natstest.StartOperatorServer(t)
	url := s.ClientURL()

	t.Run("empty address", func(t *testing.T) {
		require.Nil(t, checks.NewNatsCheckWithConfig("", stream.NatsConfig{}))
	})
	credsFile := filepath.Join(t.TempDir(), "user.creds")
	require.NoError(t, os.WriteFile(credsFile, creds, 0o600))

	t.Run("no credentials", func(t *testing.T) {
		require.Error(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{})(context.Background()))
	})

	t.Run("username and password", func(t *testing.T) {
		require.Error(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{
			AuthUsername: "user",
			AuthPassword: "pass",
		})(context.Background()))
	})

	t.Run("creds file", func(t *testing.T) {
		require.NoError(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{
			CredsFile: credsFile,
		})(context.Background()))
	})

	t.Run("creds content", func(t *testing.T) {
		require.NoError(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{
			Creds:        string(creds),
			AuthUsername: "ignored",
			AuthPassword: "ignored",
		})(context.Background()))
	})

	t.Run("invalid creds", func(t *testing.T) {
		err := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{Creds: "not a creds file"})(context.Background())
		require.ErrorContains(t, err, "could not configure the nats connection")
	})

	t.Run("probes use the same connection", func(t *testing.T) {
		rotatedFile := filepath.Join(t.TempDir(), "user.creds")
		require.NoError(t, os.WriteFile(rotatedFile, creds, 0o600))
		check := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{CredsFile: rotatedFile})
		require.NoError(t, check(context.Background()))

		// A rotation writes the file. The connection is open, so the probe does not read the file.
		require.NoError(t, os.WriteFile(rotatedFile, creds[:len(creds)/3], 0o600))
		require.NoError(t, check(context.Background()))
	})

	t.Run("probe retries invalid creds", func(t *testing.T) {
		laterFile := filepath.Join(t.TempDir(), "user.creds")
		check := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{CredsFile: laterFile})
		require.ErrorContains(t, check(context.Background()), "could not configure the nats connection")

		require.NoError(t, os.WriteFile(laterFile, creds, 0o600))
		require.NoError(t, check(context.Background()))
	})

	t.Run("probe fails after the server stops", func(t *testing.T) {
		s, creds := natstest.StartOperatorServer(t)
		url := s.ClientURL()
		check := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{Creds: string(creds)})
		require.NoError(t, check(context.Background()))

		s.Shutdown()
		require.Error(t, check(context.Background()))
	})
}

func TestNewNatsCheckWithConfigUnresponsiveServer(t *testing.T) {
	// The listener accepts connections, but it never sends the NATS INFO message.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = l.Close() })
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { _ = c.Close() })
		}
	}()

	check := checks.NewNatsCheckWithConfig("nats://"+l.Addr().String(), stream.NatsConfig{})
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			start := time.Now()
			require.Error(t, check(ctx))
			require.Less(t, time.Since(start), time.Second)
		}()
	}
	wg.Wait()
}

func TestNewNatsStoreCheck(t *testing.T) {
	s, creds := natstest.StartOperatorServer(t)
	url := s.ClientURL()
	cfg := stream.NatsConfig{Creds: string(creds)}

	t.Run("nats-js-kv", func(t *testing.T) {
		require.NoError(t, checks.NewNatsStoreCheck("nats-js-kv", []string{url}, cfg)(context.Background()))
	})

	t.Run("nats-js", func(t *testing.T) {
		require.NoError(t, checks.NewNatsStoreCheck("nats-js", []string{url}, cfg)(context.Background()))
	})

	t.Run("nats store without credentials", func(t *testing.T) {
		require.Error(t, checks.NewNatsStoreCheck("nats-js-kv", []string{url}, stream.NatsConfig{})(context.Background()))
	})

	t.Run("other store types", func(t *testing.T) {
		for _, storeType := range []string{"memory", "redis-sentinel", "noop", ""} {
			require.Nil(t, checks.NewNatsStoreCheck(storeType, []string{url}, stream.NatsConfig{}), storeType)
		}
	})

	t.Run("no nodes", func(t *testing.T) {
		require.Nil(t, checks.NewNatsStoreCheck("nats-js-kv", nil, stream.NatsConfig{}))
	})

	t.Run("one node is down", func(t *testing.T) {
		// Nothing listens on the first node.
		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		down := "nats://" + l.Addr().String()
		require.NoError(t, l.Close())

		require.NoError(t, checks.NewNatsStoreCheck("nats-js-kv", []string{down, url}, cfg)(context.Background()))
	})

	t.Run("all nodes are down", func(t *testing.T) {
		s, creds := natstest.StartOperatorServer(t)
		check := checks.NewNatsStoreCheck("nats-js-kv", []string{s.ClientURL()}, stream.NatsConfig{Creds: string(creds)})
		s.Shutdown()
		require.Error(t, check(context.Background()))
	})
}
