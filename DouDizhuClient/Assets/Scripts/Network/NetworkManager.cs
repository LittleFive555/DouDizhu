using System;
using System.Net.Sockets;
using UnityEngine;
using Network.Tcp;
using Network.Proto;
using Google.Protobuf;
using System.Threading.Tasks;
using System.Collections.Concurrent;
using System.Threading;

namespace Network
{
    public struct NetworkResult
    {
        public bool IsSuccess { get; }
        public string ErrorMessage { get; }

        private NetworkResult(bool isSuccess, string errorMessage)
        {
            IsSuccess = isSuccess;
            ErrorMessage = errorMessage;
        }

        public static NetworkResult Success() => new NetworkResult(true, null);
        public static NetworkResult Failure(string errorMessage) => new NetworkResult(false, errorMessage);
    }

    public struct NetworkResult<T>
    {
        public bool IsSuccess { get; }
        public T Data { get; }
        public string ErrorMessage { get; }

        private NetworkResult(bool isSuccess, T data, string errorMessage)
        {
            IsSuccess = isSuccess;
            Data = data;
            ErrorMessage = errorMessage;
        }

        public static NetworkResult<T> Success(T data) => new NetworkResult<T>(true, data, null);
        public static NetworkResult<T> Failure(string errorMessage) => new NetworkResult<T>(false, default, errorMessage);
    }

    public class NetworkManager
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

        private const int REQUEST_TIMEOUT = 5;

        private readonly ConcurrentDictionary<long, TaskCompletionSource<GameMsgRespPacket>> m_PendingRequests = new();

        public NetworkManager()
        {
            m_MessageReadWriter = new LengthPrefixReadWriter();
        }

        public async Task ConnectAsync(string serverIp, int port)
        {
            try
            {
                m_TcpClient = new TcpClient();
                await m_TcpClient.ConnectAsync(serverIp, port);
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

        public async Task<NetworkResult<CommonResponse>> RequestAsync<TReq>(GameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            return await RequestAsync<TReq, CommonResponse>(requestType, request);
        }

        public async Task<NetworkResult<TResp>> RequestAsync<TReq, TResp>(GameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage where TResp : IMessage
        {
            if (!m_IsConnected || m_TcpClient == null || !m_TcpClient.Connected)
                return NetworkResult<TResp>.Failure("TCP连接未建立，无法发送消息");

            try
            {
                GameClientMessage gameClientMessage = PackRequest(requestType, request);
                var tcs = new TaskCompletionSource<GameMsgRespPacket>();
                m_PendingRequests.TryAdd(gameClientMessage.Header.MessageId, tcs);
                
                await m_MessageReadWriter.WriteTo(m_NetworkStream, gameClientMessage.ToByteArray());
                Debug.Log("消息已发送: " + request.ToString());

                // 设置超时
                using var timeoutCts = new CancellationTokenSource(TimeSpan.FromSeconds(REQUEST_TIMEOUT));
                timeoutCts.Token.Register(() => tcs.TrySetCanceled(), useSynchronizationContext: false);
                var serverPacket = await tcs.Task;
                if (serverPacket.ContentCase == GameMsgRespPacket.ContentOneofCase.Error)
                    return NetworkResult<TResp>.Failure($"服务器返回错误: ErrorCode {serverPacket.Error.Code}: {serverPacket.Error.Message}");
                
                TResp response = UnpackResponse<TResp>(serverPacket);
                return NetworkResult<TResp>.Success(response);
            }
            catch (Exception ex)
            {
                return NetworkResult<TResp>.Failure($"发送消息时发生错误: {ex.Message}");
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

                    GameServerMessage gameServerMessage = GameServerMessage.Parser.ParseFrom(responseBuffer);
                    if (gameServerMessage.ContentCase == GameServerMessage.ContentOneofCase.Notification)
                    {
                        GameNotificationPacket gameNotificationPacket = gameServerMessage.Notification;
                        OnNotified(gameNotificationPacket);
                    }
                    else if (gameServerMessage.ContentCase == GameServerMessage.ContentOneofCase.Response)
                    {
                        GameMsgRespPacket gameMsgRespPacket = gameServerMessage.Response;
                        if (m_PendingRequests.TryGetValue(gameMsgRespPacket.Header.MessageId, out var tcs))
                        {
                            tcs.SetResult(gameMsgRespPacket);
                            m_PendingRequests.TryRemove(gameMsgRespPacket.Header.MessageId, out _);
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                Debug.LogError("接收消息时发生错误: " + ex.Message);
            }
        }

        public void OnNotified(GameNotificationPacket notification)
        {

        }

        public void OnResponse(GameMsgRespPacket response)
        {

        }

        private GameClientMessage PackRequest<TReq>(GameClientMessage.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            GameClientMessage gameClientMessage = new GameClientMessage()
            {
                Header = new GameMsgHeader() { MessageId = GenerateRequestId() },
            };
            
            string propertyName = requestType.ToString();
            Type packetType = typeof(GameClientMessage);
            var property = packetType.GetProperty(propertyName);
            if (property != null)
                property.SetValue(gameClientMessage, request);
            else
                Debug.LogError($"未找到对应的属性: {propertyName}");
            return gameClientMessage;
        }

        private TResp UnpackResponse<TResp>(GameMsgRespPacket response) where TResp : IMessage
        {
            string propertyName = response.ContentCase.ToString();
            
            Type packetType = typeof(GameMsgRespPacket);
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

        private TNotification UnpackNotification<TNotification>(GameNotificationPacket notification) where TNotification : IMessage
        {
            string propertyName = notification.ContentCase.ToString();
            Type packetType = typeof(GameNotificationPacket);
            var property = packetType.GetProperty(propertyName);
            if (property != null)
            {
                // 获取属性值并转换为目标类型
                var value = property.GetValue(notification);
                if (value is TNotification typedValue)
                    return typedValue;
            }
            
            Debug.LogError($"无法从通知中获取类型 {typeof(TNotification).Name} 的数据");
            return default(TNotification);
        }

        private static long GenerateRequestId() => DateTimeOffset.UtcNow.Ticks;
    } 
}