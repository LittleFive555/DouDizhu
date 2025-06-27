using Network.Proto;
using TMPro;
using UIModule;
using UnityEngine;

namespace Gameplay.Room.View
{
    public class UIRoomPlayerSlot : UIWidget
    {
        [SerializeField]
        private TextMeshProUGUI m_TextPlayerName;

        public void SetPlayerInfo(PPlayer player)
        {
            m_TextPlayerName.text = player.Nickname;
        }

        public void SetEmpty()
        {
            m_TextPlayerName.text = "";
        }
    }
}