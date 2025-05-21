using System.Threading.Tasks;
using Network.Proto;
using Serilog;
using UnityEngine;

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
            var response = await Network.NetworkManager.Instance.RequestAsync<PRegisterRequest, PRegisterResponse>(PGameClientMessage.ContentOneofCase.RegisterReq, new PRegisterRequest()
            {
                Account = account,
                Password = password
            });
            if (response == null)
                return false;

            if (response.Result != PRegisterResponse.Types.Result.Success)
            {
                Log.Error("注册失败，错误码：{result}", response.Result);
                return false;
            }
            return true;
        }
        
        public async Task<bool> Login(string account, string password)
        {
            var response = await Network.NetworkManager.Instance.RequestAsync<PLoginRequest, PLoginResponse>(PGameClientMessage.ContentOneofCase.LoginReq, new PLoginRequest()
            {
                Account = account,
                Password = password
            });
            if (response == null)
                return false;

            if (response.Result != PLoginResponse.Types.Result.Success)
            {
                Log.Error("登录失败，错误码：{result}", response.Result);
                return false;
            }

            ID = response.PlayerId;
            return true;
        }
    }
}
