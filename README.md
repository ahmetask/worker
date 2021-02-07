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
	"fmt"
	"github.com/ahmetask/worker"
	"log"
	"os"
	"os/signal"
	"time"
)

type Job struct {
	//any interface value that you need
	Id int
}

/*implement work interface*/
func (j *Job) Do() {
	time.Sleep(5 * time.Second)
	log.Println(fmt.Sprintf("Job Finished:%d", j.Id))
}

func main() {
	// Initialize Pool
	// First Param=> Worker Count Second Param is Waiting Queue Capacity
	pool := worker.NewWorkerPool(1, 1)

	//Start worker pool
	pool.Start()

	//This is blocking if all workers are busy
	pool.Submit(&Job{Id: 1})

	//Tries to enqueue but fails if queue is full*/
	a := pool.Enqueue(&Job{Id: 2})
	log.Println(fmt.Sprintf("Queued 2: %v", a))

	b := pool.Enqueue(&Job{Id: 3})
	log.Println(fmt.Sprintf("Queued 3: %v", b))

	c := pool.Enqueue(&Job{Id: 4})
	log.Println(fmt.Sprintf("Queued 4: %v", c))

	e := pool.EnqueueWithTimeout(&Job{Id: 5}, 1*time.Second)
	log.Println(fmt.Sprintf("Queued 5 : %v", e))
	/*try to enqueue but fails if timeout occurs*/
	d := pool.EnqueueWithTimeout(&Job{Id: 6}, 10*time.Second)
	log.Println(fmt.Sprintf("Queued 6 : %v", d))

	//Waiting
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//Stop Worker Pool
	pool.Stop()
}

```

## Scheduler

```go

package main

import (
	"context"
	"github.com/ahmetask/worker"
	"log"
	"os"
	"os/signal"
	"time"
)

func ScheduledJob(ctx context.Context) {
	log.Println("ScheduledJob Started")

	//You can access foo value by using context
	log.Println("Context Value:" + ctx.Value("foo").(string))
	time.Sleep(1 * time.Second)

	log.Println("ScheduledJob Finished")
}

func main() {
	//Start Scheduler
	scheduler := worker.NewScheduler()

	//Initialize Context
	ctx := context.Background()

	//Add Value to the context if you want
	c1 := context.WithValue(ctx, "foo", "bar")

	// Initialize Scheduler
	// First: Context Second: Function Third: Interval Fourth: active status
	// Returns trigger channel and activation channel.
	// If you don't want to start immediately pass active status as false and use active channel
	trigger, active := scheduler.Add(c1, ScheduledJob, time.Second*10, true)

	// Soft Start/Stop
	/*
		time.AfterFunc(10*time.Second, func() {
			log.Println("ScheduledJob Starting")
			active <- true // or false for stopping the scheduler
		})
	*/

	// If you need to trigger somehow use trigger channel
	//Specific Trigger
	/*
		time.AfterFunc(20*time.Second, func() {
			log.Println("ScheduledJob Triggered")
			trigger <- true
		})
	*/

	//Waiting
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//Stop Scheduler
	scheduler.Stop()

	//close scheduler channels (it's not required)
	close(trigger)
	close(active)
}

```