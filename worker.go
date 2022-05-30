package worker

import (
	"log"
	"runtime"
	"sync"
)

type worker struct {
	id          int
	done        *sync.WaitGroup
	readyPool   chan chan Work // get work from the boss
	work        chan Work
	currentWork Work
	quit        chan bool
}

func newWorker(id int, readyPool chan chan Work, done *sync.WaitGroup) *worker {
	return &worker{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan Work),
		quit:      make(chan bool),
	}
}

func (w *worker) process(work Work) {
	// Do the work
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic running process: %v\n%s\n", r, buf)
		}
	}()
	work.Do()
}

func (w *worker) start() {
	w.done.Add(1) // wait for me
	go func() {
		for {
			w.readyPool <- w.work // hey I am ready to work on new job
			select {
			case work := <-w.work: // hey I am waiting for new job
				w.currentWork = work
				w.process(work) // ok I am on it
			case <-w.quit:
				w.done.Done() // ok I am here I finished my all jobs
				return
			}
		}
	}()
}

func (w *worker) stop() {
	// tell worker to stop after current process
	w.quit <- true

	// stop current work
	if w.currentWork != nil {
		w.currentWork.Stop()
	}
}
