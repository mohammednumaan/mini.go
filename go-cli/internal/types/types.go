package types

type Task struct {
	Title     string
	Completed bool
}

func NewTask(title string) *Task {
	return &Task{
		Title:     title,
		Completed: false,
	}
}
