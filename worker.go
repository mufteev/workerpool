package workerpool

type worker struct {
	taskChan chan *Task
	quitChan chan bool
	ID       int
}

func newWorker(ch chan *Task, ID int) *worker {
	return &worker{
		ID:       ID,
		taskChan: ch,
		quitChan: make(chan bool),
	}
}

func (w *worker) startBackground() {
	for {
		select {
		case task := <-w.taskChan:
			process(w.ID, task)
		case <-w.quitChan:
			return
		}
	}
}

func (w *worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}
