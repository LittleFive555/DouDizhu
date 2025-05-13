using System.Reflection;
using UnityEngine;

namespace UIModule
{
    public enum EnumUILayer
    {
        Background,
        View,
        Popup,
        Floating,
        Guide,
        System,
    }

    public class UIManager
    {
        private static UIManager m_Instance;
        public static UIManager Instance
        {
            get
            {
                if (m_Instance == null)
                    m_Instance = new UIManager();
                return m_Instance;
            }
        }

        public void ShowUI<TUIComponent>() where TUIComponent : UIComponentBase
        {
            var componentType = typeof(TUIComponent);
            var componentAttribute = componentType.GetCustomAttribute(typeof(UIComponentAttribute)) as UIComponentAttribute;
            if (componentAttribute == null)
            {
                Debug.LogError($"未找到{componentType.FullName}的属性定义");
                return;
            }
        }

        public void ShowUI<TUIComponent, TArgs>(TArgs args) where TUIComponent : UIComponentBase where TArgs : struct
        {

        }

        public void HideUI<TUIComponent>() where TUIComponent : UIComponentBase
        {

        }

        public void IsUIShowing<TUIComponent>() where TUIComponent : UIComponentBase
        {

        }
    }
}
