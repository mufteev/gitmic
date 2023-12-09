package semaphore

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	quitChan chan struct{}

	mxAccept       *sync.Mutex
	mxTicket       *sync.Mutex
	ticketReleaser *time.Ticker

	t1             time.Duration // T1
	t2             time.Duration // T2
	baseTickets    int           // bN
	currentTickets int           // cN
}

func (s *Semaphore) Acquire() {
	s.mxTicket.Lock()
	s.currentTickets++
	fmt.Printf("currentTickets - %d\n", s.currentTickets)

	if s.currentTickets == s.baseTickets {
		fmt.Printf("sleep - %d\n", s.currentTickets)
		s.mxAccept.Lock()
		s.ticketReleaser.Stop()
		time.AfterFunc(s.t2, func() {
			s.currentTickets = 0
			s.ticketReleaser.Reset(s.t1)
			s.mxAccept.Unlock()
		})
	}
	s.mxTicket.Unlock()

	s.mxAccept.Lock()
	defer s.mxAccept.Unlock()
}

func (s *Semaphore) Release() {
	s.mxTicket.Lock()
	defer s.mxTicket.Unlock()
	if s.currentTickets > 0 {
		s.currentTickets--
	}
}

func NewSemaphore(tickets int, t1, t2 time.Duration) *Semaphore {
	return &Semaphore{
		currentTickets: 0,
		t1:             t1,
		t2:             t2,
		baseTickets:    tickets,
		mxAccept:       &sync.Mutex{},
		mxTicket:       &sync.Mutex{},
		ticketReleaser: time.NewTicker(t1),
		quitChan:       make(chan struct{}),
	}
}

func (s *Semaphore) Run() {
	for {
		select {
		case <-s.quitChan:
			return

		case <-s.ticketReleaser.C:
			s.mxTicket.Lock()
			s.currentTickets = 0
			s.mxTicket.Unlock()
		}
	}
}

func (s *Semaphore) Stop() {
	go func() {
		s.quitChan <- struct{}{}
		close(s.quitChan)
	}()
}
