package workerpool

type Pool[T any] struct {
	collector chan *Task[T]

	Tasks   []*Task[T]
	Workers []*worker[T]

	workerCount int
}

func NewPool[T any](workerCount, collectorCount int) *Pool[T] {
	return &Pool[T]{
		workerCount: workerCount,
		Workers:     make([]*worker[T], workerCount),
		collector:   make(chan *Task[T], collectorCount),
	}
}

func (p *Pool[T]) AddTask(t *Task[T]) {
	p.collector <- t
}

func (p *Pool[T]) RunBackground() {
	for i := 0; i < p.workerCount; i++ {
		w := newWorker(p.collector, i)
		p.Workers[i] = w

		go w.startBackground()
	}
}

func (p *Pool[T]) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.Workers[i].stop()
	}
}
