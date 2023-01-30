package workerPool

import "errors"

type taskFunc func() interface{}

type workerPool struct {
	workers []worker
	out     chan interface{}
	in      chan taskFunc
}

type WorkerPool interface {
	AddWorkers(count int64) error
	RunAll() (chan<- taskFunc, <-chan interface{})
	CloseAll()
}

func New() WorkerPool {
	return &workerPool{
		out: make(chan interface{}),
		in:  make(chan taskFunc),
	}
}

func (w *workerPool) AddWorkers(count int64) error {
	if count < 1 {
		return errors.New("workers count can't be less 1")
	}
	for i := 0; i < int(count); i++ {
		w.workers = append(w.workers, worker{
			ID:   i,
			Task: w.in,
			Out:  w.out,
		})

	}

	return nil
}

func (w *workerPool) RunAll() (chan<- taskFunc, <-chan interface{}) {

	for _, wk := range w.workers {
		go func(w worker) {
			w.Run()
		}(wk)
	}

	return w.in, w.out
}

func (w *workerPool) CloseAll() {
	close(w.in)
	close(w.out)
	for _, wk := range w.workers {
		wk.Stop()
	}
}
