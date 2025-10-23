using System.Collections.Generic;
using UnityEngine;
using Network;
using Network.Proto;
using TMPro;
using UIModule;
using Gameplay.Chat.View;
using Gameplay.Room.Controller;

namespace Gameplay.Room.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Room/UIRoom.prefab")]
    public class UIRoom : UIComponentBase<UIRoom.Args>
    {
        public struct Args
        {
            public PRoom Room;
        }

        [SerializeField]
        private TextMeshProUGUI m_TextRoomName;

        [SerializeField]
        private UIChatBox m_ChatBox;

        [SerializeField]
        private List<UIRoomPlayerSlot> m_PlayerSlots = new List<UIRoomPlayerSlot>();

        private static readonly IReadOnlyList<PChatChannel> CHANNEL_SET = new List<PChatChannel>() { PChatChannel.Room, PChatChannel.All };

        public override void OnShowBegin(Args args)
        {
            SetRoomName(args.Room.Name);
            RefreshPlayers(args.Room.Players);
            m_ChatBox.InitChannelSet(CHANNEL_SET, PChatChannel.Room);
            NetworkManager.Instance.RegisterNotificationHandler<PRoomChangedNotification>(PMsgId.RoomChanged, OnReceivedRoomChanged);
            NetworkManager.Instance.RegisterNotificationHandler<PRoomDisbandedNotification>(PMsgId.RoomDisbanded, OnReceivedRoomDisbanded);
        }

        public override void OnHideBegin()
        {
            NetworkManager.Instance.UnregisterNotificationHandler<PRoomChangedNotification>(PMsgId.RoomChanged, OnReceivedRoomChanged);
            NetworkManager.Instance.UnregisterNotificationHandler<PRoomDisbandedNotification>(PMsgId.RoomDisbanded, OnReceivedRoomDisbanded);

            base.OnHideBegin();
        }

        private void SetRoomName(string name)
        {
            m_TextRoomName.text = name;
        }

        private void RefreshPlayers(IReadOnlyList<PPlayer> players)
        {
            for (int i = 0; i < m_PlayerSlots.Count; i++)
            {
                if (i < players.Count)
                    m_PlayerSlots[i].SetPlayerInfo(players[i]);
                else
                    m_PlayerSlots[i].SetEmpty();
            }
        }

        private void OnReceivedRoomChanged(PRoomChangedNotification notification)
        {
            if (!string.IsNullOrEmpty(notification.Room.Name))
                SetRoomName(notification.Room.Name);
            
            if (!string.IsNullOrEmpty(notification.Room.OwnerId))
            {
                // TODO 房主变更
            }

            if (notification.Room.Players != null && notification.Room.Players.Count > 0)
            {
                RefreshPlayers(notification.Room.Players);
            }
        }

        private void OnReceivedRoomDisbanded(PRoomDisbandedNotification _)
        {
            Hide();
        }

        [OnClick("BtnBack")]
        private void OnClickClose()
        {
            _ = RoomController.LeaveRoom();
        }
    }
}