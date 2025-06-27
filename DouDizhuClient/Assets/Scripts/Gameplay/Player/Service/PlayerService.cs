using System.Threading.Tasks;
using Gameplay.Player.Model;
using Network.Proto;

namespace Gameplay.Player.Service
{
    public class PlayerService
    {

        public static async Task<bool> Register(string account, string password)
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
        
        public static async Task<bool> Login(string account, string password)
        {
            var response = await Network.NetworkManager.Instance.RequestAsync<PLoginRequest, PLoginResponse>(PMsgId.Login, new PLoginRequest()
            {
                Account = account,
                Password = password
            });
            if (!response.IsSuccess)
                return false;

            PlayerManager.Instance.SetPlayerInfo(response.Data.Info);
            return true;
        }
    }
}