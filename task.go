package workerpool

import (
	"context"
	"errors"
	"time"
)

type task struct {
	Err chan error
	Res chan interface{}

	f func() (interface{}, error)

	timeout time.Duration
}

var (
	defaultTimeout = time.Millisecond * time.Duration(500)
	errTaskTimeout = errors.New("error timeout task from worker")
)

func NewTask(f func() (interface{}, error), timeout *time.Duration) *task {
	if timeout == nil {
		timeout = &defaultTimeout
	}

	return &task{
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

func process(workerId int, t *task) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	responseCh := make(chan response)

	go func() {
		res, err := t.f()
		responseCh <- response{res, err}
	}()

	select {
	case res := <-responseCh:
		if res.Err != nil {
			t.Err <- res.Err
		} else {
			t.Res <- res.Res
		}
	case <-ctx.Done():
		t.Err <- errTaskTimeout
	}

	close(t.Err)
	close(t.Res)
}
