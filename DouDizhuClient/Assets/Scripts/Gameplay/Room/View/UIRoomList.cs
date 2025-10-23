using System.Collections.Generic;
using UnityEngine;
using Network.Proto;
using UIModule;
using Gameplay.Player.Model;
using Gameplay.Chat.View;
using Gameplay.Room.Controller;

namespace Gameplay.Room.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Room/UIRoomList.prefab")]
    public class UIRoomList : UIComponentBase
    {
        [SerializeField]
        private UIRoomListItem m_RoomItemPrefab;

        [SerializeField]
        private Transform m_RoomItemParent;

        [SerializeField]
        private UIChatBox m_ChatBox;

        private Dictionary<uint, PRoom> m_RoomsInfo = new Dictionary<uint, PRoom>();
        private Dictionary<uint, UIRoomListItem> m_RoomItems = new Dictionary<uint, UIRoomListItem>();

        private bool m_IsRequesting = false;

        protected override void Awake()
        {
            base.Awake();

            m_RoomItemPrefab.gameObject.SetActive(false);
        }

        public override void OnShowBegin(object args)
        {
            base.OnShowBegin(args);

            RefreshRoomList();

            m_ChatBox.InitChannelSet(new List<PChatChannel>() { PChatChannel.All }, PChatChannel.All);
        }

        [OnClick("BtnCreateRoom")]
        private void CreateRoom()
        {
            _ = RoomController.CreateRoom(string.Format("{0}'s Room", PlayerManager.Instance.Name));
        }

        [OnClick("BtnRefresh")]
        private async void RefreshRoomList()
        {
            if (m_IsRequesting)
                return;

            m_IsRequesting = true;
            var roomList = await RoomController.RefreshRoomList();
            if (roomList == null)
            {
                m_IsRequesting = false;
                return;
            }

            m_RoomsInfo.Clear();
            foreach (var room in roomList)
                m_RoomsInfo.Add(room.Id, room);

            foreach (var roomItem in m_RoomItems.Values)
                Destroy(roomItem.gameObject);
            m_RoomItems.Clear();

            foreach (var room in m_RoomsInfo.Values)
            {
                var roomItem = Instantiate(m_RoomItemPrefab, m_RoomItemParent);
                roomItem.gameObject.SetActive(true);
                roomItem.SetRoomInfo(room);

                m_RoomItems.Add(room.Id, roomItem);
            }

            m_IsRequesting = false;
        }
    }
}