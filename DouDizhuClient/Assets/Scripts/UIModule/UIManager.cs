using System.Reflection;
using UnityEngine;
using EdenMeng.AssetManager;
using System.Collections.Generic;
using System;
using Serilog;

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

        private List<ShowingUIInfo> m_ShowingUIInfos = new List<ShowingUIInfo>();

        private Dictionary<EnumUILayer, List<string>> m_UIStacks = new Dictionary<EnumUILayer, List<string>>()
        {
            { EnumUILayer.Background, new List<string>() },
            { EnumUILayer.View, new List<string>() },
            { EnumUILayer.Popup, new List<string>() },
            { EnumUILayer.Floating, new List<string>() },
            { EnumUILayer.Guide, new List<string>() },
            { EnumUILayer.System, new List<string>() },
        };

#region ShowUI
        public void ShowUI<TUIComponent>() where TUIComponent : UIComponentBase
        {
            ShowUIImpl(typeof(TUIComponent), typeof(TUIComponent).FullName, null);
        }
        
        public void ShowUI<TUIComponent>(string identifier) where TUIComponent : UIComponentBase
        {
            ShowUIImpl(typeof(TUIComponent), identifier, null);
        }

        public void ShowUI<TUIComponent, TArgs>(TArgs? args) where TUIComponent : UIComponentBase where TArgs : struct
        {
            ShowUIImpl(typeof(TUIComponent), typeof(TUIComponent).FullName, args);
        }

        public void ShowUI<TUIComponent, TArgs>(string identifier, TArgs? args) where TUIComponent : UIComponentBase where TArgs : struct
        {
            ShowUIImpl(typeof(TUIComponent), identifier, args);
        }

        public void ShowUI(Type componentType)
        {
            ShowUIImpl(componentType, componentType.FullName, null);
        }

        public void ShowUI(Type componentType, object? args)
        {
            ShowUIImpl(componentType, componentType.FullName, args);
        }

        public void ShowUI(Type componentType, string identifier)
        {
            ShowUIImpl(componentType, identifier, null);
        }

        public void ShowUI(Type componentType, string identifier, object? args)
        {
            ShowUIImpl(componentType, identifier, args);
        }

        private void ShowUIImpl(Type componentType, string identifier, object? args)
        {
            var componentAttribute = componentType.GetCustomAttribute(typeof(UIComponentAttribute)) as UIComponentAttribute;
            if (componentAttribute == null)
            {
                Log.Error("未找到{component}的属性定义", componentType.FullName);
                return;
            }
            var uiObjAsset = AssetManager.LoadAsset<GameObject>(componentAttribute.ResPath);
            var uiObj = GameObject.Instantiate(uiObjAsset);
            var uiComponent = uiObj.GetComponent(componentType) as UIComponentBase;
            if (uiComponent == null)
            {
                Log.Error("未找到{component}的UI组件", componentType.FullName);
                return;
            }
            var showingUIInfo = new ShowingUIInfo()
            {
                Identifier = identifier,
                Layer = componentAttribute.OpenLayer,
                UIComponent = uiComponent,
            };
            uiComponent.Initialize(identifier);
            m_ShowingUIInfos.Add(showingUIInfo);
            UIRoot.Instance.AppendToLayer(componentAttribute.OpenLayer, uiObj);
            Log.Information("显示UI：{component}，标识：{identifier}", componentType.FullName, identifier);

            uiComponent.OnShowBegin(args);
            uiComponent.OnShowFinish(args);
            m_UIStacks[componentAttribute.OpenLayer].Add(identifier);
        }
#endregion

#region HideUI
        public void HideUI<TUIComponent>() where TUIComponent : UIComponentBase
        {
            HideUI(typeof(TUIComponent));
        }

        public void HideUI(Type componentType)
        {
            HideUI(componentType, componentType.FullName);
        }

        public void HideUI<TUIComponent>(string identifier) where TUIComponent : UIComponentBase
        {
            HideUI(typeof(TUIComponent), identifier);
        }

        public void HideUI(Type componentType, string identifier)
        {
            var showingUIInfos = m_ShowingUIInfos.FindAll(info => info.UIComponent.GetType() == componentType && info.Identifier == identifier);
            foreach (var info in showingUIInfos)
                HideUIImpl(info);
        }
#endregion

#region HideUIType
        public void HideUIType<TUIComponent>() where TUIComponent : UIComponentBase
        {
            HideUIType(typeof(TUIComponent));
        }

        public void HideUIType(Type componentType)
        {
            var showingUITypeInfos = m_ShowingUIInfos.FindAll(info => info.UIComponent.GetType() == componentType);
            foreach (var info in showingUITypeInfos)
                HideUIImpl(info);
        }
#endregion

        private void HideUIImpl(ShowingUIInfo showingUIInfo)
        {
            Log.Information("隐藏UI：{component}，标识：{identifier}", showingUIInfo.UIComponent.GetType().FullName, showingUIInfo.Identifier);
            showingUIInfo.UIComponent.OnHideBegin();
            showingUIInfo.UIComponent.OnHideFinish();
            m_ShowingUIInfos.Remove(showingUIInfo);
            UIRoot.Instance.RemoveFromLayer(showingUIInfo.Layer, showingUIInfo.UIComponent.gameObject);
            UnityEngine.Object.Destroy(showingUIInfo.UIComponent.gameObject);
        }

#region IsUIShowing
        public bool IsUIShowing<TUIComponent>() where TUIComponent : UIComponentBase
        {
            return IsUIShowing(typeof(TUIComponent));
        }

        public bool IsUIShowing(Type componentType)
        {
            return IsUIShowing(componentType, componentType.FullName);
        }

        public bool IsUIShowing(Type componentType, string identifier)
        {
            return m_ShowingUIInfos.FindAll(info => info.UIComponent.GetType() == componentType && info.Identifier == identifier).Count > 0;
        }
#endregion

#region IsUITypeShowing
        public bool IsUITypeShowing<TUIComponent>() where TUIComponent : UIComponentBase
        {
            return IsUITypeShowing(typeof(TUIComponent));
        }

        public bool IsUITypeShowing(Type componentType)
        {
            return m_ShowingUIInfos.FindAll(info => info.UIComponent.GetType() == componentType).Count > 0;
        }
#endregion
    }
}
