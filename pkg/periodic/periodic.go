package periodic

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Task struct {
	Name      string
	Period    time.Duration
	RunOnInit bool
	Function  func()
}

func NewJobRunner(tasks ...Task) *jobRunner {
	return &jobRunner{
		tasks: tasks,

		shutdownInit: make(chan struct{}),
		shutdownDone: make(chan struct{}),
	}
}

type jobRunner struct {
	tasks []Task

	shutdownInit chan struct{}
	shutdownDone chan struct{}
}

func (jr jobRunner) Run() {
	ctx, cancelFn := context.WithCancel(context.TODO())

	for _, task := range jr.tasks {
		go func(ctx context.Context, task Task) {
			logrus.WithField("task", task.Name).Info("started")

			var sleep time.Duration

			if !task.RunOnInit {
				sleep = task.Period
			}

		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				case <-time.After(sleep):
					startSpan := time.Now()

					logrus.WithField("task", task.Name).Info("execution")
					task.Function()

					sleep = task.Period - time.Since(startSpan)
					if sleep < 0 {
						sleep = 0
					}
				}

			}

			logrus.WithField("task", task.Name).Info("finished")

		}(ctx, task)
	}

	<-jr.shutdownInit

	cancelFn()

	jr.shutdownDone <- struct{}{}
}

func (jr jobRunner) Shutdown() {
	jr.shutdownInit <- struct{}{}
	<-jr.shutdownDone
}
