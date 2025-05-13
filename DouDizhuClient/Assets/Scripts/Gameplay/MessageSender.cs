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
    [SerializeField]
    private TMP_InputField IDInput;
    [SerializeField]
    private TMP_InputField NameInput;

    public static string ID;
    public static string Name;

    void Start()
    {
        // 设置按钮点击事件
        sendButton.onClick.AddListener(SendMessage);

        NetworkManager.Instance.RegisterNotificationHandler<ChatMsgNotification>(GameNotificationPacket.ContentOneofCase.ChatMsg, OnReceivedChatMsg);
    }

    void SendMessage()
    {
        SendMessageImpl(messageInput.text);
    }

    private async void SendMessageImpl(string message)
    {
        if (string.IsNullOrEmpty(message))
            return;
        
        ID = IDInput.text;
        Name = NameInput.text;

        var response = await NetworkManager.Instance.RequestAsync(GameClientMessage.ContentOneofCase.ChatMsg, new ChatMsgRequest() { Content = message });
        if (response.IsSuccess)
        {
            messageInput.text = ""; // 清空输入框
        }
    }

    private void OnReceivedChatMsg(ChatMsgNotification notification)
    {
        Debug.Log($"收到聊天消息: {notification.Content}");
    }
}
