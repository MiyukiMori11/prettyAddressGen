package workerPool

import (
	"fmt"
)

type worker struct {
	ID   int
	Task <-chan TaskFunc
	Out  chan<- interface{}
	Quit chan struct{}
	done bool
}

func (wk *worker) Run() {
	fmt.Printf("worker %d started\n", wk.ID)
loop:
	for {
		select {
		case <-wk.Quit:
			wk.makeDone()
			break loop
		case task, ok := <-wk.Task:
			if ok {
				wk.Out <- task()
			} else {
				wk.makeDone()
				break loop
			}

		}
	}

	fmt.Printf("worker %d finished\n", wk.ID)
}

func (wk *worker) Stop() {
	wk.Quit <- struct{}{}

}

func (wk *worker) IsDone() bool {
	return wk.done
}

func (wk *worker) makeDone() {
	wk.done = true
}
