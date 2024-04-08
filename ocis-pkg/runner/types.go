package runner

// Runable represent a task that can be executed by the Runner.
// It expected to be a long running task with an indefinite execution time,
// so it's suitable for servers or services.
// The task can eventually return an error, or nil if the execution finishes
// without errors
type Runable func() error

// Stopper represent a function that will stop the Runable.
// The stopper acts as a notification to the runable to know that the task
// needs to be finished now.
//
// The stopper won't need to crash the runable or force the runable to stop,
// instead, it will let the runable to know it has to stop and let it finish.
// This means that the runable might still run for a while.
//
// It's recommended the stopper to run asynchronously. This means that the
// stopper might need to spawn a goroutine. The intention is avoid blocking
// the running thread.
//
// Usually, the stoppers are the servers's `Shutdown()` or `Close()` methods,
// that will cause the server to start its shutdown procedure. As said, there
// is no need to force the shutdown, so graceful shutdowns are preferred if
// they're available
type Stopper func()

// Result represents the result of a runner.
// The result contains the provided runner's id (for easier identification
// in case of multiple results) and the runner's error, which is the result
// of the Runable function (might be nil if no error happened)
type Result struct {
	RunnerID    string
	RunnerError error
}
