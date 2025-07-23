using System;
using System.Collections.Generic;
using System.IO;
using System.Reflection;
using Config.Define;
using Serilog;
using UnityEngine;

namespace Config
{
    public class ConfigsManager
    {
        private static ConfigsManager m_Instance;
        public static ConfigsManager Instance
        {
            get
            {
                if (m_Instance == null)
                    m_Instance = new ConfigsManager();
                return m_Instance;
            }
        }

        private const string ConfigClassNamespace = "Config.Define";
        private const string ConfigDataPath = "Configs";

        private Dictionary<Type, Dictionary<int, DBaseData<int>>> m_ConfigDictInt = new Dictionary<Type, Dictionary<int, DBaseData<int>>>();
        private Dictionary<Type, Dictionary<string, DBaseData<string>>> m_ConfigDictString = new Dictionary<Type, Dictionary<string, DBaseData<string>>>();

        public void LoadConfigs()
        {
            Assembly assembly = Assembly.GetExecutingAssembly();
            Type[] types = assembly.GetTypes();
            foreach (var type in types)
            {
                if (type.IsSubclassOf(typeof(DBaseData<int>)))
                    m_ConfigDictInt.Add(type, LoadConfig<int>(type));
                else if (type.IsSubclassOf(typeof(DBaseData<string>)))
                    m_ConfigDictString.Add(type, LoadConfig<string>(type));
            }
        }

        public TValue GetConst<TValue>(string id)
        {
            var config = GetConfig<DConst>(id);
            if (config == null)
                return default;

            return (TValue)Convert.ChangeType(config.Value, typeof(TValue));
        }

        public T GetConfig<T>(int id) where T : DBaseData<int>
        {
            if (m_ConfigDictInt.TryGetValue(typeof(T), out var dictInt) && dictInt.TryGetValue(id, out var config))
                return (T)config;

            Log.Error($"ConfigsManager: Config {typeof(T).Name} with id {id} not found");
            return null;
        }

        public T GetConfig<T>(string id) where T : DBaseData<string>
        {
            if (string.IsNullOrEmpty(id))
                return null;

            if (m_ConfigDictString.TryGetValue(typeof(T), out var dictString) && dictString.TryGetValue(id, out var config))
                return (T)config;

            Log.Error($"ConfigsManager: Config {typeof(T).Name} with id {id} not found");
            return null;
        }

        private static string ReadRawText(string fileName)
        {
            string readData;
            string fileFullPath = Path.Combine(Application.streamingAssetsPath, ConfigDataPath, fileName);
            using (StreamReader sr = File.OpenText(fileFullPath))
            {
                readData = sr.ReadToEnd();
                sr.Close();
            }
            return readData;
        }

        private Dictionary<TIndex, DBaseData<TIndex>> LoadConfig<TIndex>(Type type)
        {
            Dictionary<TIndex, DBaseData<TIndex>> dict = new Dictionary<TIndex, DBaseData<TIndex>>();
            string typeName = type.Name;
            string fileName = typeName.Substring(1);
            string json = ReadRawText(fileName);
            var listType = Type.GetType($"{ConfigClassNamespace}.{typeName}List");
            var obj = JsonUtility.FromJson(json, listType);
            var field = listType.GetField("Content");
            DBaseData<TIndex>[] list = (DBaseData<TIndex>[])field.GetValue(obj);
            foreach (var item in list)
            {
                dict.Add(item.ID, item);
            }
            return dict;
        }
    }
}