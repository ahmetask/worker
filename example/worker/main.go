package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ahmetask/worker"
)

const workerCount = 4
const queueCapacity = 4

// Job1 type
type Job1 struct {
	ID int
}

// Do implement work interface
func (j *Job1) Do() {
	time.Sleep(5 * time.Second) //nolint:gomnd
	log.Printf("Job1 Finished:%d\n", j.ID)
}

// Stop job1
func (j *Job1) Stop() {
	log.Printf("Stop Job1:%d\n", j.ID)
}

// Job2 type
type Job2 struct {
	ID int64
}

// Do implement work interface
func (j *Job2) Do() {
	time.Sleep(5 * time.Second) //nolint:gomnd
	log.Printf("Job2 Finished:%d\n", j.ID)
}

// Stop job2
func (j *Job2) Stop() {
	log.Printf("Stop Job2:%d\n", j.ID)
}

func main() {
	pool := worker.NewWorkerPool(workerCount, queueCapacity)
	pool.Start()

	queued := pool.EnqueueWithTimeout(&Job1{
		ID: 1,
	}, 1*time.Second)
	log.Printf("Enqueue Job Finished: %d Result:%v\n", 1, queued)

	queued = pool.Enqueue(&Job2{
		ID: 2,
	})
	log.Printf("Enqueue Job Finished: %d Result:%v\n", 2, queued) //nolint:gomnd

	pool.Submit(&Job2{
		ID: 3,
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	pool.Stop()
}
