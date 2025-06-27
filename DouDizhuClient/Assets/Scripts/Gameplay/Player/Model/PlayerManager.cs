using Network.Proto;

namespace Gameplay.Player.Model
{
    public class PlayerManager
    {
        public PPlayer PlayerInfo { get; private set; }
        public bool IsLogin => PlayerInfo != null;
        public string ID => PlayerInfo.Id;
        public string Name => PlayerInfo.Nickname;
        private static PlayerManager s_Instance;
        public static PlayerManager Instance
        {
            get
            {
                if (s_Instance == null)
                    s_Instance = new PlayerManager();
                return s_Instance;
            }
        }

        public void SetPlayerInfo(PPlayer player)
        {
            PlayerInfo = player;
        }
    }
}
