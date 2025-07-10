using System.Linq;
using System.Collections.Generic;
using UnityEngine;
using TMPro;
using Network;
using Network.Proto;
using UIModule;
using Serilog;
using Gameplay.Chat.Service;

namespace Gameplay.Chat.View
{
    public class UIChatBox : UIWidget
    {
        [SerializeField]
        private TMP_InputField m_MessageInput;

        [SerializeField]
        private TMP_Dropdown m_ChannelDropdown;

        [SerializeField]
        private Transform m_MessageParent;

        [SerializeField]
        private UIChatMessage m_MessagePrototype;

        private List<PChatMsgNotification> m_Messages = new List<PChatMsgNotification>();

        private List<PChatChannel> m_ChannelSet;
        private PChatChannel m_CurrentChannel = PChatChannel.None;

        protected override void Awake()
        {
            base.Awake();

            m_MessagePrototype.gameObject.SetActive(false);
            m_ChannelDropdown.onValueChanged.AddListener(OnChannelChanged);
        }

        private void Start()
        {
            NetworkManager.Instance.RegisterNotificationHandler<PChatMsgNotification>(PMsgId.ChatMsg, OnReceivedChatMsg);
        }

        private void OnDestroy()
        {
            NetworkManager.Instance.UnregisterNotificationHandler<PChatMsgNotification>(PMsgId.ChatMsg, OnReceivedChatMsg);
        }

        public void InitChannelSet(IReadOnlyList<PChatChannel> channels, PChatChannel defaultChannel)
        {
            m_ChannelSet = new List<PChatChannel>(channels);
            m_ChannelDropdown.ClearOptions();
            m_ChannelDropdown.AddOptions(channels.Select(c => new TMP_Dropdown.OptionData(c.ToString())).ToList());
            ChangeChannel(defaultChannel);
        }

        public void ChangeChannel(PChatChannel channel)
        {
            m_CurrentChannel = channel;
            m_ChannelDropdown.value = m_ChannelSet.IndexOf(channel);
        }

        private void OnChannelChanged(int index)
        {
            m_CurrentChannel = m_ChannelSet[index];
        }

        private void OnReceivedChatMsg(PChatMsgNotification notification)
        {
            Log.Information("收到聊天消息 {playerName}: {content}", notification.From.Nickname, notification.Content);
            m_Messages.Add(notification);
            var messageInstance = Instantiate(m_MessagePrototype, m_MessageParent);
            messageInstance.gameObject.SetActive(true);
            messageInstance.SetMessage(notification.From.Nickname, notification.Content, notification.Channel);
        }

        [OnClick("BtnSend")]
        private async void SendMessage()
        {
            var message = m_MessageInput.text;
            if (string.IsNullOrEmpty(message))
                return;
                
            if (m_ChannelSet == null || !m_ChannelSet.Contains(m_CurrentChannel))
            {
                Log.Error("当前频道{channel}目前不可用，可用列表为：{channels}", m_CurrentChannel, string.Join(", ", m_ChannelSet));
                return;
            }

            var success = await ChatService.SendMessage(message, m_CurrentChannel);
            if (success)
            {
                m_MessageInput.text = "";
            }
        }
    }
}
