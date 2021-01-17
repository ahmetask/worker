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
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
	}
}

func (s *Scheduler) Add(ctx context.Context, j Job, interval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations = append(s.cancellations, cancel)

	s.wg.Add(1)
	go s.process(ctx, j, interval)
}

func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellations {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) process(ctx context.Context, j Job, interval time.Duration) {
	ticker := time.NewTicker(interval)
	first := make(chan bool, 1)
	first <- true
	for {
		select {
		case <-ticker.C:
			j(ctx)
		case <-first:
			j(ctx)
		case <-ctx.Done():
			s.wg.Done()
			ticker.Stop()
			close(first)
			return
		}
	}
}
