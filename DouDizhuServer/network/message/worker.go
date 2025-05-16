package message

type Worker struct {
	queue   <-chan *Message
	handler func(message *Message) error
}

func NewWorker(queue <-chan *Message, handler func(message *Message) error) *Worker {
	return &Worker{
		queue:   queue,
		handler: handler,
	}
}

func (w *Worker) Run() {
	for msg := range w.queue {
		w.handler(msg)
	}
}
