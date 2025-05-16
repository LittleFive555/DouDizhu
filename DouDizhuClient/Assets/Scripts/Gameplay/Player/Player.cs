using System.Threading.Tasks;
using Network.Proto;
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
            if (response.IsSuccess)
            {
                if (response.Data.Result != PRegisterResponse.Types.Result.Success)
                {
                    Debug.LogError("注册失败，错误码：" + response.Data.Result);
                    return false;
                }
                return true;
            }
            return false;
        }
        public async Task<bool> Login(string account, string password)
        {
            var response = await Network.NetworkManager.Instance.RequestAsync<PLoginRequest, PLoginResponse>(PGameClientMessage.ContentOneofCase.LoginReq, new PLoginRequest()
            {
                Account = account,
                Password = password
            });
            if (response.IsSuccess)
            {
                if (response.Data.Result != PLoginResponse.Types.Result.Success)
                {
                    Debug.LogError("登录失败，错误码：" + response.Data.Result);
                    return false;
                }

                ID = response.Data.PlayerId;
                return true;
            }
            return false;
        }
    }
}
