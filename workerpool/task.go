package workerpool

import (
	"context"
	"errors"
	"time"
)

type Task struct {
	Err chan error
	Res chan interface{}

	f func() (interface{}, error)

	timeout time.Duration
}

var (
	defaultDuration = time.Second * time.Duration(20)
	errTimeoutTask  = errors.New("error timeout task from worker")
)

func NewTask(f func() (interface{}, error), timeout *time.Duration) *Task {

	if timeout == nil {
		timeout = &defaultDuration
	}

	return &Task{
		f:       f,
		timeout: *timeout,
		Res:     make(chan interface{}, 2),
		Err:     make(chan error, 2),
	}
}

type response struct {
	Res interface{}
	Err error
}

func process(t *Task) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	responseCh := make(chan response)

	go func() {
		res, err := t.f()
		responseCh <- response{res, err}
	}()

	select {
	case <-ctx.Done():
		t.Err <- errTimeoutTask

	case res := <-responseCh:
		if res.Err != nil {
			t.Err <- res.Err
		} else {
			t.Res <- res.Res
		}
	}

	close(t.Err)
	close(t.Res)
}
