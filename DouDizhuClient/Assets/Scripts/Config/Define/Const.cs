using System;

namespace Config.Define
{
    [Serializable]
    public class DConst : DBaseData<string>
    {
        public string Value;
    }
    
    public class DConstList
    {
        public DConst[] Content;
    }
}
