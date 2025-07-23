using System;

namespace UIModule
{
    public class ShowingUIInfo
    {
        public readonly string Identifier;
        public readonly EnumUILayer Layer;
        public readonly UIComponentBase UIComponent;
        public readonly long ShowTime;

        public ShowingUIInfo(string identifier, EnumUILayer layer, UIComponentBase uiComponent)
        {
            Identifier = identifier;
            ShowTime = DateTime.Now.Ticks;
            Layer = layer;
            UIComponent = uiComponent;
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
