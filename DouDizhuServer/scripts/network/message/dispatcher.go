package message

type INotificationDispatcher interface {
	EnqueueNotification(notification *Notification)
}

type MessageDispatcher struct {
	messageQueue   chan *Message    // 消息队列
	messageWorkers []*MessageWorker // 工作协程池
}

// 初始化消息分发器
func NewMessageDispatcher(workerCount int, handler func(message *Message)) *MessageDispatcher {
	md := &MessageDispatcher{
		messageQueue: make(chan *Message, 10000), // 带缓冲的队列
	}

	// 初始化工作协程
	for i := 0; i < workerCount; i++ {
		worker := NewMessageWorker(md.messageQueue, handler)
		md.messageWorkers = append(md.messageWorkers, worker)
		go worker.Run()
	}

	return md
}

// 接收消息
func (md *MessageDispatcher) EnqueueMessage(message *Message) {
	md.messageQueue <- message
}
