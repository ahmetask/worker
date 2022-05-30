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
