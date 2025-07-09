using UnityEngine;
using TMPro;
using Network.Proto;
using UIModule;

namespace Gameplay.Chat.View
{
    public class UIChatMessage : UIWidget
    {
        private TMP_Text m_TextMessage;

        protected override void Awake()
        {
            base.Awake();

            m_TextMessage = GetComponent<TMP_Text>();
        }

        public void SetMessage(string name, string message, PChatChannel channel)
        {
            m_TextMessage.text = $"[{name}] {message}";
            if (channel == PChatChannel.All)
                m_TextMessage.color = Color.white;
            else if (channel == PChatChannel.Room)
                m_TextMessage.color = Color.blue;
        }
    }
}
