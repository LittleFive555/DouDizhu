using System.Threading.Tasks;
using Gameplay.Player;
using TMPro;
using UIModule;
using UnityEngine;

namespace Gameplay.Login.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Login/UIRegister.prefab")]
    public class UIRegister : UIComponentBase
    {
        [SerializeField]
        private TMP_InputField m_AccountInput;
        [SerializeField]
        private TMP_InputField m_PasswordInput;

        public void OnClickConfirm()
        {
            OnClickConfirmAsync();
        }

        public async Task OnClickConfirmAsync()
        {
            string account = m_AccountInput.text;
            string password = m_PasswordInput.text;
            bool result = await PlayerManager.Instance.Register(account, password);
            if (result)
                Hide();
        }
    }
}
