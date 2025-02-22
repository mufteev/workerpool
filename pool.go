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

func (p *Pool) RunBackground() {
	for i := range p.workerCount {
		w := newWorker(p.collector)
		p.Workers[i] = w

		go w.startBackground()
	}
}

func (p *Pool) Stop() {
	for i := range p.workerCount {
		p.Workers[i].stop()
	}
}
