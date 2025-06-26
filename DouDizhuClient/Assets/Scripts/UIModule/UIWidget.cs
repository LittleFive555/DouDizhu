using System.Reflection;
using UnityEngine;
using UnityEngine.UI;

namespace UIModule
{
    public class UIWidget : MonoBehaviour
    {
        protected virtual void Awake()
        {
            BindClickMethods();
        }

        private void BindClickMethods()
        {
            var methods = GetType().GetMethods(BindingFlags.Instance | BindingFlags.NonPublic | BindingFlags.Public);
            foreach (var method in methods)
            {
                var attribute = method.GetCustomAttribute<OnClickAttribute>();
                if (attribute == null)
                    continue;

                var button = transform.Find(attribute.Path)?.GetComponent<Button>();
                if (button == null)
                    continue;

                button.onClick.AddListener(() => method.Invoke(this, null));
            }
        }
    }
}