package workerpool

import "gitmic/semaphore"

type worker struct {
	sem *semaphore.Semaphore

	taskChan chan *Task
	quitChan chan bool

	n int
}

func newWorker(n int, sem *semaphore.Semaphore, ch chan *Task) *worker {
	return &worker{
		n:        n,
		taskChan: ch,
		sem:      sem,
		quitChan: make(chan bool),
	}
}

func (w *worker) startBackground() {
	for {
		select {
		case <-w.quitChan:
			return

		case task := <-w.taskChan:
			process(task)
			w.sem.Release()
		}
	}
}

func (w *worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}
