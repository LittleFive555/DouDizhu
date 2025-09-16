using Config.Define;

namespace Config
{
    public partial class ConfigsManager
    {
        public string GetString(string key)
        {
            return GetConfig<DStrings>(key)?.Value;
        }

        public string GetString(string key, params object[] args)
        {
            var str = GetString(key);
            if (str == null)
                return null;

            return string.Format(str, args);
        }
    }
}