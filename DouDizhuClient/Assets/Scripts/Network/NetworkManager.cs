using System;
using System.Net.Sockets;
using System.Threading.Tasks;
using System.Collections.Concurrent;
using System.Threading;
using System.Collections.Generic;
using Google.Protobuf;
using System.Text;
using Network.Tcp;
using Network.Proto;
using Serilog;
using Gameplay.Player;
using Org.BouncyCastle.Crypto;
using Org.BouncyCastle.Crypto.Parameters;
using Org.BouncyCastle.Crypto.Generators;
using Org.BouncyCastle.Asn1.Sec;
using Org.BouncyCastle.Security;
using Org.BouncyCastle.Crypto.Agreement;
using Org.BouncyCastle.Crypto.Digests;
using Org.BouncyCastle.Crypto.Engines;
using Org.BouncyCastle.Crypto.Paddings;
using Org.BouncyCastle.Crypto.Modes;
using Org.BouncyCastle.Math;

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
        
        public byte[] PublicKeyBytes { get; private set; }
        private ECPrivateKeyParameters privateKey;
        private ECDomainParameters domainParams;
        private byte[] m_DerivedSecureKey;

        private readonly ConcurrentQueue<MessageDataToSend> m_MessageDataToSend = new();
        
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

            PClientMsg clientMsg;
            try
            {
                clientMsg = PackClientMsg(msgId, request);
            }
            catch (Exception ex)
            {
                Log.Error(ex, "打包消息时发生错误");
                return NetworkResult<TResp>.Error(null, ex.Message);
            }

            var tcs = new TaskCompletionSource<PServerMsg>();
            m_PendingRequests.TryAdd(clientMsg.Header.UniqueId, tcs);

            SendMessage(new MessageDataToSend(msgId, clientMsg.ToByteArray(), request.ToString()));

            // 设置超时
            using var timeoutCts = new CancellationTokenSource(TimeSpan.FromSeconds(REQUEST_TIMEOUT));
            timeoutCts.Token.Register(() => tcs.TrySetCanceled(), useSynchronizationContext: false);
            var serverPacket = await tcs.Task;

            try
            {
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
                if (IsSensitiveMessage(msgId))
                    Log.Information("收到服务器响应: [{requestType}]", msgId);
                else
                    Log.Information("收到服务器响应: [{requestType}] {response}", msgId, response);

                return NetworkResult<TResp>.Success(response);
            }
            catch (Exception ex)
            {
                Log.Error(ex, "解析服务器响应时发生错误");
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
                Log.Error(ex, "TCP连接失败");
                m_IsConnected = false;
            }
        }

        private void SendMessage(MessageDataToSend messageDataToSend)
        {
            if (m_MessageDataToSend.Count > 0)
            {
                m_MessageDataToSend.Enqueue(messageDataToSend);
            }
            else
            {
                m_MessageDataToSend.Enqueue(messageDataToSend);
                _ = SendMessageFromQueueAsync();
            }
        }

        private async Task SendMessageFromQueueAsync()
        {
            while (m_MessageDataToSend.Count > 0)
            {
                if (!m_IsConnected) // TODO: 等待连接成功后重试
                {
                    Log.Error("连接已断开");
                    break;
                }

                if (m_MessageDataToSend.TryDequeue(out var toSend))
                {
                    try
                    {
                        await m_MessageReadWriter.WriteTo(m_NetworkStream, toSend.MessageBytes);
                    }
                    catch (Exception ex)
                    {
                        Log.Error(ex, "发送消息时发生错误");
                        continue;
                    }

                    if (IsSensitiveMessage(toSend.MsgId))
                        Log.Information("消息已发送: [{requestType}]", toSend.MsgId);
                    else
                        Log.Information("消息已发送: [{requestType}] {request}", toSend.MsgId, toSend.PayloadString);
                }
            }
        }

        private async Task ReceiveLoopAsync()
        {
            while (m_IsConnected)
            {
                byte[] responseBuffer;
                try
                {
                    responseBuffer = await m_MessageReadWriter.ReadFrom(m_NetworkStream);
                }
                catch (Exception ex)
                {
                    m_IsConnected = false;
                    Log.Error(ex, "接收消息时发生错误");
                    break;
                }

                if (responseBuffer == null || responseBuffer.Length == 0)
                {
                    m_IsConnected = false;
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

        public async Task Handshake()
        {
            // 1. 生成客户端临时密钥对
            var ecParams = SecNamedCurves.GetByName("secp256r1");
            domainParams = new ECDomainParameters(
                ecParams.Curve, ecParams.G, ecParams.N, ecParams.H);
            
            var keyPair = GenerateKeyPair(domainParams);
            privateKey = (ECPrivateKeyParameters)keyPair.Private;
            var publicKey = (ECPublicKeyParameters)keyPair.Public;
            PublicKeyBytes = publicKey.Q.GetEncoded(false);

            byte[] salt = new byte[16];
            new SecureRandom().NextBytes(salt);
            byte[] info = Encoding.UTF8.GetBytes("global encryption key");

            // 2. 发送握手请求
            var handshakeRequest = new PHandshakeRequest()
            {
                PublicKey = Convert.ToBase64String(PublicKeyBytes),
                Salt = ByteString.CopyFrom(salt),
                Info = ByteString.CopyFrom(info),
            };
            var response = await RequestAsync<PHandshakeRequest, PHandshakeResponse>(PMsgId.Handshake, handshakeRequest);
            if (!response.IsSuccess)
                return;

            // 3. 解密服务器响应
            var serverPublicKey = Convert.FromBase64String(response.Data.PublicKey);
            var sharedSecret = DeriveSharedSecret(serverPublicKey);

            // 4. 派生安全密钥
            m_DerivedSecureKey = DeriveSecureKey(sharedSecret, salt, info, 32);
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
                    Log.Error(ex, "通知处理程序发生错误");
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

            byte[] payload = request.ToByteArray();
            // 加密
            if (m_DerivedSecureKey != null)
            {
                (byte[] iv, byte[] ciphertext) = Encrypt(payload, m_DerivedSecureKey);
                header.Iv = ByteString.CopyFrom(iv);
                payload = ciphertext;
            }
            PClientMsg clientMsg = new PClientMsg()
            {
                Header = header,
                Payload = ByteString.CopyFrom(payload),
            };
            
            return clientMsg;
        }

        private TResp UnpackServerMsg<TResp>(PServerMsg serverMsg) where TResp : class, IMessage
        {
            byte[] payload = GetPlaintextPayload(serverMsg);
            if (serverMsg.MsgType == PServerMsgType.Error)
            {
                return PError.Parser.ParseFrom(payload) as TResp;
            }
            else if (serverMsg.MsgType == PServerMsgType.Response)
            {
                var parser = MsgBodyParser.GetResponseParser(serverMsg.Header.MsgId);
                if (parser == null)
                    throw new Exception($"未找到响应解析器: {serverMsg.Header.MsgId}");
                
                return parser.ParseFrom(payload) as TResp;
            }
            else if (serverMsg.MsgType == PServerMsgType.Notification)
            {
                var parser = MsgBodyParser.GetNotifyParser(serverMsg.Header.MsgId);
                if (parser == null)
                    throw new Exception($"未找到通知解析器: {serverMsg.Header.MsgId}");
                return parser.ParseFrom(payload) as TResp;
            }
            throw new Exception($"未知的消息类型: {serverMsg.MsgType}");
        }

        private bool IsSensitiveMessage(PMsgId msgId)
        {
            if (msgId == PMsgId.Login
            || msgId == PMsgId.Register
            || msgId == PMsgId.Handshake) {
                return true;
            }
            return false;
        }

        private static long GenerateRequestId() => DateTimeOffset.UtcNow.Ticks;

        private static AsymmetricCipherKeyPair GenerateKeyPair(ECDomainParameters domainParams)
        {
            var generator = new ECKeyPairGenerator();
            var keyGenParams = new ECKeyGenerationParameters(domainParams, new SecureRandom());
            generator.Init(keyGenParams);
            return generator.GenerateKeyPair();
        }

        private byte[] DeriveSharedSecret(byte[] otherPartyPublicKeyBytes)
        {
            // 导入对方公钥
            var curve = domainParams.Curve;
            var otherPartyPoint = curve.DecodePoint(otherPartyPublicKeyBytes);
            var otherPartyPublicKey = new ECPublicKeyParameters(otherPartyPoint, domainParams);
            
            // 计算共享密钥
            var agreement = new ECDHBasicAgreement();
            agreement.Init(privateKey);
            var sharedSecret = agreement.CalculateAgreement(otherPartyPublicKey);
            
            return ToFixedLengthBytes(sharedSecret, 32); // 32 字节对齐
        }
        
        byte[] ToFixedLengthBytes(BigInteger value, int length)
        {
            byte[] bytes = value.ToByteArrayUnsigned();
            if (bytes.Length == length) return bytes;
            
            byte[] result = new byte[length];
            Buffer.BlockCopy(
                src: bytes,
                srcOffset: Math.Max(0, bytes.Length - length),
                dst: result,
                dstOffset: Math.Max(0, length - bytes.Length),
                count: Math.Min(bytes.Length, length)
            );
            return result;
        }

        private byte[] DeriveSecureKey(byte[] sharedSecret, byte[] salt, byte[] info, int outputLength)
        {
            HkdfParameters parameters = new HkdfParameters(sharedSecret, salt, info);
            HkdfBytesGenerator hkdf = new HkdfBytesGenerator(new Sha256Digest());
            hkdf.Init(parameters);
            byte[] output = new byte[outputLength];
            hkdf.GenerateBytes(output, 0, outputLength);
            return output;
        }

        private byte[] GetPlaintextPayload(PServerMsg serverMsg)
        {
            if (serverMsg.Header.Iv != null && serverMsg.Header.Iv.Length > 0)
            {
                byte[] iv = serverMsg.Header.Iv.ToByteArray();
                byte[] ciphertext = serverMsg.Payload.ToByteArray();
                return Decrypt(ciphertext, m_DerivedSecureKey, iv);
            }
            return serverMsg.Payload.ToByteArray();
        }

        private (byte[] iv, byte[] ciphertext) Encrypt(byte[] plaintext, byte[] key)
        {
            // 生成随机IV(块大小为16字节)
            byte[] iv = new byte[16];
            new SecureRandom().NextBytes(iv);

            // 初始化加密引擎 CBC模式
            AesEngine engine = new AesEngine();  // 不建议存储引擎实例
            PaddedBufferedBlockCipher cipher = new PaddedBufferedBlockCipher(new CbcBlockCipher(engine), new Pkcs7Padding());
            cipher.Init(true, new ParametersWithIV(new KeyParameter(key), iv));

            // 执行加密
            byte[] output = new byte[cipher.GetOutputSize(plaintext.Length)];
            int len = cipher.ProcessBytes(plaintext, 0, plaintext.Length, output, 0);
            cipher.DoFinal(output, len);

            return (iv, output);
        }

        private byte[] Decrypt(byte[] ciphertext, byte[] key, byte[] iv)
        {
            AesEngine engine = new AesEngine();
            PaddedBufferedBlockCipher cipher = new PaddedBufferedBlockCipher(new CbcBlockCipher(engine), new Pkcs7Padding());
            cipher.Init(false, new ParametersWithIV(new KeyParameter(key), iv));

            byte[] output = new byte[cipher.GetOutputSize(ciphertext.Length)];
            int len = cipher.ProcessBytes(ciphertext, 0, ciphertext.Length, output, 0);
            int finalLen = cipher.DoFinal(output, len);
            
            // 创建正确大小的数组，只包含实际的明文数据（去除填充）
            byte[] result = new byte[len + finalLen];
            Array.Copy(output, 0, result, 0, len + finalLen);
            
            return result;
        }

        private struct MessageDataToSend
        {
            public PMsgId MsgId;
            public byte[] MessageBytes;
            public string PayloadString;

            public MessageDataToSend(PMsgId msgId, byte[] messageBytes, string payloadString)
            {
                MsgId = msgId;
                MessageBytes = messageBytes;
                PayloadString = payloadString;
            }
        }
    }
}