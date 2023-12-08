package workerpool

import (
	"gitmic/semaphore"
)

type Pool struct {
	sem *semaphore.Semaphore

	collector   chan *Task
	workers     []*worker
	workerCount int
}

func NewPool(workerCount, collectorCount int, sem *semaphore.Semaphore) *Pool {
	return &Pool{
		sem:         sem,
		workerCount: workerCount,
		workers:     make([]*worker, workerCount),
		collector:   make(chan *Task, collectorCount),
	}
}

func (p *Pool) AddTask(t *Task) {
	p.sem.Acquire()
	p.collector <- t
}

func (p *Pool) RunBackground() {
	go p.sem.Run()

	for i := 0; i < p.workerCount; i++ {
		w := newWorker(i, p.sem, p.collector)
		p.workers[i] = w

		go w.startBackground()
	}
}

func (p *Pool) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.workers[i].stop()
	}
}
