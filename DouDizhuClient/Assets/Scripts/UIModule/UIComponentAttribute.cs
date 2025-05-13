using System;

namespace UIModule
{
    public class UIComponentAttribute : Attribute
    {
        public string ResPath;
        public EnumUILayer OpenLayer;
    }
}
