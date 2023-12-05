package semaphore

type semaphore struct {
	sem chan struct{}
}

func (s *semaphore) Acquire() {
	s.sem <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.sem
}

func NewSemaphore(tickets int) *semaphore {
	return &semaphore{
		sem: make(chan struct{}, tickets),
	}
}
