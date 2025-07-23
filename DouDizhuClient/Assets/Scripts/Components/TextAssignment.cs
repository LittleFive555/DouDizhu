using Config;
using Config.Define;
using Serilog;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

namespace Components
{
    public class TextAssignment : MonoBehaviour
    {
        [SerializeField]
        private string m_TextKey;

        private TMP_Text m_TMPText;
        private Text m_Text;

        private void Awake()
        {
            if (TryGetComponent(out m_TMPText))
                return;
            if (TryGetComponent(out m_Text))
                return;
        }

        private void Start()
        {
            var config = ConfigsManager.Instance.GetConfig<DStrings>(m_TextKey);
            if (config == null)
                return;

            AssignText(config.Value);
        }

        public void AssignText(string text)
        {
            if (m_TMPText != null)
                m_TMPText.text = text;
            else if (m_Text != null)
                m_Text.text = text;
        }
    }
}