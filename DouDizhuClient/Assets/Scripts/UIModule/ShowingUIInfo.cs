using System;
using Serilog;

namespace UIModule
{
    public class ShowingUIInfo
    {
        public readonly string Identifier;
        public readonly EnumUILayer Layer;
        public readonly UIComponentBase UIComponent;
        public readonly long ShowTime;

        private int m_CoveredCounter;

        public ShowingUIInfo(string identifier, EnumUILayer layer, UIComponentBase uiComponent)
        {
            Identifier = identifier;
            ShowTime = DateTime.Now.Ticks;
            Layer = layer;
            UIComponent = uiComponent;

            m_CoveredCounter = 0;
        }

        public void Covered()
        {
            if (m_CoveredCounter == 0)
            {
                try
                {
                    UIComponent.OnCovered();
                    UIComponent.gameObject.SetActive(false);
                }
                catch (Exception e)
                {
                    // TODO 是不是可以考虑直接关闭界面？
                    Log.Error(e, $"UIComponent.OnCovered() 异常, Identifier: {Identifier}, Type: {UIComponent.GetType().FullName}");
                }
            }
            m_CoveredCounter++;
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
                    UIComponent.gameObject.SetActive(true);
                    UIComponent.OnUncovered();
                }
                catch (Exception e)
                {
                    // TODO 是不是可以考虑直接关闭界面？
                    Log.Error(e, $"UIComponent.OnUncovered() 异常, Identifier: {Identifier}, Type: {UIComponent.GetType().FullName}");
                }
            }
        }

        public override bool Equals(object obj)
        {
            if (obj is ShowingUIInfo other)
            {
                return Identifier == other.Identifier && Layer == other.Layer && ShowTime == other.ShowTime && UIComponent.GetType() == other.UIComponent.GetType();
            }
            return false;
        }

        public override int GetHashCode()
        {
            return Identifier.GetHashCode() ^ ShowTime.GetHashCode() ^ Layer.GetHashCode() ^ UIComponent.GetType().GetHashCode();
        }

        public override string ToString()
        {
            return $"<{UIComponent.GetType().Name}>[{Identifier}]";
        }
    }
}
