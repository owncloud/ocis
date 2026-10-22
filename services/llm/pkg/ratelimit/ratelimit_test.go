package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	microstore "go-micro.dev/v4/store"
)

func TestLimiter_AllowsUpToMax(t *testing.T) {
	l := New(microstore.NewMemoryStore(), time.Minute, 3)

	for i := 0; i < 3; i++ {
		allowed, err := l.Allow("user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("request %d should have been allowed", i+1)
		}
	}
}

func TestLimiter_DeniesOverMax(t *testing.T) {
	l := New(microstore.NewMemoryStore(), time.Minute, 2)

	for i := 0; i < 2; i++ {
		if _, err := l.Allow("user-1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	allowed, err := l.Allow("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("third request should have been denied")
	}
}

func TestLimiter_TracksUsersIndependently(t *testing.T) {
	l := New(microstore.NewMemoryStore(), time.Minute, 1)

	allowed, err := l.Allow("user-1")
	if err != nil || !allowed {
		t.Fatalf("user-1 first request should be allowed, got allowed=%v err=%v", allowed, err)
	}

	allowed, err = l.Allow("user-2")
	if err != nil || !allowed {
		t.Fatalf("user-2 first request should be allowed, got allowed=%v err=%v", allowed, err)
	}

	allowed, err = l.Allow("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("user-1 second request should have been denied")
	}
}

func TestLimiter_WindowExpiryClearsOldEntries(t *testing.T) {
	l := New(microstore.NewMemoryStore(), 10*time.Millisecond, 1)

	allowed, err := l.Allow("user-1")
	if err != nil || !allowed {
		t.Fatalf("first request should be allowed, got allowed=%v err=%v", allowed, err)
	}

	time.Sleep(20 * time.Millisecond)

	allowed, err = l.Allow("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("request after window expiry should have been allowed")
	}
}

func TestLimiter_DisabledWhenMaxRequestsIsZero(t *testing.T) {
	l := New(microstore.NewMemoryStore(), time.Minute, 0)

	for i := 0; i < 50; i++ {
		allowed, err := l.Allow("user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("request %d should have been allowed when the limiter is disabled", i+1)
		}
	}
}

// Two limiters sharing one store stand in for two llm replicas.
func TestLimiter_SharedAcrossReplicas(t *testing.T) {
	shared := microstore.NewMemoryStore()
	a := New(shared, time.Minute, 4)
	b := New(shared, time.Minute, 4)

	for i := 0; i < 2; i++ {
		for _, l := range []*Limiter{a, b} {
			allowed, err := l.Allow("user-1")
			if err != nil || !allowed {
				t.Fatalf("request should be allowed, got allowed=%v err=%v", allowed, err)
			}
		}
	}

	for _, l := range []*Limiter{a, b} {
		allowed, err := l.Allow("user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Fatal("request beyond the shared limit should have been denied")
		}
	}
}

func TestLimiter_ConcurrentRequestsNeverExceedMax(t *testing.T) {
	shared := microstore.NewMemoryStore()
	const max = 5

	var admitted int32
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// each goroutine is its own "replica" with no shared mutex
			allowed, err := New(shared, time.Minute, max).Allow("user-1")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed {
				atomic.AddInt32(&admitted, 1)
			}
		}(i)
	}
	wg.Wait()

	if admitted > max {
		t.Fatalf("admitted %d requests, limit is %d", admitted, max)
	}
}

func TestLimiter_DeniedRequestsDoNotCount(t *testing.T) {
	l := New(microstore.NewMemoryStore(), 50*time.Millisecond, 1)

	if allowed, _ := l.Allow("user-1"); !allowed {
		t.Fatal("first request should be allowed")
	}
	for i := 0; i < 5; i++ {
		if allowed, _ := l.Allow("user-1"); allowed {
			t.Fatal("request should be denied")
		}
	}

	time.Sleep(80 * time.Millisecond)

	if allowed, _ := l.Allow("user-1"); !allowed {
		t.Fatal("denied requests must not extend the window")
	}
}

func TestLimiter_KeyPrefixDoesNotLeak(t *testing.T) {
	l := New(microstore.NewMemoryStore(), time.Minute, 1)

	if allowed, _ := l.Allow("user"); !allowed {
		t.Fatal("first request should be allowed")
	}
	if allowed, _ := l.Allow("user/other"); !allowed {
		t.Fatal("a key that merely starts with another key must be tracked separately")
	}
}
