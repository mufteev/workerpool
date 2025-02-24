package workerpool

type worker struct {
	taskChan chan Task
	quitChan chan bool
}

func newWorker(ch chan Task) *worker {
	return &worker{
		taskChan: ch,
		quitChan: make(chan bool),
	}
}

func (w *worker) startBackground() {
	for {
		select {
		case task := <-w.taskChan:
			task()
		case <-w.quitChan:
			return
		}
	}
}

func (w *worker) stop() {
	w.quitChan <- true
}
