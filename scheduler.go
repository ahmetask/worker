package worker

import (
	"context"
	"sync"
	"time"
)

type Job func(ctx context.Context)

type Scheduler struct {
	wg            *sync.WaitGroup
	cancellations []context.CancelFunc
	triggers      []chan bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
		triggers:      make([]chan bool, 0),
	}
}

func (s *Scheduler) Add(ctx context.Context, j Job, interval time.Duration) chan bool {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations = append(s.cancellations, cancel)

	ch := make(chan bool)

	s.triggers = append(s.triggers, ch)
	s.wg.Add(1)
	go s.process(ctx, j, interval, ch)
	return ch
}

func (s *Scheduler) TriggerAll() {
	for _, ch := range s.triggers {
		ch <- true
	}
}

func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellations {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) process(ctx context.Context, j Job, interval time.Duration, trigger chan bool) {
	ticker := time.NewTicker(interval)
	first := make(chan bool, 1)

	first <- true
	for {
		select {
		case <-first:
			j(ctx)
		case <-ticker.C:
			j(ctx)
		case <-trigger:
			j(ctx)
		case <-ctx.Done():
			s.wg.Done()
			ticker.Stop()
			close(first)
			return
		}
	}
}
