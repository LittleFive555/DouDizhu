﻿using UnityEngine;

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
        public string Identifier { get; private set; }
        public struct EmptyArgs { }

        public virtual void Initialize(string identifier)
        {
            if (m_IsInitialized)
                return;
            Identifier = identifier;
            m_IsInitialized = true;
        }

        public virtual void OnShowBegin(object args)
        {

        }

        public virtual void OnShowFinish(object args)
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
            UIManager.Instance.HideUI(GetType(), Identifier);
        }
    }
}
