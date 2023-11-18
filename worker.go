package workerpool

type worker struct {
	taskChan chan *task
	quitChan chan bool
	Id       int
}

func newWorker(ch chan *task, Id int) *worker {
	return &worker{
		Id:       Id,
		taskChan: ch,
		quitChan: make(chan bool),
	}
}

func (w *worker) startBackground() {
	for {
		select {
		case task := <-w.taskChan:
			process(w.Id, task)
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
