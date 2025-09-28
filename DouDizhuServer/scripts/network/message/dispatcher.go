package message

type IMessageDispatcher interface {
	EnqueueMessage(message *Message)
}

type INotificationDispatcher interface {
	EnqueueNotification(notification *Notification)
}
