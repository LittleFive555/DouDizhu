using System;

namespace Config.Define
{
    [Serializable]
    public class DConst : BaseData<string>
    {
        public string Value;
        public string Type;
    }
    
    public class DConstList
    {
        public DConst[] Content;
    }
}
