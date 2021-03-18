package shutdowner

import (
	"os"
	"os/signal"
)

func WaitTermination() {
	done := make(chan struct{}, 1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		done <- struct{}{}
	}()

	<-done
}
