package workerPool

import (
	"fmt"
)

type worker struct {
	ID   int
	Task <-chan taskFunc
	Out  chan interface{}
	Quit chan struct{}
}

func (wk *worker) Run() {
	fmt.Printf("worker %d started\n", wk.ID)
loop:
	for {
		select {
		case <-wk.Quit:
			break loop
		case task := <-wk.Task:
			wk.Out <- task()
		}
	}

	fmt.Printf("worker %d finished\n", wk.ID)
}

func (wk *worker) Stop() {
	go func() {
		wk.Quit <- struct{}{}
	}()

}
