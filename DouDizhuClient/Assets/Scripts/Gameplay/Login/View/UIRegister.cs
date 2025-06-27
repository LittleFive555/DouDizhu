using System.Threading.Tasks;
using UnityEngine;
using TMPro;
using UIModule;
using Gameplay.Player.Service;

namespace Gameplay.Login.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Login/UIRegister.prefab")]
    public class UIRegister : UIComponentBase
    {
        [SerializeField]
        private TMP_InputField m_AccountInput;
        [SerializeField]
        private TMP_InputField m_PasswordInput;

        public async Task OnClickConfirmAsync()
        {
            string account = m_AccountInput.text;
            string password = m_PasswordInput.text;
            bool result = await PlayerService.Register(account, password);
            if (result)
                Hide();
        }

        [OnClick("Button")]
        private async void OnClickConfirm()
        {
            await OnClickConfirmAsync();
        }
    }
}
