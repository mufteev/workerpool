package workerpool

type worker struct {
	taskChan chan *Task
	quitChan chan bool
	Id       int
}

func newWorker(ch chan *Task, Id int) *worker {
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
