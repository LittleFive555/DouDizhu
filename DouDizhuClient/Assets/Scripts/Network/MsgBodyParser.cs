using System.Collections.Generic;
using Google.Protobuf;
using Network.Proto;

namespace Network
{
    public class MsgBodyParser
    {
        private static Dictionary<PMsgId, MessageParser> m_MsgParserInfos = new()
        {
            { PMsgId.Handshake, PHandshakeResponse.Parser },
            { PMsgId.Login, PLoginResponse.Parser },
            { PMsgId.Register, PEmptyResponse.Parser },
            { PMsgId.ChatMsg, PChatMsgNotification.Parser },
        };

        private static Dictionary<PMsgId, MessageParser> m_NotifyParsers = new()
        {
            { PMsgId.ChatMsg, PChatMsgNotification.Parser },
        };
        
        public static MessageParser GetResponseParser(PMsgId msgId)
        {
            if (m_MsgParserInfos.TryGetValue(msgId, out var parser))
                return parser;
            return null;
        }
        
        public static MessageParser GetNotifyParser(PMsgId msgId)
        {
            if (m_NotifyParsers.TryGetValue(msgId, out var parser))
                return parser;
            return null;
        }
    }
}