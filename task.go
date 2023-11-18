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
	defaultTimeout = time.Millisecond * time.Duration(500)
	errTaskTimeout = errors.New("error timeout task from worker")
)

func NewTask(f func() (interface{}, error), timeout *time.Duration) *Task {
	if timeout == nil {
		timeout = &defaultTimeout
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

func process(workerId int, task *Task) {
	ctx, cancel := context.WithTimeout(context.Background(), task.timeout)
	defer cancel()

	responseCh := make(chan response)

	go func() {
		res, err := task.f()
		responseCh <- response{res, err}
	}()

	select {
	case res := <-responseCh:
		if res.Err != nil {
			task.Err <- res.Err
		} else {
			task.Res <- res.Res
		}
	case <-ctx.Done():
		task.Err <- errTaskTimeout
	}

	close(task.Err)
	close(task.Res)
}
