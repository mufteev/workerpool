package workerpool

type Worker struct {
	taskChan chan *Task
	quitChan chan bool
	Id       int
}

func NewWorker(ch chan *Task, Id int) *Worker {
	return &Worker{
		Id:       Id,
		taskChan: ch,
		quitChan: make(chan bool),
	}
}

func (w *Worker) StartBackground() {
	for {
		select {
		case task := <-w.taskChan:
			process(w.Id, task)
		case <-w.quitChan:
			return
		}
	}
}

func (w *Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}
