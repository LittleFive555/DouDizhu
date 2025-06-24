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
using Gameplay.Player;
using System.Security.Cryptography;

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
        private IMessageReadWriter m_MessageReadWriter;
        private bool m_IsConnected = false;
        public bool IsConnected => m_IsConnected;

        private string m_SessionId;
        private byte[] m_SharedSecret;

        
        private const int SERVER_PORT = 8080;


        private const int REQUEST_TIMEOUT = 5;

        private readonly ConcurrentDictionary<long, TaskCompletionSource<PServerMsg>> m_PendingRequests = new();
        private readonly Dictionary<PMsgId, Action<IMessage>> m_NotificationHandlers = new();

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

        public async Task<NetworkResult<PEmptyResponse>> RequestAsync<TReq>(PMsgId msgId, TReq request) where TReq : class, IMessage
        {
            return await RequestAsync<TReq, PEmptyResponse>(msgId, request);
        }

        public async Task<NetworkResult<TResp>> RequestAsync<TReq, TResp>(PMsgId msgId, TReq request) where TReq : class, IMessage where TResp : class, IMessage
        {
            if (!m_IsConnected || m_TcpClient == null || !m_TcpClient.Connected)
            {
                Log.Error("TCP连接未建立，无法发送消息");
                return NetworkResult<TResp>.Error(null, "TCP连接未建立，无法发送消息");
            }

            try
            {
                PClientMsg clientMsg = PackClientMsg(msgId, request);
                var tcs = new TaskCompletionSource<PServerMsg>();
                m_PendingRequests.TryAdd(clientMsg.Header.UniqueId, tcs);
                
                // TODO 对消息进行加密
                await m_MessageReadWriter.WriteTo(m_NetworkStream, clientMsg.ToByteArray());
                if (IsSecretMessage(msgId))
                    Log.Information("消息已发送: [{requestType}]", msgId);
                else
                    Log.Information("消息已发送: [{requestType}] {request}", msgId, request);

                // 设置超时
                using var timeoutCts = new CancellationTokenSource(TimeSpan.FromSeconds(REQUEST_TIMEOUT));
                timeoutCts.Token.Register(() => tcs.TrySetCanceled(), useSynchronizationContext: false);
                var serverPacket = await tcs.Task;
                if (serverPacket.MsgType == PServerMsgType.Error)
                {
                    PError error = UnpackServerMsg<PError>(serverPacket);
                    if (error.Type == PError.Types.Type.ServerError)
                        Log.Error("服务器内部错误: {errorCode}: {errorMessage}", error.ErrorCode, error.Message);
                    else
                        Log.Error("服务器返回游戏逻辑错误: {errorCode}: {errorMessage}", error.ErrorCode, error.Message);
                    return NetworkResult<TResp>.Error(error.ErrorCode, error.Message);
                }
                
                TResp response = UnpackServerMsg<TResp>(serverPacket);
                if (IsSecretMessage(msgId))
                    Log.Information("收到服务器响应: [{requestType}]", msgId);
                else
                    Log.Information("收到服务器响应: [{requestType}] {response}", msgId, response);
                return NetworkResult<TResp>.Success(response);
            }
            catch (Exception ex)
            {
                Log.Error("发送消息时发生错误: {exception}", ex.Message);
                return NetworkResult<TResp>.Error(null, ex.Message);
            }
        }

        public void RegisterNotificationHandler<TNotification>(PMsgId msgId, Action<TNotification> handler) where TNotification : IMessage
        {
            var adapter = new NotificationHandlerAdapter<TNotification>(handler);
            if (m_NotificationHandlers.ContainsKey(msgId))
                m_NotificationHandlers[msgId] += adapter.Handle;
            else
                m_NotificationHandlers[msgId] = adapter.Handle;
        }

        public void UnregisterNotificationHandler<TNotification>(PMsgId msgId, Action<TNotification> handler) where TNotification : IMessage
        {
            // 查找对应的适配器
            var adapter = new NotificationHandlerAdapter<TNotification>(handler);
            if (m_NotificationHandlers.ContainsKey(msgId))
            {
                m_NotificationHandlers[msgId] -= adapter.Handle;
                if (m_NotificationHandlers[msgId] == null)
                    m_NotificationHandlers.Remove(msgId);
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

                    PServerMsg serverMsg = PServerMsg.Parser.ParseFrom(responseBuffer);
                    if (serverMsg.MsgType == PServerMsgType.Notification)
                        OnNotified(serverMsg);
                    else if (serverMsg.MsgType == PServerMsgType.Response)
                        OnResponse(serverMsg);
                    else if (serverMsg.MsgType == PServerMsgType.Error)
                        OnResponse(serverMsg);
                }
            }
            catch (Exception ex)
            {
                Log.Error("接收消息时发生错误: {exception}", ex.Message);
            }
        }

        public async Task Handshake()
        {
            // 1. 生成客户端临时密钥对
            using var ecdh = ECDiffieHellman.Create(ECCurve.NamedCurves.nistP256);
            var clientPublicKey = ecdh.ExportSubjectPublicKeyInfo();

            // 2. 发送握手请求
            var handshakeRequest = new PHandshakeRequest()
            {
                PublicKey = Convert.ToBase64String(clientPublicKey),
            };
            var response = await RequestAsync<PHandshakeRequest, PHandshakeResponse>(PMsgId.Handshake, handshakeRequest);
            if (!response.IsSuccess)
                return;

            // 3. 解密服务器响应
            var serverPublicKey = Convert.FromBase64String(response.Data.PublicKey);
            using var serverEcdh = ECDiffieHellman.Create();
            serverEcdh.ImportSubjectPublicKeyInfo(serverPublicKey, out _);
            m_SharedSecret = ecdh.DeriveKeyMaterial(serverEcdh.PublicKey);

            Log.Information("握手成功，共享密钥为: {sharedSecret}", Convert.ToBase64String(m_SharedSecret));
        }

        private void OnResponse(PServerMsg serverMsg)
        {
            if (m_PendingRequests.TryGetValue(serverMsg.Header.UniqueId, out var tcs))
            {
                tcs.SetResult(serverMsg);
                m_PendingRequests.TryRemove(serverMsg.Header.UniqueId, out _);
            }
        }

        private void OnNotified(PServerMsg serverMsg)
        {
            if (m_NotificationHandlers.TryGetValue(serverMsg.Header.MsgId, out var handler))
            {
                try
                {
                    IMessage notification = UnpackServerMsg<IMessage>(serverMsg);
                    handler(notification);
                }
                catch (Exception ex)
                {
                    Log.Error("通知处理程序发生错误: {exception}", ex.Message);
                }
            }
            else
                Log.Error("未找到通知处理程序: {notificationType}", serverMsg.Header.MsgId);
        }

        private PClientMsg PackClientMsg<TReq>(PMsgId msgId, TReq request) where TReq : class, IMessage
        {
            PMsgHeader header = new PMsgHeader()
            {
                UniqueId = GenerateRequestId(),
                Timestamp = DateTime.UtcNow.Millisecond,
                MsgId = msgId,
            };
            if (m_SessionId != null)
                header.SessionId = m_SessionId;
            if (!string.IsNullOrEmpty(PlayerManager.Instance.ID))
                header.PlayerId = PlayerManager.Instance.ID;
            if (m_SharedSecret != null)
                header.Encrypted = true;

            PClientMsg clientMsg = new PClientMsg()
            {
                Header = header,
                Payload = request.ToByteString(), // TODO 加密
            };
            
            return clientMsg;
        }

        private TResp UnpackServerMsg<TResp>(PServerMsg serverMsg) where TResp : class, IMessage
        {
            if (serverMsg.MsgType == PServerMsgType.Error)
            {
                PError error = PError.Parser.ParseFrom(serverMsg.Payload);
                return error as TResp;
            }
            else if (serverMsg.MsgType == PServerMsgType.Response)
            {
                var parser = MsgBodyParser.GetResponseParser(serverMsg.Header.MsgId);
                if (parser == null)
                    throw new Exception($"未找到响应解析器: {serverMsg.Header.MsgId}");
                
                return parser.ParseFrom(serverMsg.Payload) as TResp;
            }
            else if (serverMsg.MsgType == PServerMsgType.Notification)
            {
                var parser = MsgBodyParser.GetNotifyParser(serverMsg.Header.MsgId);
                if (parser == null)
                    throw new Exception($"未找到通知解析器: {serverMsg.Header.MsgId}");
                return parser.ParseFrom(serverMsg.Payload) as TResp;
            }
            throw new Exception($"未知的消息类型: {serverMsg.MsgType}");
        }

        private bool IsSecretMessage(PMsgId msgId)
        {
            if (msgId == PMsgId.Login
            || msgId == PMsgId.Register
            || msgId == PMsgId.Handshake) {
                return true;
            }
            return false;
        }

        private static long GenerateRequestId() => DateTimeOffset.UtcNow.Ticks;
    }
}