package worker

import (
	"sync"
	"time"
)

type Pool struct {
	singleJob         chan Work
	internalQueue     chan Work
	readyPool         chan chan Work //boss says hey i have a new job at my desk workers who available can get it in this way he does not have to ask current status of workers
	workers           []*worker
	dispatcherStopped sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

func NewWorkerPool(opts ...opts) *Pool {

	cfg := buildWorkerPoolConfig(opts...)

	maxWorkers := cfg.maxWorkers
	jobQueueCapacity := cfg.jobQueueCapacity

	if jobQueueCapacity <= 0 {
		jobQueueCapacity = 100
	}

	workersStopped := sync.WaitGroup{}

	readyPool := make(chan chan Work, maxWorkers)
	workers := make([]*worker, maxWorkers)

	// create workers
	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(i+1, readyPool, &workersStopped)
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

func (q *Pool) Start() {
	//tell workers to get ready
	for _, w := range q.workers {
		w.Start()
	}
	// open factory
	go q.dispatch()
}

func (q *Pool) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *Pool) dispatch() {
	//open factory gate
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.singleJob:
			workerXChannel := <-q.readyPool //free worker x founded
			workerXChannel <- job           // here is your job worker x
		case job := <-q.internalQueue:
			workerXChannel := <-q.readyPool //free worker x founded
			workerXChannel <- job           // here is your job worker x
		case <-q.quit:
			// free all workers
			for _, w := range q.workers {
				w.Stop()
			}
			// wait for all workers to finish their job
			q.workersStopped.Wait()
			//close factory gate
			q.dispatcherStopped.Done()
			return
		}
	}
}

/*This is blocking if all workers are busy*/
func (q *Pool) Submit(job Work) {
	// daily - fill the board with new works
	q.singleJob <- job
}

/*Tries to enqueue but fails if queue is full*/
func (q *Pool) Enqueue(job Work) bool {
	select {
	case q.internalQueue <- job:
		return true
	default:
		return false
	}
}

/*try to enqueue but fails if timeout occurs*/
func (q *Pool) EnqueueWithTimeout(job Work, timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}

	ch := make(chan bool)
	t := time.AfterFunc(timeout, func() { ch <- false })
	defer func() {
		t.Stop()
		close(ch)
	}()

	for {
		select {
		case q.internalQueue <- job:
			return true
		case <-ch:
			return false
		}
	}
}
