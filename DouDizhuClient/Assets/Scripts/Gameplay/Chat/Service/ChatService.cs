using System.Threading.Tasks;
using Network;
using Network.Proto;

namespace Gameplay.Chat.Service
{
    public class ChatService
    {
        public static async Task<bool> SendMessage(string message, PChatChannel channel)
        {
            var args = new PChatMsgRequest()
            {
                Content = message,
                Channel = channel
            };
            var response = await NetworkManager.Instance.RequestAsync(PMsgId.ChatMsg, args);
            if (response.IsSuccess)
            {
                return true;
            }
            return false;
        }
    }
}