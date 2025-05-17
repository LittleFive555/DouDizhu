using System;
using System.Net.Sockets;
using UnityEngine;
using System.Threading.Tasks;
using System.Collections.Concurrent;
using System.Threading;
using System.Collections.Generic;
using Google.Protobuf;
using Network.Tcp;
using Network.Proto;

namespace Network
{
    public partial class NetworkManager
    {
        private static NetworkManager m_Instance;
        public static NetworkManager Instance
        {
            get
            {
                if (m_Instance == null)
                    m_Instance = new NetworkManager();
                return m_Instance;
            }
        }

        private TcpClient m_TcpClient;
        private NetworkStream m_NetworkStream;
        private MessageReadWriter m_MessageReadWriter;
        private bool m_IsConnected = false;
        public bool IsConnected => m_IsConnected;

        
        private const int SERVER_PORT = 8080;


        private const int REQUEST_TIMEOUT = 5;

        private readonly ConcurrentDictionary<long, TaskCompletionSource<PGameMsgRespPacket>> m_PendingRequests = new();
        private readonly Dictionary<PGameNotificationPacket.ContentOneofCase, Action<IMessage>> m_NotificationHandlers = new();

        public NetworkManager()
        {
            m_MessageReadWriter = new LengthPrefixReadWriter();
        }

        public async Task ConnectAsync(string serverHost)
        {
            await ConnectAsync(serverHost, SERVER_PORT);
        }

        public void Disconnect()
        {
            if (m_NetworkStream != null)
            {
                m_NetworkStream.Close();
                m_NetworkStream = null;
            }
            if (m_TcpClient != null)
            {
                m_TcpClient.Close();
                m_TcpClient = null;
            }
            m_IsConnected = false;
            Debug.Log("TCP连接已关闭");
        }

        public async Task<NetworkResult<PCommonResponse>> RequestAsync<TReq>(PGameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            return await RequestAsync<TReq, PCommonResponse>(requestType, request);
        }

