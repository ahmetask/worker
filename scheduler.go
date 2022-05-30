package worker

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

// Job scheduled job func
type Job func(ctx context.Context)

// Scheduler instance
type Scheduler struct {
	wg            *sync.WaitGroup
	cancellations []context.CancelFunc
	triggers      []chan bool
}

// NewScheduler creates Scheduler instance
func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
		triggers:      make([]chan bool, 0),
	}
}

// Add scheduled job
func (s *Scheduler) Add(ctx context.Context, j Job,
	interval time.Duration, active bool) (triggerCh, activeCh chan bool) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations = append(s.cancellations, cancel)

	triggerChannel := make(chan bool)

	activeChannel := make(chan bool)

	s.triggers = append(s.triggers, triggerChannel)
	s.wg.Add(1)
	go s.process(ctx, j, interval, triggerChannel, activeChannel, active)
	return triggerChannel, activeChannel
}

// TriggerAll trigger all scheduled jobs
func (s *Scheduler) TriggerAll() {
	for _, ch := range s.triggers {
		ch <- true
	}
}

// Stop scheduled jobs
func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellations {
		cancel()
	}
	s.wg.Wait()
}

// Run Scheduler
func (s *Scheduler) run(ctx context.Context, j Job, isActive bool) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic scheduled job: %v\n%s\n", r, buf)
		}
	}()
	if isActive {
		j(ctx)
	}
}

func (s *Scheduler) process(ctx context.Context, j Job,
	interval time.Duration, trigger chan bool, activeCh chan bool, active bool) {
	ticker := time.NewTicker(interval)
	first := make(chan bool, 1)
	first <- true
	isActive := active

	for {
		select {
		case a := <-activeCh:
			isActive = a
		case <-first:
			s.run(ctx, j, isActive)
		case <-ticker.C:
			s.run(ctx, j, isActive)
		case <-trigger:
			s.run(ctx, j, isActive)
			<-ticker.C
		case <-ctx.Done():
			s.wg.Done()
			ticker.Stop()
			close(first)
			return
		}
	}
}
