package runner_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
)

var _ = Describe("GroupRunner", func() {
	var (
		gr *runner.GroupRunner
	)

	BeforeEach(func() {
		gr = runner.NewGroup()

		task1Ch := make(chan error)
		task1 := TimedTask(task1Ch, 30*time.Second)
		gr.Add(runner.New("task1", task1, func() {
			task1Ch <- nil
			close(task1Ch)
		}))

		task2Ch := make(chan error)
		task2 := TimedTask(task2Ch, 20*time.Second)
		gr.Add(runner.New("task2", task2, func() {
			task2Ch <- nil
			close(task2Ch)
		}))
	})

	Describe("Add", func() {
		It("Duplicated runner id panics", func() {
			Expect(func() {
				gr.Add(runner.New("task1", func() error {
					time.Sleep(6 * time.Second)
					return nil
				}, func() {
				}))
			}).To(Panic())
		})

		It("Add after run panics", func(ctx SpecContext) {
			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan []*runner.Result)
			go func(ch2 chan []*runner.Result) {
				ch2 <- gr.Run(myCtx)
				close(ch2)
			}(ch2)

			// context is done in 1 sec, so all task should be interrupted and finish
			Eventually(ctx, ch2).Should(Receive(ContainElements(
				&runner.Result{RunnerID: "task1", RunnerError: nil},
				&runner.Result{RunnerID: "task2", RunnerError: nil},
			)))

			task3Ch := make(chan error)
			task3 := TimedTask(task3Ch, 6*time.Second)
			Expect(func() {
				gr.Add(runner.New("task3", task3, func() {
					task3Ch <- nil
					close(task3Ch)
				}))
			}).To(Panic())
		}, SpecTimeout(5*time.Second))

		It("Add after runAsync panics", func(ctx SpecContext) {
			ch2 := make(chan *runner.Result)
			gr.RunAsync(ch2)

			Expect(func() {
				task3Ch := make(chan error)
				task3 := TimedTask(task3Ch, 6*time.Second)
				gr.Add(runner.New("task3", task3, func() {
					task3Ch <- nil
					close(task3Ch)
				}))
			}).To(Panic())
		}, SpecTimeout(5*time.Second))
	})

	Describe("Run", func() {
		It("Context is done", func(ctx SpecContext) {
			task3Ch := make(chan error)
			task3 := TimedTask(task3Ch, 6*time.Second)
			gr.Add(runner.New("task3", task3, func() {
				task3Ch <- nil
				close(task3Ch)
			}))

			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan []*runner.Result)
			go func(ch2 chan []*runner.Result) {
				ch2 <- gr.Run(myCtx)
				close(ch2)
			}(ch2)

			// context is done in 1 sec, so all task should be interrupted and finish
			Eventually(ctx, ch2).Should(Receive(ContainElements(
				&runner.Result{RunnerID: "task1", RunnerError: nil},
				&runner.Result{RunnerID: "task2", RunnerError: nil},
				&runner.Result{RunnerID: "task3", RunnerError: nil},
			)))
		}, SpecTimeout(5*time.Second))

		It("One task finishes early", func(ctx SpecContext) {
			task3Ch := make(chan error)
			task3 := TimedTask(task3Ch, 1*time.Second)
			gr.Add(runner.New("task3", task3, func() {
				task3Ch <- nil
				close(task3Ch)
			}))

			// context will be done in 10 second
			myCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan []*runner.Result)
			go func(ch2 chan []*runner.Result) {
				ch2 <- gr.Run(myCtx)
				close(ch2)
			}(ch2)

			// task3 finishes in 1 sec, so the rest should also be interrupted
			Eventually(ctx, ch2).Should(Receive(ContainElements(
				&runner.Result{RunnerID: "task1", RunnerError: nil},
				&runner.Result{RunnerID: "task2", RunnerError: nil},
				&runner.Result{RunnerID: "task3", RunnerError: nil},
			)))
		}, SpecTimeout(5*time.Second))

		It("Context done and group timeout reached", func(ctx SpecContext) {
			gr := runner.NewGroup(runner.WithInterruptDuration(2 * time.Second))

			gr.Add(runner.New("task1", func() error {
				time.Sleep(6 * time.Second)
				return nil
			}, func() {
			}))

			gr.Add(runner.New("task2", func() error {
				time.Sleep(6 * time.Second)
				return nil
			}, func() {
			}))

			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan []*runner.Result)
			go func(ch2 chan []*runner.Result) {
				ch2 <- gr.Run(myCtx)
				close(ch2)
			}(ch2)

			// context finishes in 1 sec, tasks will be interrupted
			// group timeout will be reached after 2 extra seconds
			Eventually(ctx, ch2).Should(Receive(ContainElements(
				&runner.Result{RunnerID: "_unknown_", RunnerError: runner.NewGroupTimeoutError(2 * time.Second)},
				&runner.Result{RunnerID: "_unknown_", RunnerError: runner.NewGroupTimeoutError(2 * time.Second)},
			)))
		}, SpecTimeout(5*time.Second))

		It("Interrupted and group timeout reached", func(ctx SpecContext) {
			gr := runner.NewGroup(runner.WithInterruptDuration(2 * time.Second))

			gr.Add(runner.New("task1", func() error {
				time.Sleep(6 * time.Second)
				return nil
			}, func() {
			}))

			gr.Add(runner.New("task2", func() error {
				time.Sleep(6 * time.Second)
				return nil
			}, func() {
			}))

			// context will be done in 10 second
			myCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// spawn a new goroutine and return the result in the channel
			ch2 := make(chan []*runner.Result)
			go func(ch2 chan []*runner.Result) {
				ch2 <- gr.Run(myCtx)
				close(ch2)
			}(ch2)
			gr.Interrupt()

			// tasks will be interrupted
			// group timeout will be reached after 2 extra seconds
			Eventually(ctx, ch2).Should(Receive(ContainElements(
				&runner.Result{RunnerID: "_unknown_", RunnerError: runner.NewGroupTimeoutError(2 * time.Second)},
				&runner.Result{RunnerID: "_unknown_", RunnerError: runner.NewGroupTimeoutError(2 * time.Second)},
			)))
		}, SpecTimeout(5*time.Second))

		It("Doble run panics", func(ctx SpecContext) {
			// context will be done in 1 second
			myCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			Expect(func() {
				gr.Run(myCtx)
				gr.Run(myCtx)
			}).To(Panic())
		}, SpecTimeout(5*time.Second))
	})

	Describe("RunAsync", func() {
		It("Wait in channel", func(ctx SpecContext) {
			task3Ch := make(chan error)
			task3 := TimedTask(task3Ch, 1*time.Second)
			gr.Add(runner.New("task3", task3, func() {
				task3Ch <- nil
				close(task3Ch)
			}))

			ch2 := make(chan *runner.Result)
			gr.RunAsync(ch2)

			// task3 finishes in 1 sec, so the rest should also be interrupted
			Eventually(ctx, ch2).Should(Receive())
			Eventually(ctx, ch2).Should(Receive())
			Eventually(ctx, ch2).Should(Receive())
		}, SpecTimeout(5*time.Second))

		It("Double runAsync panics", func(ctx SpecContext) {
			ch2 := make(chan *runner.Result)
			Expect(func() {
				gr.RunAsync(ch2)
				gr.RunAsync(ch2)
			}).To(Panic())
		}, SpecTimeout(5*time.Second))

		It("Interrupt async", func(ctx SpecContext) {
			task3Ch := make(chan error)
			task3 := TimedTask(task3Ch, 6*time.Second)
			gr.Add(runner.New("task3", task3, func() {
				task3Ch <- nil
				close(task3Ch)
			}))

			ch2 := make(chan *runner.Result)
			gr.RunAsync(ch2)
			gr.Interrupt()

			// tasks will be interrupted
			Eventually(ctx, ch2).Should(Receive())
			Eventually(ctx, ch2).Should(Receive())
			Eventually(ctx, ch2).Should(Receive())
		}, SpecTimeout(5*time.Second))

		It("Interrupt async group timeout reached", func(ctx SpecContext) {
			gr := runner.NewGroup(runner.WithInterruptDuration(2 * time.Second))

			gr.Add(runner.New("task1", func() error {
				time.Sleep(6 * time.Second)
				return nil
			}, func() {
			}))

			gr.Add(runner.New("task2", func() error {
				time.Sleep(6 * time.Second)
				return nil
			}, func() {
			}))

			ch2 := make(chan *runner.Result)
			gr.RunAsync(ch2)
			gr.Interrupt()

			// group timeout will be reached after 2 extra seconds
			Eventually(ctx, ch2).Should(Receive(Equal(&runner.Result{RunnerID: "_unknown_", RunnerError: runner.NewGroupTimeoutError(2 * time.Second)})))
			Eventually(ctx, ch2).Should(Receive(Equal(&runner.Result{RunnerID: "_unknown_", RunnerError: runner.NewGroupTimeoutError(2 * time.Second)})))
		}, SpecTimeout(5*time.Second))
	})
})
