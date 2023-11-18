package workerpool

type Pool struct {
	collector chan *Task

	Tasks   []*Task
	Workers []*Worker

	workerCount int
}

func NewPool(workerCount, collectorCount int) *Pool {
	return &Pool{
		workerCount: workerCount,
		Workers:     make([]*Worker, workerCount),
		collector:   make(chan *Task, collectorCount),
	}
}

func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

func (p *Pool) RunBackgground() {
	for i := 0; i < p.workerCount; i++ {
		worker := NewWorker(p.collector, i)
		p.Workers[i] = worker

		go worker.StartBackground()
	}
}

func (p *Pool) Stop() {
	for i := 0; i < p.workerCount; i++ {
		p.Workers[i].Stop()
	}
}
