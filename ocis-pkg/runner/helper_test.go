package runner_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
)

var _ = Describe("Helper", func() {
	Describe("InterruptedTimeoutRunner", func() {
		It("Context done, no timeout", func(ctx SpecContext) {
			r1 := runner.New("task", func() error {
				time.Sleep(10 * time.Millisecond)
				return nil
			}, func() {
			})

			r2 := runner.InterruptedTimeoutRunner(r1, 2*time.Second)

			// context will be done in 1 second (task will finishes before)
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r2.Run(myCtx)
				close(ch2)
			}(ch2)

			expectedResult := &runner.Result{
				RunnerID:    "task",
				RunnerError: nil,
			}

			Eventually(ctx, ch2).Should(Receive(Equal(expectedResult)))
		}, SpecTimeout(5*time.Second))

		It("Context done, timeout reached", func(ctx SpecContext) {
			r1 := runner.New("task", func() error {
				time.Sleep(10 * time.Second)
				return nil
			}, func() {
			})

			r2 := runner.InterruptedTimeoutRunner(r1, 2*time.Second)

			// context will be done in 1 second (task will finishes before)
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan *runner.Result)
			go func(ch2 chan *runner.Result) {
				ch2 <- r2.Run(myCtx)
				close(ch2)
			}(ch2)

			var expectedResult *runner.Result
			Eventually(ctx, ch2).Should(Receive(&expectedResult))
			Expect(expectedResult.RunnerID).To(Equal("task"))
			Expect(expectedResult.RunnerError.Error()).To(ContainSubstring("Timeout reached"))
		}, SpecTimeout(5*time.Second))

		It("Interrupted, timeout reached", func(ctx SpecContext) {
			r1 := runner.New("task", func() error {
				time.Sleep(10 * time.Second)
				return nil
			}, func() {
			})

			r2 := runner.InterruptedTimeoutRunner(r1, 2*time.Second)

			ch2 := make(chan *runner.Result)
			r2.RunAsync(ch2)
			r2.Interrupt()

			var expectedResult *runner.Result
			Eventually(ctx, ch2).Should(Receive(&expectedResult))
			Expect(expectedResult.RunnerID).To(Equal("task"))
			Expect(expectedResult.RunnerError.Error()).To(ContainSubstring("Timeout reached"))
		}, SpecTimeout(5*time.Second))
	})
})
