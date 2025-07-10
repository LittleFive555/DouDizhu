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

            // 房间相关
            { PMsgId.CreateRoom, PCreateRoomResponse.Parser },
            { PMsgId.GetRoomList, PGetRoomListResponse.Parser },
            { PMsgId.EnterRoom, PEnterRoomResponse.Parser },
            { PMsgId.LeaveRoom, PEmptyResponse.Parser },

            // 聊天相关
            { PMsgId.ChatMsg, PChatMsgNotification.Parser },
        };

        private static Dictionary<PMsgId, MessageParser> m_NotifyParsers = new()
        {
            // 聊天相关
            { PMsgId.ChatMsg, PChatMsgNotification.Parser },

            // 房间相关
            { PMsgId.RoomChanged, PRoomChangedNotification.Parser },
            { PMsgId.RoomDisbanded, PRoomDisbandedNotification.Parser },
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