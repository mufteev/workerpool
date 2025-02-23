package workerpool

import (
	"errors"
	"sync/atomic"
)

type Task func()

type Pool struct {
	collector chan Task
	isStopped atomic.Bool

	Tasks   []*Task
	Workers []*worker

	workerCount int
}

var (
	ErrTaskNil         = errors.New("task is nil")
	ErrPoolStopped     = errors.New("pool is stopped")
	ErrInvalidPoolSize = errors.New("workerCount and collectorCount must be greater than 0")
)

func NewPool(workerCount, collectorCount int) (*Pool, error) {
	if workerCount <= 0 || collectorCount <= 0 {
		return nil, ErrInvalidPoolSize
	}

	return &Pool{
		isStopped:   atomic.Bool{},
		workerCount: workerCount,
		Workers:     make([]*worker, workerCount),
		collector:   make(chan Task, collectorCount),
	}, nil
}

func (p *Pool) AddTask(t Task) error {
	if t == nil {
		return ErrTaskNil
	}

	if p.isStopped.Load() {
		return ErrPoolStopped
	}

	p.collector <- t

	return nil
}

func (p *Pool) RunBackground() {
	for i := range p.workerCount {
		w := newWorker(p.collector)
		p.Workers[i] = w

		go w.startBackground()
	}
}

func (p *Pool) Stop() {
	p.isStopped.Store(true)

	for i := range p.workerCount {
		p.Workers[i].stop()
	}
}
