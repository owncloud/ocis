package runner

import (
	"context"
)

// Runner represents the one executing a long running task, such as a server
// or a service.
// The ID of the runner is public to make identification easier, and the
// Result that it will generated will contain the same ID, so we can
// know which runner provided which result.
type Runner struct {
	ID        string
	fn        Runable
	interrupt Stopper
}

// New will create a new runner.
// The runner will be created with the provided id (the id must be unique,
// otherwise undefined behavior might occur), and will run the provided
// runable task, using the "interrupt" function to stop that task if needed.
//
// Note that it's your responsibility to provide a proper stopper for the task.
// The runner will just call that method assuming it will be enough to
// eventually stop the task at some point.
func New(id string, fn Runable, interrupt Stopper) *Runner {
	return &Runner{
		ID:        id,
		fn:        fn,
		interrupt: interrupt,
	}
}

// Run will execute the task associated to this runner in a synchronous way.
// The task will be spawned in a new goroutine, and the current thread will
// wait until the task finishes.
//
// The task will finish "naturally". The stopper will be called in the
// following ways:
// - Manually calling this runner's `Interrupt` method
// - When the provided context is done
// As said, it's expected that calling the provided stopper will be enough to
// make the task to eventually complete.
//
// Once the task finishes, the result will be returned.
//
// Some nice things you can do:
// - Use signal.NotifyContext(...) to call the stopper and provide a clean
// shutdown procedure when an OS signal is received
// - Use context.WithDeadline(...) or context.WithTimeout(...) to run the task
// for a limited time
func (r *Runner) Run(ctx context.Context) *Result {
	ch := make(chan *Result)

	go r.doTask(ch, true)

	select {
	case result := <-ch:
		return result
	case <-ctx.Done():
		r.interrupt()
	}

	return <-ch
}

// RunAsync will execute the task associated to this runner asynchronously.
// The task will be spawned in a new goroutine and this method will finish.
// The task's result will be written in the provided channel when it's
// available, so you can wait for it if needed. It's up to you to decide
// to use a blocking or non-blocking channel, but the task will always finish
// before writing in the channel.
//
// To interrupt the running task, the only option is to call the `Interrupt`
// method at some point.
func (r *Runner) RunAsync(ch chan<- *Result) {
	go r.doTask(ch, false)
}

// Interrupt will execute the stopper function, which should notify the task
// in order for it to finish.
// The stopper will be called immediately, although it's expected the
// consequences to take a while (task might need a while to stop)
func (r *Runner) Interrupt() {
	r.interrupt()
}

// doTask will perform this runner's task and write the result in the provided
// channel. The channel will be closed if requested.
func (r *Runner) doTask(ch chan<- *Result, closeChan bool) {
	err := r.fn()

	result := &Result{
		RunnerID:    r.ID,
		RunnerError: err,
	}
	ch <- result

	if closeChan {
		close(ch)
	}
}
