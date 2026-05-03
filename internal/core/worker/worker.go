package worker

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Run() {
	messageWorker := NewMessageWorker()
	clientWorker := NewClientWorker()

	messageWorker.Run()
	clientWorker.Run()

	select {}
}

// func RegisterWorker() {
// }
