using System.Threading.Tasks;
using UnityEngine;
using TMPro;
using UIModule;
using Gameplay.Player.Service;
using UnityEngine.UI;
using Config;

namespace Gameplay.Login.View
{
    [UIComponent(OpenLayer = EnumUILayer.View, ResPath = "Assets/Res/Gameplay/UI/Login/UIRegister.prefab")]
    public class UIRegister : UIComponentBase
    {
        [SerializeField]
        private TMP_InputField m_AccountInput;
        [SerializeField]
        private TextMeshProUGUI m_TextAccountTip;
        [SerializeField]
        private TMP_InputField m_PasswordInput;
        [SerializeField]
        private TextMeshProUGUI m_TextPasswordTip;

        protected override void Awake()
        {
            base.Awake();

            m_AccountInput.onSelect.AddListener(OnSelectAccountInput);
            m_AccountInput.onDeselect.AddListener(OnDeselectAccountInput);
            m_PasswordInput.onSelect.AddListener(OnSelectPasswordInput);
            m_PasswordInput.onDeselect.AddListener(OnDeselectPasswordInput);

            m_TextAccountTip.gameObject.SetActive(false);
            m_TextPasswordTip.gameObject.SetActive(false);
        }

        private void Start()
        {
            int accountMinLength = ConfigsManager.Instance.GetConst<int>("AccountMinLength");
            int accountMaxLength = ConfigsManager.Instance.GetConst<int>("AccountMaxLength");
            int passwordMinLength = ConfigsManager.Instance.GetConst<int>("PasswordMinLength");
            int passwordMaxLength = ConfigsManager.Instance.GetConst<int>("PasswordMaxLength");
            m_TextAccountTip.text = StringsHelper.GetString("AccountTip", accountMinLength, accountMaxLength);
            m_TextPasswordTip.text = StringsHelper.GetString("PasswordTip", passwordMinLength, passwordMaxLength);
        }

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

        private void OnSelectAccountInput(string _)
        {
            m_TextAccountTip.gameObject.SetActive(true);
        }

        private void OnDeselectAccountInput(string _)
        {
            m_TextAccountTip.gameObject.SetActive(false);
        }

        private void OnSelectPasswordInput(string _)
        {
            m_TextPasswordTip.gameObject.SetActive(true);
        }

        private void OnDeselectPasswordInput(string _)
        {
            m_TextPasswordTip.gameObject.SetActive(false);
        }
    }
}
