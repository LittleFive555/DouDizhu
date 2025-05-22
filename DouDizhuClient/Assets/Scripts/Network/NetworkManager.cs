using System;
using System.Net.Sockets;
using System.Threading.Tasks;
using System.Collections.Concurrent;
using System.Threading;
using System.Collections.Generic;
using Google.Protobuf;
using Network.Tcp;
using Network.Proto;
using Serilog;

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
            Log.Information("TCP连接已关闭");
        }

        public async Task<NetworkResult<PEmptyResponse>> RequestAsync<TReq>(PGameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            return await RequestAsync<TReq, PEmptyResponse>(requestType, request);
        }

        public async Task<NetworkResult<TResp>> RequestAsync<TReq, TResp>(PGameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage where TResp : class, IMessage
        {
            if (!m_IsConnected || m_TcpClient == null || !m_TcpClient.Connected)
            {
                Log.Error("TCP连接未建立，无法发送消息");
                return NetworkResult<TResp>.Error(null, "TCP连接未建立，无法发送消息");
            }

            try
            {
                PGameClientMessage gameClientMessage = PackRequest(requestType, request);
                var tcs = new TaskCompletionSource<PGameMsgRespPacket>();
                m_PendingRequests.TryAdd(gameClientMessage.Header.MessageId, tcs);
                
                // TODO 对消息进行加密
                await m_MessageReadWriter.WriteTo(m_NetworkStream, gameClientMessage.ToByteArray());
                Log.Information("消息已发送: {request}", request);

                // 设置超时
                using var timeoutCts = new CancellationTokenSource(TimeSpan.FromSeconds(REQUEST_TIMEOUT));
                timeoutCts.Token.Register(() => tcs.TrySetCanceled(), useSynchronizationContext: false);
                var serverPacket = await tcs.Task;
                if (serverPacket.ContentCase == PGameMsgRespPacket.ContentOneofCase.Error)
                {
                    if (serverPacket.Error.Type == PError.Types.Type.ServerError)
                        Log.Error("服务器内部错误: {errorCode}: {errorMessage}", serverPacket.Error.ErrorCode, serverPacket.Error.Message);
                    else
                        Log.Error("服务器返回游戏逻辑错误: {errorCode}: {errorMessage}", serverPacket.Error.ErrorCode, serverPacket.Error.Message);
                    return NetworkResult<TResp>.Error(serverPacket.Error.ErrorCode, serverPacket.Error.Message);
                }
                
                TResp response = UnpackResponse<TResp>(serverPacket);
                Log.Information("收到服务器响应: {response}", response);
                return NetworkResult<TResp>.Success(response);
            }
            catch (Exception ex)
            {
                Log.Error("发送消息时发生错误: {exception}", ex.Message);
                return NetworkResult<TResp>.Error(null, ex.Message);
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
                Log.Information("TCP连接已建立");
                _ = ReceiveLoopAsync();
            }
            catch (SocketException ex)
            {
                Log.Error("TCP连接失败: {exception}", ex.Message);
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
                        Log.Information("服务器正常断开");
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
                Log.Error("接收消息时发生错误: {exception}", ex.Message);
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
                    Log.Error("通知处理程序发生错误: {exception}", ex.Message);
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
                    PlayerId = "" // TODO
                },
            };
            
            string propertyName = requestType.ToString();
            Type packetType = typeof(PGameClientMessage);
            var property = packetType.GetProperty(propertyName);
            if (property == null)
                throw new Exception($"未找到对应的属性: {propertyName}");
                
            property.SetValue(gameClientMessage, request);
            return gameClientMessage;
        }

        private TResp UnpackResponse<TResp>(PGameMsgRespPacket response) where TResp : IMessage
        {
            string propertyName = response.ContentCase.ToString();
            
            Type packetType = typeof(PGameMsgRespPacket);
            var property = packetType.GetProperty(propertyName);
            if (property == null)
                throw new Exception($"响应 PGameMsgRespPacket 中未找到对应的属性: {propertyName}");
                
            // 获取属性值并转换为目标类型
            var value = property.GetValue(response);
            if (value is TResp typedValue)
                return typedValue;
            
            throw new Exception($"无法从响应中获取类型 {typeof(TResp).Name} 的数据");
        }

        private static long GenerateRequestId() => DateTimeOffset.UtcNow.Ticks;
    }
}