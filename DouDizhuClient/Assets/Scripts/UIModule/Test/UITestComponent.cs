using TMPro;
using UnityEngine;

namespace UIModule.Test
{
    public class UITestComponent : MonoBehaviour
    {
        [SerializeField] private TextMeshProUGUI m_Text;
        [SerializeField] private Vector2 m_TextPosition;
        [SerializeField] private EnumUILayer m_Layer;

        private void Start()
        {
            
        }

        public void ShowTestUI(string text, Vector2 textPosition, EnumUILayer layer)
        {
            switch (layer)
            {
                case EnumUILayer.Background:
                    UIManager.Instance.ShowUI<UIBackgroundTest, UITestBase.Args>(new UITestBase.Args() { Text = text, TextPosition = textPosition });
                    break;
                case EnumUILayer.View:
                    UIManager.Instance.ShowUI<UIViewTest, UITestBase.Args>(new UITestBase.Args() { Text = text, TextPosition = textPosition });
                    break;
                case EnumUILayer.Popup:
                    UIManager.Instance.ShowUI<UIPopupTest, UITestBase.Args>(new UITestBase.Args() { Text = text, TextPosition = textPosition });
                    break;
                case EnumUILayer.Floating:
                    UIManager.Instance.ShowUI<UIFloatingTest, UITestBase.Args>(new UITestBase.Args() { Text = text, TextPosition = textPosition });
                    break;
                case EnumUILayer.Guide:
                    UIManager.Instance.ShowUI<UIGuideTest, UITestBase.Args>(new UITestBase.Args() { Text = text, TextPosition = textPosition });
                    break;
                case EnumUILayer.System:
                    UIManager.Instance.ShowUI<UISystemTest, UITestBase.Args>(new UITestBase.Args() { Text = text, TextPosition = textPosition });
                    break;
            }
        }
    }
}