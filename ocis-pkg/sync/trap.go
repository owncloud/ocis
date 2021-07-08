package sync

import (
	"context"
	"os"
	"os/signal"

	"github.com/oklog/run"
)

// Trap listens to interrupt signals and handles context cancellation and channel closing on a group run.
func Trap(gr *run.Group, cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	gr.Add(func() error {
		signal.Notify(stop, os.Interrupt)
		<-stop
		return nil
	}, func(err error) {
		cancel()
	})
}
