package registry

//
//import (
//	"context"
//	"testing"
//
//	"github.com/micro/go-micro/v2/registry"
//	"github.com/owncloud/ocis/v2/ocis-pkg/log"
//)
//
//func TestRegisterGRPCEndpoint(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	err := RegisterGRPCEndpoint(ctx, "test", "1234", "192.168.0.1:777", log.Logger{})
//	if err != nil {
//		t.Errorf("Unexpected error: %v", err)
//	}
//
//	s, err := registry.GetService("test")
//	if err != nil {
//		t.Errorf("Unexpected error: %v", err)
//	}
//
//	if len(s) != 1 {
//		t.Errorf("Expected exactly one service to be returned got %v", len(s))
//	}
//
//	if len(s[0].Nodes) != 1 {
//		t.Errorf("Expected exactly one node to be returned got %v", len(s[0].Nodes))
//	}
//
//	testSvc := s[0]
//	if testSvc.Name != "test" {
//		t.Errorf("Expected service name to be 'test' got %v", s[0].Name)
//	}
//
//	testNode := testSvc.Nodes[0]
//
//	if testNode.Address != "192.168.0.1:777" {
//		t.Errorf("Expected node address to be '192.168.0.1:777' got %v", testNode.Address)
//	}
//
//	if testNode.Id != "test-1234" {
//		t.Errorf("Expected node id to be 'test-1234' got %v", testNode.Id)
//	}
//
//	cancel()
//
// // When switching over to monorepo this little test fails. We're unsure of what the cause is, but since this test
// // is testing a framework specific behavior, we're better off letting it commented out. There is also no use of
// // com.owncloud.reva anywhere in the codebase, so we're effectively only registering reva as a go-micro service,
// // but not sending any message.
//	s, err = registry.GetService("test")
//	if err != nil {
//		t.Errorf("Unexpected error: %v", err)
//	}
//
//	if len(s) != 0 {
//		t.Errorf("Deregister on cancelation failed. Result-length should be zero, got %v", len(s))
//	}
//}

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	mRegistry "go-micro.dev/v4/registry"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// syncBuffer is written by the registration goroutine and read by the test
// goroutine, so it needs its own lock.
type syncBuffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

// TestRegisterServiceDoesNotLogEveryRefresh pins owncloud/ocis#8274: the
// periodic re-registration must not write a debug line on every tick.
//
// Two things make this test non-vacuous, and both are needed:
//
//   - Importing ocis-pkg/log runs its init(), which sets zerolog's *global*
//     level from MICRO_LOG_LEVEL (defaulting to go-micro's "error", which maps
//     onto zerolog's fatal). A logger built here therefore writes nothing at
//     all unless the global level is lowered first, and an assertion that no
//     debug line was written would pass for the wrong reason.
//   - The assertion only means something while the refresh loop is running.
//     The in-memory registry prunes any node it has not seen within its TTL,
//     so a service still resolvable well past one TTL is a service that is
//     being re-registered. If the loop were dead the lookup would fail and the
//     test would stop before asserting anything.
func TestRegisterServiceDoesNotLogEveryRefresh(t *testing.T) {
	const (
		interval = 100 * time.Millisecond
		ttl      = 1200 * time.Millisecond
		// Long enough for several refreshes and for the registry's own pruner
		// (it ticks once a second) to have removed an unrefreshed node.
		observe = 2500 * time.Millisecond
	)

	previous := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	t.Cleanup(func() { zerolog.SetGlobalLevel(previous) })

	t.Setenv(_registryEnv, "memory")
	t.Setenv(_registryRegisterIntervalEnv, interval.String())
	t.Setenv(_registryRegisterTTLEnv, ttl.String())

	ready := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ready.Close()

	buf := &syncBuffer{}
	logger := log.Logger{Logger: zerolog.New(buf).Level(zerolog.DebugLevel)}

	svc := &mRegistry.Service{
		Name:    "com.owncloud.test.refresh-logging",
		Version: "test",
		Nodes: []*mRegistry.Node{{
			Id:      "com.owncloud.test.refresh-logging-0",
			Address: "127.0.0.1:9999",
		}},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := RegisterService(ctx, logger, svc, strings.TrimPrefix(ready.URL, "http://")); err != nil {
		t.Fatalf("RegisterService returned an unexpected error: %v", err)
	}

	// The logger reaches the buffer at all -- otherwise every assertion below
	// would pass on an empty log.
	if got := buf.String(); !strings.Contains(got, "registering external service") {
		t.Fatalf("the test logger never reached the buffer; log was:\n%s", got)
	}

	time.Sleep(observe)

	// The node outlived its own TTL, so the refresh loop ran.
	found, err := GetRegistry().GetService(svc.Name)
	if err != nil {
		t.Fatalf("service was not refreshed within its TTL, so the assertion below would be vacuous: %v", err)
	}
	if len(found) == 0 {
		t.Fatal("service was not refreshed within its TTL, so the assertion below would be vacuous")
	}

	// The actual regression guard for #8274.
	if got := buf.String(); strings.Contains(got, "refreshing external service-registration") {
		t.Errorf("a debug line is written on every re-registration (#8274); log was:\n%s", got)
	}
}
