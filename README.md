# Worker

Worker is a Golang library for scheduling and worker pool.

```go
 go get github.com/ahmetask/worker
```

## Usage

# Worker Pool
```go
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
```

## Scheduler

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ahmetask/worker"
)

// ScheduledJob1 job1
func ScheduledJob1(ctx context.Context) {
	log.Println("ScheduledJob1 Started")
	time.Sleep(1 * time.Second)
	log.Println("ScheduledJob1 Finished")
}

// ScheduledJob2 job2
func ScheduledJob2(_ context.Context) {
	log.Println("ScheduledJob2 Started")
	time.Sleep(1 * time.Second)
	log.Println("ScheduledJob2 Finished")
}

// ScheduledJob3 job3
func ScheduledJob3(_ context.Context) {
	log.Println("ScheduledJob3 Started")
	time.Sleep(3 * time.Second) //nolint:gomnd
	log.Println("ScheduledJob3 Finished")
}

func main() {
	scheduler := worker.NewScheduler()

	context1 := context.Background()
	trigger1, active := scheduler.Add(context1, ScheduledJob1, time.Second*10, false) //nolint:gomnd

	// Soft Start/Stop
	time.AfterFunc(10*time.Second, func() { //nolint:gomnd
		fmt.Println("ScheduledJob1 Starting")
		active <- true // or false for stopping the scheduler
	})

	// Specific Trigger
	time.AfterFunc(15*time.Second, func() { //nolint:gomnd
		fmt.Println("ScheduledJob1 Triggered")
		trigger1 <- true
	})

	context2 := context.Background()
	scheduler.Add(context2, ScheduledJob2, time.Minute*1, true)

	context3 := context.Background()
	scheduler.Add(context3, ScheduledJob3, time.Minute*1, true)

	// Manual Trigger
	time.AfterFunc(30*time.Second, func() { //nolint:gomnd
		fmt.Println("Trigger All")
		scheduler.TriggerAll()
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	scheduler.Stop()
}
```