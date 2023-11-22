package workerpool

type pool struct {
	collector chan *Task

	Tasks   []*Task
	Workers []*worker

	workerCount int
}

func NewPool(workerCount, collectorCount int) *pool {
	return &pool{
		workerCount: workerCount,
		Workers:     make([]*worker, workerCount),
		collector:   make(chan *Task, collectorCount),
	}
}

func (p *pool) AddTask(t *Task) {
	p.collector <- t
}

func (p *pool) RunBackgground() {
	for i := 0; i < p.workerCount; i++ {
		w := newWorker(p.collector, i)
		p.Workers[i] = w

		go w.startBackground()
	}
}

func (p *pool) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.Workers[i].stop()
	}
}
