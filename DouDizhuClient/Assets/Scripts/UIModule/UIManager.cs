using System.Reflection;
using UnityEngine;
using EdenMeng.AssetManager;
using System.Collections.Generic;
using System;
using Serilog;
using System.Text;

namespace UIModule
{
    public enum EnumUILayer
    {
        /// <summary>
        /// 背景层
        /// </summary>
        Background,
        /// <summary>
        /// 全屏视图层，打开后会隐藏View和Popup层
        /// </summary>
        View,
        /// <summary>
        /// 弹窗层，打开后会在View层之上，但不隐藏View层
        /// </summary>
        Popup,
        /// <summary>
        /// 浮动层，打开后始终在View和Popup层之上
        /// </summary>
        Floating,
        /// <summary>
        /// 引导层，用于显示引导UI，始终在View、Popup和Floating层之上
        /// </summary>
        Guide,
        /// <summary>
        /// 系统层，层级最高，用于显示系统提示等，在所有其他层之上
        /// </summary>
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

        public void ShowUI(Type componentType, object args)
        {
            ShowUIImpl(componentType, componentType.FullName, args);
        }

        public void ShowUI(Type componentType, string identifier)
        {
            ShowUIImpl(componentType, identifier, null);
        }

        public void ShowUI(Type componentType, string identifier, object args)
        {
            ShowUIImpl(componentType, identifier, args);
        }

        private void ShowUIImpl(Type componentType, string identifier, object args)
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
            var showingUIInfo = new ShowingUIInfo(identifier, componentAttribute.OpenLayer, uiComponent);
            uiComponent.Initialize(showingUIInfo);

            m_ShowingUIInfos.Add(showingUIInfo);
            UIRoot.Instance.AppendToLayer(componentAttribute.OpenLayer, uiObj);
            m_UIStacks[componentAttribute.OpenLayer].Add(identifier);
            Log.Information("显示UI：{component}, 堆栈：{stack}\n 详细堆栈：{stackAll}", componentType, GetUIStack(), GetUIStackAllLayers());

            try
            {
                uiComponent.OnShowBegin(args);
            }
            catch (Exception ex)
            {
                Log.Error(ex, "{showingUIInfo}在 OnShowBegin() 时发生错误", showingUIInfo);
            }

            try
            {
                uiComponent.OnShowFinish(args);
            }
            catch (Exception ex)
            {
                Log.Error(ex, "{showingUIInfo}在 OnShowFinish() 时发生错误", showingUIInfo);
            }

            AfterShowUI(showingUIInfo);
        }

        private void AfterShowUI(ShowingUIInfo uiInfo)
        {
            if (uiInfo.Layer == EnumUILayer.View)
            {
                int currentIndex = m_ShowingUIInfos.IndexOf(uiInfo);
                for (int i = 0; i < currentIndex; i++)
                {
                    var showingUIInfo = m_ShowingUIInfos[i];
                    if (showingUIInfo.Layer != EnumUILayer.View && showingUIInfo.Layer != EnumUILayer.Popup)
                        continue;
                    showingUIInfo.UIComponent.Covered();
                }
            }
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
            var showingUIInfo = m_ShowingUIInfos.FindLast(info => info.UIComponent.GetType() == componentType && info.Identifier == identifier);
            if (showingUIInfo == null)
            {
                Log.Warning("未找到UI：<{component}>[{identifier}]，当前堆栈：{stack}", componentType.FullName, identifier, GetUIStack());
                return;
            }
            HideUI(showingUIInfo);
        }

        public void HideUI(ShowingUIInfo showingUIInfo)
        {
            HideUIImpl(showingUIInfo);
        }
        
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

        private void HideUIImpl(ShowingUIInfo showingUIInfo)
        {
            if (showingUIInfo == null)
                return;
                
            BeforeHideUI(showingUIInfo);

            m_ShowingUIInfos.Remove(showingUIInfo);
            UIRoot.Instance.RemoveFromLayer(showingUIInfo.Layer, showingUIInfo.UIComponent.gameObject);
            m_UIStacks[showingUIInfo.Layer].Remove(showingUIInfo.Identifier);

            Log.Information("隐藏UI：{component}, 堆栈：{stack}\n 详细堆栈：{stackAll}", showingUIInfo, GetUIStack(), GetUIStackAllLayers());

            try
            {
                showingUIInfo.UIComponent.OnHideBegin();
            }
            catch (Exception ex)
            {
                Log.Error(ex, "{showingUIInfo}在 OnHideBegin() 时发生错误", showingUIInfo);
            }

            try
            {
                showingUIInfo.UIComponent.OnHideFinish();
            }
            catch (Exception ex)
            {
                Log.Error(ex, "{showingUIInfo}在 OnHideFinish() 时发生错误", showingUIInfo);
            }
            UnityEngine.Object.Destroy(showingUIInfo.UIComponent.gameObject);
        }

        private void BeforeHideUI(ShowingUIInfo uiInfo)
        {
            if (uiInfo.Layer == EnumUILayer.View)
            {
                int currentIndex = m_ShowingUIInfos.IndexOf(uiInfo);
                for (int i = 0; i < currentIndex; i++)
                {
                    var showingUIInfo = m_ShowingUIInfos[i];
                    if (showingUIInfo.Layer != EnumUILayer.View && showingUIInfo.Layer != EnumUILayer.Popup)
                        continue;
                    showingUIInfo.UIComponent.Uncovered();
                }
            }
        }
#endregion

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

        public string GetUIStack()
        {
            return string.Join(", ", m_ShowingUIInfos);
        }

        public string GetUIStackByLayer(EnumUILayer layer)
        {
            return string.Join(", ", m_UIStacks[layer]);
        }

        public string GetUIStackAllLayers()
        {
            StringBuilder stack = new StringBuilder();
            foreach (var layer in m_UIStacks)
            {
                stack.Append($"{layer.Key}: ");
                stack.Append(string.Join(", ", layer.Value));
                stack.Append("\n");
            }
            return stack.ToString();
        }
    }
}
