package message

// 消息工作协程
type MessageWorker struct {
	queue   <-chan *Message
	handler func(message *Message)
}

func NewMessageWorker(queue <-chan *Message, handler func(message *Message)) *MessageWorker {
	return &MessageWorker{
		queue:   queue,
		handler: handler,
	}
}

func (w *MessageWorker) Run() {
	for msg := range w.queue {
		w.handler(msg)
	}
}
