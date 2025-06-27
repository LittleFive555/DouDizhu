using System.Threading.Tasks;
using UnityEngine;
using TMPro;
using Network;
using Network.Proto;
using UIModule;
using Serilog;

namespace Gameplay.Chat.View
{
    public class UIChatBox : UIWidget
    {
        [SerializeField]
        private TMP_InputField messageInput;

        public static string ID;
        public static string Name;

        private void Start()
        {
            NetworkManager.Instance.RegisterNotificationHandler<PChatMsgNotification>(PMsgId.ChatMsg, OnReceivedChatMsg);
        }

        private async Task SendMessageImpl(string message)
        {
            if (string.IsNullOrEmpty(message))
                return;

            var response = await NetworkManager.Instance.RequestAsync(PMsgId.ChatMsg, new PChatMsgRequest() { Content = message });
            if (response.IsSuccess)
            {
                messageInput.text = ""; // 清空输入框
            }
        }

        private void OnReceivedChatMsg(PChatMsgNotification notification)
        {
            Log.Information("收到聊天消息 {playerName}: {content}", notification.From.Nickname, notification.Content);
        }

        [OnClick("BtnSend")]
        private async void SendMessage()
        {
            await SendMessageImpl(messageInput.text);
        }
    }
}
