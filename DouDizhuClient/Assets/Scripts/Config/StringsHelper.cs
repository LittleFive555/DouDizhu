using Config.Define;

namespace Config
{
    public static class StringsHelper
    {
        public static string GetString(string key)
        {
            return ConfigsManager.Instance.GetConfig<DStrings>(key)?.Value;
        }

        public static string GetString(string key, params object[] args)
        {
            var str = GetString(key);
            if (str == null)
                return null;

            return string.Format(str, args);
        }
    }
}