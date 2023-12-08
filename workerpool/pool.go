package workerpool

import (
	"gitmic/semaphore"
)

type Pool struct {
	// workOnDuring  chan struct{}
	collector   chan *Task
	workers     []*worker
	workerCount int

	sem *semaphore.Semaphore
	// timesOnDuring int
	// timeWork      time.Duration
	// timeSleep time.Duration
	// collectorCount int
}

func NewPool(workerCount, collectorCount int, sem *semaphore.Semaphore) *Pool {
	return &Pool{
		// timeWork:      timeWork,
		// timeSleep:     timeSleep,
		// timesOnDuring: timesOnDuring,
		// workOnDuring:  make(chan struct{}, timesOnDuring),
		sem:         sem,
		workerCount: workerCount,
		workers:     make([]*worker, workerCount),
		collector:   make(chan *Task, collectorCount),
	}
}

// func acceptor(p *Pool, timesOnDuring <-chan struct{}) {
// 	mx := sync.Mutex{}
// 	// tickerWork := time.NewTicker(p.timeWork)
// 	tickerCount := p.timesOnDuring

// 	for {

// 		if tickerCount == 0 {
// 			fmt.Printf("lock - %s\n", p.timeSleep)

// 			// tickerWork.Stop()
// 			mx.Lock()
// 			time.AfterFunc(p.timeSleep, func() {
// 				tickerCount = p.timesOnDuring

// 				// tickerWork.Reset(p.timeWork)
// 				mx.Unlock()
// 			})
// 		}

// 		mx.Lock()
// 		tickerCount--
// 		p.workOnDuring <- struct{}{}
// 		mx.Unlock()
// 		// select {
// 		// case <-tickerWork.C:
// 		// 	fmt.Println("sleep")
// 		// 	time.Sleep(p.timeSleep)
// 		// default:
// 		// }
// 	}
// }

func (p *Pool) AddTask(t *Task) {
	p.sem.Acquire()
	p.collector <- t
}

func (p *Pool) RunBackground() {
	for i := 0; i < p.workerCount; i++ {
		w := newWorker(i, p.sem, p.collector)
		p.workers[i] = w

		go w.startBackground()
	}

	// go acceptor(p, p.workOnDuring)
}

func (p *Pool) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.workers[i].stop()
	}
}
