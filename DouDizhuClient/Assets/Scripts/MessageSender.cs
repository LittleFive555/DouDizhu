using UnityEngine;
using UnityEngine.UI;
using TMPro;
using Network;
using Network.Proto;

public class MessageSender : MonoBehaviour
{
    [SerializeField]
    private TMP_InputField messageInput;
    [SerializeField]
    private Button sendButton;

    void Start()
    {
        // 替换为你的服务器IP和端口
        string serverIp = "127.0.0.1";
        int port = 8080;

        // 连接服务器
        NetworkManager.Instance.Connect(serverIp, port);

        // 设置按钮点击事件
        sendButton.onClick.AddListener(SendMessage);
    }

    void SendMessage()
    {
        SendMessageImpl(messageInput.text);
    }

    private void SendMessageImpl(string message)
    {
        if (string.IsNullOrEmpty(message))
            return;

        var response = NetworkManager.Instance.Request<ChatMsgRequest, ChatMsgResponse>(GameMsgReqPacket.ContentOneofCase.ChatMsg, new ChatMsgRequest() { Content = message });
        if (response.IsSuccess)
        {
            messageInput.text = ""; // 清空输入框
        }
    }

    void OnDestroy()
    {
        NetworkManager.Instance.Disconnect();
    }
}
