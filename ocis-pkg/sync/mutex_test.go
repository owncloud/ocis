package sync

import (
	"fmt"
	"runtime"
	"testing"
)

func HammerMutex(m *NamedRWMutex, loops int, c chan bool) {
	for i := range loops {
		id := fmt.Sprintf("%v", i)
		m.Lock(id)
		m.Unlock(id)
	}
	c <- true
}

func TestNamedRWMutex(t *testing.T) {
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	defer runtime.SetMutexProfileFraction(0)
	m := NewNamedRWMutex()
	c := make(chan bool)
	r := 10

	for range r {
		go HammerMutex(&m, 2000, c)
	}
	for range r {
		<-c
	}
}
