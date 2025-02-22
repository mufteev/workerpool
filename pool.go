package workerpool

type Pool struct {
	collector chan *Task

	Tasks   []*Task
	Workers []*worker

	workerCount int
}

func NewPool(workerCount, collectorCount int) *Pool {
	return &Pool{
		workerCount: workerCount,
		Workers:     make([]*worker, workerCount),
		collector:   make(chan *Task, collectorCount),
	}
}

func (p *Pool) AddTask(t *Task) {
	p.collector <- t
}

func (p *Pool) RunBackgground() {
	for i := 0; i < p.workerCount; i++ {
		w := newWorker(p.collector, i)
		p.Workers[i] = w

		go w.startBackground()
	}
}

func (p *Pool) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.Workers[i].stop()
	}
}
