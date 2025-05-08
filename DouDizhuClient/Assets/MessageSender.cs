using System.Text;
using System.Net.Sockets;
using System;
using UnityEngine;
using UnityEngine.UI;
using TMPro;
using Network.Proto;
using Google.Protobuf;

public class MessageSender : MonoBehaviour
{
    private TcpClient tcpClient;
    private NetworkStream networkStream;

    [SerializeField]
    private TMP_InputField messageInput;
    [SerializeField]
    private Button sendButton;

    // Start is called once before the first execution of Update after the MonoBehaviour is created
    void Start()
    {
        // 替换为你的服务器IP和端口
        string serverIp = "127.0.0.1";
        int port = 8080;

        tcpClient = new TcpClient();
        try
        {
            tcpClient.Connect(serverIp, port);
            networkStream = tcpClient.GetStream();
            Debug.Log("TCP连接已建立");

            // 设置按钮点击事件
            sendButton.onClick.AddListener(SendMessage);
        }
        catch (SocketException ex)
        {
            Debug.LogError("TCP连接失败: " + ex.Message);
        }
    }

    void SendMessage()
    {
        if (tcpClient != null && tcpClient.Connected)
        {
            string message = messageInput.text;
            if (!string.IsNullOrEmpty(message))
            {
                GameMsgReqPacket gameMsgReqPacket = new GameMsgReqPacket()
                {
                    Header = new GameMsgHeader() {},
                    ChatMsg = new ChatMsgRequest() {
                        Content = message
                    }
                };
                gameMsgReqPacket.WriteTo(networkStream);
                Debug.Log("消息已发送: " + message);
                networkStream.Flush();

                // 读取响应数据
                byte[] responseBuffer = new byte[1024];
                int bytesRead = networkStream.Read(responseBuffer, 0, responseBuffer.Length);
                if (bytesRead > 0)
                {
                    // 创建新的字节数组，只包含实际读取的数据
                    byte[] actualData = new byte[bytesRead];
                    Array.Copy(responseBuffer, actualData, bytesRead);
                    
                    // 解析响应数据
                    GameMsgRespPacket response = GameMsgRespPacket.Parser.ParseFrom(actualData);
                    
                    // 处理响应
                    if (response.ContentCase == GameMsgRespPacket.ContentOneofCase.Error)
                    {
                        Debug.LogError("服务器返回错误: " + response.Error.Message);
                    }
                    else if (response.ContentCase == GameMsgRespPacket.ContentOneofCase.ChatMsg)
                    {
                        Debug.Log("消息发送成功，收到服务器确认");
                    }
                }
                else
                {
                    Debug.LogError("未收到服务器响应");
                }
                messageInput.text = ""; // 清空输入框
            }
        }
        else
        {
            Debug.LogError("TCP连接未建立，无法发送消息");
        }
    }

    void OnDestroy()
    {
        // 关闭连接
        if (networkStream != null)
        {
            networkStream.Close();
        }
        if (tcpClient != null)
        {
            tcpClient.Close();
            Debug.Log("TCP连接已关闭");
        }
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
