package runner_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
)

// TimedTask will create a task with the specified duration
// The task will finish naturally after the given duration, or
// when it receives from the provided channel
//
// For the related stopper, just reuse the same channel:
//
//	func() {
//	  ch <- nil
//	  close(ch)
//	}
func TimedTask(ch chan error, dur time.Duration) runner.Runable {
	return func() error {
		timer := time.NewTimer(dur)
		defer timer.Stop()

		var result error
		select {
		case <-timer.C:
			// finish the task in 15 secs
		case result = <-ch:
			// or finish when we receive from the channel
		}
		return result
	}
}

var _ = Describe("Runner", func() {
	Describe("Run", func() {
		It("Context is done", func(ctx SpecContext) {
			// task will wait until it receives from the channel
			// stopper will just send something through the
			// channel, so the task can finish
			// Worst case, the task will finish after 15 secs
			ch := make(chan error)
			r := runner.New("run001", TimedTask(ch, 15*time.Second), func() {
				ch <- nil
				close(ch)
			})

			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r.Run(myCtx)
				close(ch2)
			}(ch2)

			expectedResult := &runner.Result{
				RunnerID:    "run001",
				RunnerError: nil,
			}

			// a result should be available in ch2 within the 5 secs spec
			// (task's context finishes in 1 sec so we expect a 1 sec delay)
			Eventually(ctx, ch2).Should(Receive(Equal(expectedResult)))
		}, SpecTimeout(5*time.Second))

		It("Context is done and interrupt after", func(ctx SpecContext) {
			// task will wait until it receives from the channel
			// stopper will just send something through the
			// channel, so the task can finish
			// Worst case, the task will finish after 15 secs
			ch := make(chan error)
			r := runner.New("run001", TimedTask(ch, 15*time.Second), func() {
				ch <- nil
				close(ch)
			})

			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r.Run(myCtx)
				close(ch2)
			}(ch2)

			expectedResult := &runner.Result{
				RunnerID:    "run001",
				RunnerError: nil,
			}

			// a result should be available in ch2 within the 5 secs spec
			// (task's context finishes in 1 sec so we expect a 1 sec delay)
			Eventually(ctx, ch2).Should(Receive(Equal(expectedResult)))

			r.Interrupt() // this shouldn't do anything
		}, SpecTimeout(5*time.Second))

		It("Task finishes naturally", func(ctx SpecContext) {
			e := errors.New("overslept!")
			r := runner.New("run002", func() error {
				time.Sleep(50 * time.Millisecond)
				return e
			}, func() {
			})

			// context will be done in 1 second (task will finishes before)
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r.Run(myCtx)
				close(ch2)
			}(ch2)

			expectedResult := &runner.Result{
				RunnerID:    "run002",
				RunnerError: e,
			}

			// a result should be available in ch2 within the 5 secs spec
			// (task finish naturally in 50 msec)
			Eventually(ctx, ch2).Should(Receive(Equal(expectedResult)))
		}, SpecTimeout(5*time.Second))

		It("Task doesn't finish", func(ctx SpecContext) {
			r := runner.New("run003", func() error {
				time.Sleep(20 * time.Second)
				return nil
			}, func() {
			})

			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r.Run(myCtx)
				close(ch2)
			}(ch2)

			// Task will finish naturally in 60 secs
			// Task's context will finish in 1 sec, but task won't receive
			// the notification and it will keep going
			Consistently(ctx, ch2).WithTimeout(4500 * time.Millisecond).ShouldNot(Receive())
		}, SpecTimeout(5*time.Second))

		It("Task doesn't finish and times out", func(ctx SpecContext) {
			r := runner.New("run003", func() error {
				time.Sleep(20 * time.Second)
				return nil
			}, func() {
			}, runner.WithInterruptDuration(3*time.Second))

			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r.Run(myCtx)
				close(ch2)
			}(ch2)

			var expectedResult *runner.Result
			// Task will finish naturally in 60 secs
			// Task's context will finish in 1 sec, but task won't receive
			// the notification and it will keep going
			// Task will time out in 3 seconds after being interrupted (when
			// context is done), so test should finish in 4 seconds
			Eventually(ctx, ch2).Should(Receive(&expectedResult))
			Expect(expectedResult.RunnerID).To(Equal("run003"))

			var timeoutError *runner.TimeoutError
			Expect(errors.As(expectedResult.RunnerError, &timeoutError)).To(BeTrue())
			Expect(timeoutError.RunnerID).To(Equal("run003"))
			Expect(timeoutError.Duration).To(Equal(3 * time.Second))
		}, SpecTimeout(5*time.Second))

		It("Run mutiple times panics", func(ctx SpecContext) {
			e := errors.New("overslept!")
			r := runner.New("run002", func() error {
				time.Sleep(50 * time.Millisecond)
				return e
			}, func() {
			})

			// context will be done in 1 second (task will finishes before)
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			Expect(func() {
				r.Run(myCtx)
				r.Run(myCtx)
			}).To(Panic())
		}, SpecTimeout(5*time.Second))
	})

	Describe("RunAsync", func() {
		It("Wait in channel", func(ctx SpecContext) {
			ch := make(chan *runner.Result)
			e := errors.New("Task has finished")

			r := runner.New("run004", func() error {
				time.Sleep(50 * time.Millisecond)
				return e
			}, func() {
			})

			r.RunAsync(ch)
			expectedResult := &runner.Result{
				RunnerID:    "run004",
				RunnerError: e,
			}

			Eventually(ctx, ch).Should(Receive(Equal(expectedResult)))
		}, SpecTimeout(5*time.Second))

		It("Run multiple times panics", func(ctx SpecContext) {
			ch := make(chan *runner.Result)
			e := errors.New("Task has finished")

			r := runner.New("run004", func() error {
				time.Sleep(50 * time.Millisecond)
				return e
			}, func() {
			})

			r.RunAsync(ch)

			Expect(func() {
				r.RunAsync(ch)
			}).To(Panic())
		}, SpecTimeout(5*time.Second))

		It("Interrupt async", func(ctx SpecContext) {
			ch := make(chan *runner.Result)
			e := errors.New("Task interrupted")

			taskCh := make(chan error)
			r := runner.New("run005", TimedTask(taskCh, 20*time.Second), func() {
				taskCh <- e
				close(taskCh)
			})

			r.RunAsync(ch)
			r.Interrupt()

			expectedResult := &runner.Result{
				RunnerID:    "run005",
				RunnerError: e,
			}

			Eventually(ctx, ch).Should(Receive(Equal(expectedResult)))
		}, SpecTimeout(5*time.Second))

		It("Interrupt async times out", func(ctx SpecContext) {
			ch := make(chan *runner.Result)
			e := errors.New("Task interrupted")

			r := runner.New("run005", func() error {
				time.Sleep(30 * time.Second)
				return e
			}, func() {
			}, runner.WithInterruptDuration(3*time.Second))

			r.RunAsync(ch)
			r.Interrupt()

			var expectedResult *runner.Result

			// Task will timeout after 3 second of receiving the interruption
			Eventually(ctx, ch).Should(Receive(&expectedResult))
			Expect(expectedResult.RunnerID).To(Equal("run005"))
			Expect(expectedResult.RunnerError.Error()).To(ContainSubstring("timed out"))
		}, SpecTimeout(5*time.Second))

		It("Interrupt async multiple times", func(ctx SpecContext) {
			ch := make(chan *runner.Result)
			e := errors.New("Task interrupted")

			taskCh := make(chan error)
			r := runner.New("run005", TimedTask(taskCh, 20*time.Second), func() {
				taskCh <- e
				close(taskCh)
			})

			r.RunAsync(ch)
			r.Interrupt()
			r.Interrupt()
			r.Interrupt()

			expectedResult := &runner.Result{
				RunnerID:    "run005",
				RunnerError: e,
			}

			Eventually(ctx, ch).Should(Receive(Equal(expectedResult)))
		}, SpecTimeout(5*time.Second))
	})

	Describe("Finished", func() {
		It("Finish channel closes", func(ctx SpecContext) {

			r := runner.New("run006", func() error {
				time.Sleep(50 * time.Millisecond)
				return nil
			}, func() {
			})

			// context will be done in 1 second (task will finishes before)
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r.Run(myCtx)
				close(ch2)
			}(ch2)

			finishedCh := r.Finished()

			Eventually(ctx, finishedCh).Should(BeClosed())
		}, SpecTimeout(5*time.Second))
	})
})
