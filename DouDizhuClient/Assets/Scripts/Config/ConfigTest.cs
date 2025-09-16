using System;
using Serilog;
using UnityEngine;

namespace Config.Editor
{
    public class ConfigTest : MonoBehaviour
    {
        private void Start()
        {
            ConfigsManager.Instance.LoadConfigs();
            TestConst();
        }

        private void TestConst()
        {
            TestReadConst<int>("TestInt");
            TestReadConst<float>("TestFloat");
            TestReadConst<bool>("TestBool");
            TestReadConst<long>("TestLong");
            TestReadConst<double>("TestDouble");
            TestReadConst<string>("TestString");
            TestReadConst<Vector2>("TestVector2");
            TestReadConst<Vector3>("TestVector3");
            TestReadConst<DateTime>("TestDateTime");
        }

        private void TestReadConst<T>(string id)
        {
            try
            {
                var value = ConfigsManager.Instance.GetConst<T>(id);
                Log.Information($"TestReadConst<{typeof(T).Name}> with id {id} success, value: {value}");

            }
            catch (Exception ex)
            {
                Log.Error(ex, $"TestReadConst<{typeof(T).Name}> with id {id} failed");
            }
        }
    }
}