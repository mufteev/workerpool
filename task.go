package workerpool

type Task struct {
	f func()
}

func NewTask(f func()) *Task {
	return &Task{
		f: f,
	}
}
