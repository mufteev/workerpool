package workerpool

import (
	"context"
	"errors"
	"time"
)

type Task[T any] struct {
	Err chan error
	Res chan T

	f func() (T, error)

	timeout time.Duration
}

var (
	defaultTimeout = time.Millisecond * time.Duration(500)
	errTaskTimeout = errors.New("error timeout task from worker")
)

func NewTask[T any](f func() (T, error), timeout *time.Duration) *Task[T] {
	if timeout == nil {
		timeout = &defaultTimeout
	}

	return &Task[T]{
		f:       f,
		timeout: *timeout,
		Res:     make(chan T, 2),
		Err:     make(chan error, 2),
	}
}

type response[T any] struct {
	Res T
	Err error
}

func process[T any](workerId int, t *Task[T]) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	responseCh := make(chan response[T])

	go func() {
		res, err := t.f()
		responseCh <- response[T]{res, err}
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
