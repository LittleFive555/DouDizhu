using System.Threading.Tasks;
using UnityEngine;
using TMPro;
using UIModule;
using Gameplay.Player;
using Gameplay.Chat.View;

namespace Gameplay.Login.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Login/UILogin.prefab")]
    public class UILogin : UIComponentBase
    {
        [SerializeField]
        private TMP_InputField m_AccountInput;
        [SerializeField]
        private TMP_InputField m_PasswordInput;

        private async Task OnClickLoginAsync()
        {
            string account = m_AccountInput.text;
            string password = m_PasswordInput.text;
            bool result = await PlayerManager.Instance.Login(account, password);
            if (result)
            {
                Hide();
                UIManager.Instance.ShowUI<UIMessageSender>();
            }
        }

        [OnClick("Panel/RegisterButton")]
        private void OnClickRegister()
        {
            UIManager.Instance.ShowUI<UIRegister>();
        }

        [OnClick("Panel/LoginButton")]
        private async void OnClickLogin()
        {
            await OnClickLoginAsync();
        }
    }
}