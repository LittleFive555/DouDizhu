using System;
using Google.Protobuf;

namespace Network
{

    public partial class NetworkManager
    {
        private class NotificationHandlerAdapter<TNotification> where TNotification : IMessage
        {
            private readonly Action<TNotification> m_Handler;

            public NotificationHandlerAdapter(Action<TNotification> handler)
            {
                m_Handler = handler;
            }

            public void Handle(IMessage message)
            {
                if (message is TNotification typedMessage)
                    m_Handler(typedMessage);
                else
                    throw new Exception($"消息类型不匹配: {message.GetType()}");
            }
        }
    }
}