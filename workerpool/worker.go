package workerpool

type worker struct {
	n int

	workOnDuring <-chan struct{}
	taskChan     chan *Task
	quitChan     chan bool
}

func newWorker(n int, timesOnDuring <-chan struct{}, ch chan *Task) *worker {
	return &worker{
		n:            n,
		taskChan:     ch,
		quitChan:     make(chan bool),
		workOnDuring: timesOnDuring,
	}
}

func (w *worker) startBackground() {
	for {
		<-w.workOnDuring
		select {
		case <-w.quitChan:
			return

		case task := <-w.taskChan:
			process(task)
		}
	}
}

func (w *worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}
