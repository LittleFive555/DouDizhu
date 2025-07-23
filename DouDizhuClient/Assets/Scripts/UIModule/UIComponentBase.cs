using System;
using Serilog;
using UnityEngine;

namespace UIModule
{
    public abstract class UIComponentBase<TArgs> : UIComponentBase where TArgs : struct
    {
        public virtual void OnShowBegin(TArgs args)
        {

        }

        public sealed override void OnShowBegin(object args)
        {
            base.OnShowBegin(args);
            OnShowBegin((TArgs)args);
        }

        public virtual void OnShowFinish(TArgs args)
        {

        }

        public sealed override void OnShowFinish(object args)
        {
            base.OnShowFinish(args);
            OnShowFinish((TArgs)args);
        }
    }
    
    public abstract class UIComponentBase : UIWidget
    {
        private bool m_IsInitialized = false;
        private ShowingUIInfo m_ShowingUIInfo;
        public struct EmptyArgs { }

        public void Initialize(ShowingUIInfo showingUIInfo)
        {
            if (m_IsInitialized)
                return;
            m_ShowingUIInfo = showingUIInfo;
            m_IsInitialized = true;

            try
            {
                OnInitialize();
            }
            catch (Exception ex)
            {
                Log.Error(ex, "{showingUIInfo}在 OnInitialize() 时发生错误", m_ShowingUIInfo);
            }
        }

        public virtual void OnInitialize()
        {
            
        }

        public virtual void OnShowBegin(object args)
        {

        }

        public virtual void OnShowFinish(object args)
        {

        }

        public virtual void OnCovered()
        {

        }

        public virtual void OnUncovered()
        {
            
        }

        public virtual void OnHideBegin()
        {

        }

        public virtual void OnHideFinish()
        {

        }

        public void Hide()
        {
            UIManager.Instance.HideUI(m_ShowingUIInfo);
        }
    }
}
