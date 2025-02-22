package workerpool

type worker struct {
	taskChan chan *Task
	quitChan chan bool
}

func newWorker(ch chan *Task) *worker {
	return &worker{
		taskChan: ch,
		quitChan: make(chan bool),
	}
}

func (w *worker) startBackground() {
	for {
		select {
		case task := <-w.taskChan:
			task.f()
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
