using System;
using System.Collections.Generic;
using Config.Define;
using UnityEngine;

namespace Config
{
    public partial class ConfigsManager
    {
        Dictionary<Type, Func<string, object>> m_ConvertDict = new Dictionary<Type, Func<string, object>>()
        {
            { typeof(Vector2), ConvertVector2 },
            { typeof(Vector3), ConvertVector3 },
        };

        public TValue GetConst<TValue>(string id)
        {
            var config = GetConfig<DConst>(id);
            if (config == null)
                throw new KeyNotFoundException($"Const {id} is not found");

            var type = typeof(TValue);
            var convertibleType = typeof(IConvertible);
            if (convertibleType.IsAssignableFrom(type))
                return (TValue)Convert.ChangeType(config.Value, typeof(TValue));
            else // 自定义转换
            {
                if (m_ConvertDict.TryGetValue(type, out var convert))
                    return (TValue)convert(config.Value);
                else
                    throw new InvalidCastException($"Type {type} is not supported");
            }
        }
        
        private static object ConvertVector2(string value)
        {
            string[] values = value.TrimStart('(').TrimEnd(')').Split(',');
            if (values.Length != 2)
                throw new InvalidCastException($"\"{value}\" is not valid Vector2");
            return new Vector2(float.Parse(values[0]), float.Parse(values[1]));
        }

        private static object ConvertVector3(string value)
        {
            string[] values = value.TrimStart('(').TrimEnd(')').Split(',');
            if (values.Length != 3)
                throw new InvalidCastException($"\"{value}\" is not valid Vector3");
            return new Vector3(float.Parse(values[0]), float.Parse(values[1]), float.Parse(values[2]));
        }
    }
}