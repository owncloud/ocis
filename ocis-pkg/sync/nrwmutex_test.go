package sync

import (
	"fmt"
	"runtime"
	"testing"
)

func HammerMutex(m *NRWMutex, loops int, c chan bool) {
	for i := 0; i < loops; i++ {
		id := fmt.Sprintf("%v", i)
		m.Lock(id)
		m.Unlock(id)
	}
	c <- true
}

func TestNRWMutex(t *testing.T) {
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	defer runtime.SetMutexProfileFraction(0)
	m := NewNRWMutex()
	c := make(chan bool)
	r := 10

	for i := 0; i < r; i++ {
		go HammerMutex(&m, 2000, c)
	}
	for i := 0; i < r; i++ {
		<-c
	}
}
