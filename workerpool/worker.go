package workerpool

import "gitmic/semaphore"

type worker struct {
	n int

	// workOnDuring <-chan struct{}
	sem      *semaphore.Semaphore
	taskChan chan *Task
	quitChan chan bool
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
