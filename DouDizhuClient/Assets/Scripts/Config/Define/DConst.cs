using System;

namespace Config.Define
{
    [Serializable]
    public class DConst : BaseData<string>
    {
        public string Value;
    }
    
    public class DConstList
    {
        public DConst[] Content;
    }
}
