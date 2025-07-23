using TMPro;
using UnityEngine;

namespace UIModule.Test
{
    public class UITestBase : UIComponentBase<UITestBase.Args>
    {
        [SerializeField]
        private TextMeshProUGUI m_Text;

        public struct Args
        {
            public string Text;
            public Vector2 TextPosition;
        }

        public override void OnShowBegin(Args args)
        {
            base.OnShowBegin(args);

            m_Text.text = args.Text;
            m_Text.transform.localPosition = args.TextPosition;
        }

        [OnClick("BtnClose")]
        protected void OnClickClose()
        {
            Hide();
        }
    }
}