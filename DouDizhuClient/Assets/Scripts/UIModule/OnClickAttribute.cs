using System;

namespace UIModule
{
    [AttributeUsage(AttributeTargets.Method, AllowMultiple = false, Inherited = true)]
    public class OnClickAttribute : Attribute
    {
        public string Path;

        public OnClickAttribute(string path)
        {
            Path = path;
        }
    }
}