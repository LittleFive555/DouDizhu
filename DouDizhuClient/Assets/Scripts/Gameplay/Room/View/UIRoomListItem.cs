using Gameplay.Room.Service;
using Network.Proto;
using TMPro;
using UIModule;
using UnityEngine;

namespace Gameplay.Room.View
{
    public class UIRoomListItem : UIWidget
    {
        [SerializeField]
        private TextMeshProUGUI m_TextRoomName;

        [SerializeField]
        private TextMeshProUGUI m_TextPlayerCount;

        [SerializeField]
        private TextMeshProUGUI m_TextRoomState;

        private PRoom m_RoomInfo;

        public void SetRoomInfo(PRoom room)
        {
            m_TextRoomName.text = room.Name;
            m_TextPlayerCount.text = $"{room.Players.Count}/{room.MaxPlayerCount}";
            m_TextRoomState.text = room.State.ToString();

            m_RoomInfo = room;
        }

        [OnClick("BtnEnter")]
        private async void EnterRoom()
        {
            var result = await RoomService.EnterRoom(m_RoomInfo.Id);
            if (result != null)
            {
                UIManager.Instance.ShowUI<UIRoom, UIRoom.Args>(new UIRoom.Args() { Room = result });
            }
        }
    }
}