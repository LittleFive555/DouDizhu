using System.Text;
using System.Net.Sockets;
using UnityEngine;
using UnityEngine.UI;
using TMPro;

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
                byte[] data = Encoding.UTF8.GetBytes(message);
                networkStream.Write(data, 0, data.Length);
                Debug.Log("消息已发送: " + message);
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
