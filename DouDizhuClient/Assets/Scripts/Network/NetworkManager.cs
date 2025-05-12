using System;
using System.Net.Sockets;
using UnityEngine;
using Network.Tcp;
using Network.Proto;
using Google.Protobuf;
using System.Threading.Tasks;

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
        private static NetworkManager instance;
        public static NetworkManager Instance
        {
            get
            {
                if (instance == null)
                    instance = new NetworkManager();
                return instance;
            }
        }

        private TcpClient tcpClient;
        private NetworkStream networkStream;
        private MessageReadWriter messageReadWriter;
        private bool isConnected = false;

        public bool IsConnected => isConnected;

        public NetworkManager()
        {
            messageReadWriter = new LengthPrefixReadWriter();
        }

        public void Connect(string serverIp, int port)
        {
            try
            {
                tcpClient = new TcpClient();
                tcpClient.Connect(serverIp, port);
                networkStream = tcpClient.GetStream();
                isConnected = true;
                Debug.Log("TCP连接已建立");
            }
            catch (SocketException ex)
            {
                Debug.LogError("TCP连接失败: " + ex.Message);
                isConnected = false;
            }
        }

        public void Disconnect()
        {
            if (networkStream != null)
            {
                networkStream.Close();
                networkStream = null;
            }
            if (tcpClient != null)
            {
                tcpClient.Close();
                tcpClient = null;
            }
            isConnected = false;
            Debug.Log("TCP连接已关闭");
        }

        public async Task<NetworkResult<CommonResponse>> RequestAsync<TReq>(GameMsgReqPacket.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            return await RequestAsync<TReq, CommonResponse>(requestType, request);
        }

        public async Task<NetworkResult<TResp>> RequestAsync<TReq, TResp>(GameMsgReqPacket.ContentOneofCase requestType, TReq request) where TReq : IMessage where TResp : IMessage
        {
            if (!isConnected || tcpClient == null || !tcpClient.Connected)
                return NetworkResult<TResp>.Failure("TCP连接未建立，无法发送消息");

            try
            {
                GameMsgReqPacket gameMsgReqPacket = PackRequest(requestType, request);
                await messageReadWriter.WriteTo(networkStream, gameMsgReqPacket.ToByteArray());
                Debug.Log("消息已发送: " + request.ToString());

                // 读取响应数据
                byte[] responseBuffer = await messageReadWriter.ReadFrom(networkStream);
                if (responseBuffer == null || responseBuffer.Length == 0)
                    return NetworkResult<TResp>.Failure("未收到服务器响应");

                GameMsgRespPacket gameMsgRespPacket = GameMsgRespPacket.Parser.ParseFrom(responseBuffer);
                if (gameMsgRespPacket.ContentCase == GameMsgRespPacket.ContentOneofCase.Error)
                    return NetworkResult<TResp>.Failure($"服务器返回错误: ErrorCode {gameMsgRespPacket.Error.Code}: {gameMsgRespPacket.Error.Message}");

                TResp response = UnpackResponse<TResp>(gameMsgRespPacket);
                if (response == null)
                    return NetworkResult<TResp>.Failure($"无法解析响应数据为类型 {typeof(TResp).Name}");

                return NetworkResult<TResp>.Success(response);
            }
            catch (Exception ex)
            {
                return NetworkResult<TResp>.Failure($"发送消息时发生错误: {ex.Message}");
            }
        }

        public void OnNotified(GameNotificationPacket notification)
        {

        }

        private GameMsgReqPacket PackRequest<TReq>(GameMsgReqPacket.ContentOneofCase requestType, TReq request) where TReq : IMessage
        {
            GameMsgReqPacket gameMsgReqPacket = new GameMsgReqPacket()
            {
                Header = new GameMsgHeader() { },
            };
            
            string propertyName = requestType.ToString();
            Type packetType = typeof(GameMsgReqPacket);
            var property = packetType.GetProperty(propertyName);
            if (property != null)
                property.SetValue(gameMsgReqPacket, request);
            else
                Debug.LogError($"未找到对应的属性: {propertyName}");
            return gameMsgReqPacket;
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
    } 
}