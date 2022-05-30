package worker

import (
	"sync"
	"time"
)

// Pool worker pool
type Pool struct {
	singleJob     chan Work
	internalQueue chan Work
	/*
		boss says hey I have a new job at my desk.
		workers who available can get it in this way he does not have to ask current status of workers
	*/
	readyPool         chan chan Work
	workers           []*worker
	dispatcherStopped sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

// NewWorkerPool creates worker pool
func NewWorkerPool(maxWorkers, jobQueueCapacity int) *Pool {
	if jobQueueCapacity <= 0 {
		jobQueueCapacity = 100
	}

	workersStopped := sync.WaitGroup{}

	readyPool := make(chan chan Work, maxWorkers)
	workers := make([]*worker, maxWorkers)

	// create workers
	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(i+1, readyPool, &workersStopped)
	}

	return &Pool{
		internalQueue:     make(chan Work, jobQueueCapacity),
		singleJob:         make(chan Work),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: sync.WaitGroup{},
		workersStopped:    &workersStopped,
		quit:              make(chan bool),
	}
}

// Start worker pool
func (q *Pool) Start() {
	// tell workers to get ready
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].start()
	}
	// open factory
	go q.dispatch()
}

// Stop pool
func (q *Pool) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *Pool) dispatch() {
	// open factory gate
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.singleJob:
			workerXChannel := <-q.readyPool // if ree worker x founded
			workerXChannel <- job           // here is your job worker x
		case job := <-q.internalQueue:
			workerXChannel := <-q.readyPool // free worker x founded
			workerXChannel <- job           // here is your job worker x
		case <-q.quit:
			// free all workers
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].stop()
			}
			// wait for all workers to finish their job
			q.workersStopped.Wait()
			// close factory gate
			q.dispatcherStopped.Done()
			return
		}
	}
}

// Submit This is blocking if all workers are busy*/
func (q *Pool) Submit(job Work) {
	// daily - fill the board with new works
	q.singleJob <- job
}

// Enqueue Tries to enqueue but fails if queue is full*/
func (q *Pool) Enqueue(job Work) bool {
	select {
	case q.internalQueue <- job:
		return true
	default:
		return false
	}
}

// EnqueueWithTimeout try to enqueue but fails if timeout occurs
func (q *Pool) EnqueueWithTimeout(job Work, timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	ch := make(chan bool, 1)
	t := time.AfterFunc(timeout, func() {
		ch <- false
	})
	defer func() {
		t.Stop()
	}()

	select {
	case q.internalQueue <- job:
		return true
	case <-ch:
		return false
	}
}
