using System;
using Serilog;

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

        private int m_CoveredCounter = 0;
        
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
                Log.Error(ex, $"{m_ShowingUIInfo}在 OnInitialize() 时发生错误");
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

        public void Covered()
        {
            if (m_CoveredCounter == 0)
            {
                try
                {
                    OnCovered();
                }
                catch (Exception e)
                {
                    Log.Error(e, $"{m_ShowingUIInfo}在 OnCovered() 时发生错误");
                }
                finally
                {
                    gameObject.SetActive(false);
                }
            }
            m_CoveredCounter++;
        }

        public virtual void OnCovered()
        {

        }

        public void Uncovered()
        {
            if (m_CoveredCounter == 0)
                return;
            m_CoveredCounter--;
            if (m_CoveredCounter == 0)
            {
                try
                {
                    OnUncovered();
                }
                catch (Exception e)
                {
                    Log.Error(e, $"{m_ShowingUIInfo}在 OnUncovered() 时发生错误");
                }
                finally
                {
                    gameObject.SetActive(true);
                }
            }
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
