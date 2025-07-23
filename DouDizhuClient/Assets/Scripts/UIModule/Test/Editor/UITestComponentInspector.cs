using UnityEditor;
using UnityEngine.UIElements;

namespace UIModule.Test.Editor
{
    [CustomEditor(typeof(UITestComponent))]
    public class UITestComponentInspector : UnityEditor.Editor
    {
        public override VisualElement CreateInspectorGUI()
        {
            var root = new VisualElement();
            var textElement = new TextField("Text")
            {
                bindingPath = "m_Text"
            };
            root.Add(textElement);
            var textPositionElement = new Vector2Field("Text Position")
            {
                bindingPath = "m_TextPosition"
            };
            root.Add(textPositionElement);
            var layerElement = new EnumField("Layer")
            {
                bindingPath = "m_Layer"
            };
            root.Add(layerElement);
            var button = new Button(() =>
            {
                ((UITestComponent)target).ShowTestUI(textElement.value, textPositionElement.value, (EnumUILayer)layerElement.value);
            })
            {
                text = "Show Test UI"
            };
            root.Add(button);
            return root;
        }
    }
}