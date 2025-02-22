package workerpool

type worker[T any] struct {
	taskChan chan *Task[T]
	quitChan chan bool
	ID       int
}

func newWorker[T any](ch chan *Task[T], ID int) *worker[T] {
	return &worker[T]{
		ID:       ID,
		taskChan: ch,
		quitChan: make(chan bool),
	}
}

func (w *worker[T]) startBackground() {
	for {
		select {
		case task := <-w.taskChan:
			process(w.ID, task)
		case <-w.quitChan:
			return
		}
	}
}

func (w *worker[T]) stop() {
	go func() {
		w.quitChan <- true
	}()
}
