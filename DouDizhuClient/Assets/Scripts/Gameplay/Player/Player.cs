using System.Threading.Tasks;
using Network.Proto;

namespace Gameplay.Player
{
    public class PlayerManager
    {
        public bool IsLogin => !string.IsNullOrEmpty(ID);
        public string ID { get; private set; }
        public string Name { get; private set; }

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

        public async Task<bool> Register(string account, string password)
        {
            var response = await Network.NetworkManager.Instance.RequestAsync(PMsgId.Register, new PRegisterRequest()
            {
                Account = account,
                Password = password
            });
            if (!response.IsSuccess)
                return false;
            return true;
        }
        
        public async Task<bool> Login(string account, string password)
        {
            var response = await Network.NetworkManager.Instance.RequestAsync<PLoginRequest, PLoginResponse>(PMsgId.Login, new PLoginRequest()
            {
                Account = account,
                Password = password
            });
            if (!response.IsSuccess)
                return false;

            ID = response.Data.PlayerId;
            return true;
        }
    }
}
