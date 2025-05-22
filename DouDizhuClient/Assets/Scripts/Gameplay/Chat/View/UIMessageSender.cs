using UnityEngine;
using UnityEngine.UI;
using TMPro;
using Network;
using Network.Proto;
using UIModule;
using System.Threading.Tasks;
using Serilog;

namespace Gameplay.Chat.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Chat/UIChat.prefab")]
    public class UIMessageSender : UIComponentBase
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

            NetworkManager.Instance.RegisterNotificationHandler<PChatMsgNotification>(PGameNotificationPacket.ContentOneofCase.ChatMsg, OnReceivedChatMsg);
        }

        void SendMessage()
        {
            SendMessageImpl(messageInput.text);
        }

        private async Task SendMessageImpl(string message)
        {
            if (string.IsNullOrEmpty(message))
                return;
            
            ID = IDInput.text;
            Name = NameInput.text;

            var response = await NetworkManager.Instance.RequestAsync(PGameClientMessage.ContentOneofCase.ChatMsg, new PChatMsgRequest() { Content = message });
            if (response.IsSuccess)
            {
                messageInput.text = ""; // 清空输入框
            }
        }

        private void OnReceivedChatMsg(PChatMsgNotification notification)
        {
            Log.Information("收到聊天消息 {playerName}: {content}", notification.From.Nickname, notification.Content);
        }
    }
}