        public async Task<NetworkResult<TResp>> RequestAsync<TReq, TResp>(PGameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage where TResp : IMessage
        {
            if (!m_IsConnected || m_TcpClient == null || !m_TcpClient.Connected)
                return NetworkResult<TResp>.Failure("TCP连接未建立，无法发送消息");

            try
            {
                PGameClientMessage gameClientMessage = PackRequest(requestType, request);
                var tcs = new TaskCompletionSource<PGameMsgRespPacket>();
                m_PendingRequests.TryAdd(gameClientMessage.Header.MessageId, tcs);
                
                await m_MessageReadWriter.WriteTo(m_NetworkStream, gameClientMessage.ToByteArray());
                Debug.Log("消息已发送: " + request.ToString());

                // 设置超时
                using var timeoutCts = new CancellationTokenSource(TimeSpan.FromSeconds(REQUEST_TIMEOUT));
                timeoutCts.Token.Register(() => tcs.TrySetCanceled(), useSynchronizationContext: false);
                var serverPacket = await tcs.Task;
                if (serverPacket.ContentCase == PGameMsgRespPacket.ContentOneofCase.Error)
                    return NetworkResult<TResp>.Failure($"服务器返回错误: ErrorCode {serverPacket.Error.Code}: {serverPacket.Error.Message}");
                
                TResp response = UnpackResponse<TResp>(serverPacket);
                return NetworkResult<TResp>.Success(response);
            }
            catch (Exception ex)
            {
                return NetworkResult<TResp>.Failure($"发送消息时发生错误: {ex.Message}");
            }
        }

        public void RegisterNotificationHandler<TNotification>(PGameNotificationPacket.ContentOneofCase notificationType, Action<TNotification> handler) where TNotification : IMessage
        {
            var adapter = new NotificationHandlerAdapter<TNotification>(handler);
            if (m_NotificationHandlers.ContainsKey(notificationType))
                m_NotificationHandlers[notificationType] += adapter.Handle;
            else
                m_NotificationHandlers[notificationType] = adapter.Handle;
        }

        public void UnregisterNotificationHandler<TNotification>(PGameNotificationPacket.ContentOneofCase notificationType, Action<TNotification> handler) where TNotification : IMessage
        {
            // 查找对应的适配器
            var adapter = new NotificationHandlerAdapter<TNotification>(handler);
            if (m_NotificationHandlers.ContainsKey(notificationType))
            {
                m_NotificationHandlers[notificationType] -= adapter.Handle;
                if (m_NotificationHandlers[notificationType] == null)
                    m_NotificationHandlers.Remove(notificationType);
            }
        }

        private async Task ConnectAsync(string serverHost, int port)
        {
            try
            {
                m_TcpClient = new TcpClient();
                await m_TcpClient.ConnectAsync(serverHost, port);
                m_NetworkStream = m_TcpClient.GetStream();
                m_IsConnected = true;
                Debug.Log("TCP连接已建立");
                _ = ReceiveLoopAsync();
            }
            catch (SocketException ex)
            {
                Debug.LogError("TCP连接失败: " + ex.Message);
                m_IsConnected = false;
            }
        }

        private async Task ReceiveLoopAsync()
        {
            try
            {
                while (m_IsConnected)
                {
                    byte[] responseBuffer = await m_MessageReadWriter.ReadFrom(m_NetworkStream);
                    if (responseBuffer == null || responseBuffer.Length == 0)
                    {
                        Debug.LogError("服务器正常断开");
                        break;
                    }

                    PGameServerMessage gameServerMessage = PGameServerMessage.Parser.ParseFrom(responseBuffer);
                    if (gameServerMessage.ContentCase == PGameServerMessage.ContentOneofCase.Notification)
                        OnNotified(gameServerMessage.Notification);
                    else if (gameServerMessage.ContentCase == PGameServerMessage.ContentOneofCase.Response)
                        OnResponse(gameServerMessage.Response);
                }
            }
            catch (Exception ex)
            {
                Debug.LogError("接收消息时发生错误: " + ex.Message);
            }
        }

        private void OnResponse(PGameMsgRespPacket gameMsgRespPacket)
        {
            if (m_PendingRequests.TryGetValue(gameMsgRespPacket.Header.MessageId, out var tcs))
            {
                tcs.SetResult(gameMsgRespPacket);
                m_PendingRequests.TryRemove(gameMsgRespPacket.Header.MessageId, out _);
            }
        }

        private void OnNotified(PGameNotificationPacket notification)
        {
            if (m_NotificationHandlers.TryGetValue(notification.ContentCase, out var handler))
            {
                try
                {
                    string propertyName = notification.ContentCase.ToString();
                    var value = notification.GetType().GetProperty(propertyName).GetValue(notification);
                    handler(value as IMessage);
                }
                catch (Exception ex)
                {
                    Debug.LogError($"通知处理程序发生错误: {ex.Message}");
                }
            }
        }

        private PGameClientMessage PackRequest<TReq>(PGameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            PGameClientMessage gameClientMessage = new PGameClientMessage()
            {
                Header = new PGameMsgHeader()
                { 
                    MessageId = GenerateRequestId(),
                    Player = new PPlayer()
                    {
                        PlayerId = "", // TODO
                        Nickname = ""
                    }
                },
            };
            
            string propertyName = requestType.ToString();
            Type packetType = typeof(PGameClientMessage);
            var property = packetType.GetProperty(propertyName);
            if (property != null)
                property.SetValue(gameClientMessage, request);
            else
                Debug.LogError($"未找到对应的属性: {propertyName}");
            return gameClientMessage;
        }

        private TResp UnpackResponse<TResp>(PGameMsgRespPacket response) where TResp : IMessage
        {
            string propertyName = response.ContentCase.ToString();
            
            Type packetType = typeof(PGameMsgRespPacket);
            var property = packetType.GetProperty(propertyName);
            if (property != null)
            {
                // 获取属性值并转换为目标类型
                var value = property.GetValue(response);
                if (value is TResp typedValue)
                    return typedValue;
            }
            
            Debug.LogError($"无法从响应中获取类型 {typeof(TResp).Name} 的数据");
            return default(TResp);
        }

        private static long GenerateRequestId() => DateTimeOffset.UtcNow.Ticks;
    }
}