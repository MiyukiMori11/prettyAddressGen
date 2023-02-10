package workerPool

import (
	"errors"
	"fmt"
	"sync"
)

type TaskFunc func() interface{}

type workerPool struct {
	workers map[*worker]struct{}
	out     chan interface{}
}

type WorkerPool interface {
	Run(in <-chan TaskFunc, workers int64) (<-chan interface{}, error)
	CloseAll()
}

func New() WorkerPool {
	return &workerPool{
		workers: make(map[*worker]struct{}),
		out:     make(chan interface{}),
	}
}

func (w *workerPool) Run(in <-chan TaskFunc, workers int64) (<-chan interface{}, error) {

	if err := w.addWorkers(workers, in); err != nil {
		return nil, fmt.Errorf("can't add workers: %w", err)
	}

	go w.healthCheck()

	for wk := range w.workers {
		go func(wk *worker) {
			wk.Run()
		}(wk)
	}

	return w.out, nil
}

func (w *workerPool) healthCheck() {
	for len(w.workers) > 0 {
		for wk := range w.workers {
			if wk.IsDone() {
				delete(w.workers, wk)
			}
		}
	}

	close(w.out)
}

func (w *workerPool) addWorkers(count int64, in <-chan TaskFunc) error {
	if count < 1 {
		return errors.New("workers count can't be less 1")
	}
	for i := 0; i < int(count); i++ {
		w.workers[&worker{
			ID:   i,
			Task: in,
			Out:  w.out,
			Quit: make(chan struct{}),
		}] = struct{}{}

	}

	return nil
}

func (w *workerPool) CloseAll() {
	var wg sync.WaitGroup

	for wk := range w.workers {
		wg.Add(1)
		go func(wk *worker) {
			defer wg.Done()
			wk.Stop()
		}(wk)
	}

	wg.Wait()

}
