using System.Reflection;
using UnityEditor;
using UnityEditor.UIElements;
using UnityEngine;
using UnityEngine.UIElements;

namespace UIModule.Editor
{
    [CustomEditor(typeof(UIWidget), true)]
    public class UIWidgetInspector : UnityEditor.Editor
    {
        public override VisualElement CreateInspectorGUI()
        {
            var root = new VisualElement();

            InspectorElement.FillDefaultInspector(root, serializedObject, this);

            var uiWidget = target as UIWidget;
            if (uiWidget == null)
                return root;

            // 按钮绑定
            var clickBindingFoldout = new Foldout()
            {
                text = "Click Binding",
                value = true
            };

            bool hasClickBinding = false;
            var type = uiWidget.GetType();
            while (type != null && type != typeof(UIWidget))
            {
                var methods = type.GetMethods(BindingFlags.Instance | BindingFlags.NonPublic | BindingFlags.Public | BindingFlags.DeclaredOnly);
                foreach (var method in methods)
                {
                    var attribute = method.GetCustomAttribute<OnClickAttribute>();
                    if (attribute == null)
                        continue;

                    var objectField = new ObjectField(method.Name);
                    objectField.SetEnabled(false);
                    objectField.value = uiWidget.transform.Find(attribute.Path)?.GetComponent<UnityEngine.UI.Button>();
                    clickBindingFoldout.Add(objectField);
                    if (objectField.value == null)
                        clickBindingFoldout.Add(new Label("Above ClickBinding is not referenced")
                        {
                            style = { color = new Color(1f, 1f, 0, 1), backgroundColor = new Color(1f, 1f, 0, 0.1f) }
                        });
                    hasClickBinding = true;
                }
                type = type.BaseType;
            }

            if (hasClickBinding)
                root.Add(clickBindingFoldout);
            
            return root;
        }
    }
}