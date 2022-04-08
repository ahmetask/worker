package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ahmetask/worker"
)

type job struct {
	quit bool
}

func (j *job) doWork() {
	time.Sleep(1 * time.Second)
}

func (j *job) Do() {
	for {
		if j.quit {
			fmt.Println("goroutine is break ")
			break
		}
		j.doWork()
	}
}

func (j *job) Stop() {
	j.quit = true
}

func main() {
	pool := worker.NewWorkerPool(4, 4)
	pool.Start()

	pool.Submit(&job{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	pool.Stop()
}
