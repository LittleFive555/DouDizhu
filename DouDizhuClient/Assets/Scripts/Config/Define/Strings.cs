using System;

namespace Config.Define
{
    [Serializable]
    public class DStrings : DBaseData<string>
    {
        public string Value;
    }
    
    public class DStringsList
    {
        public DStrings[] Content;
    }
}
