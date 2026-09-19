package ratelimit

import (
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
