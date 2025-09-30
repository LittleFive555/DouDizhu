using System;
using System.Threading.Tasks;
using UnityEngine;
using UnityEngine.SceneManagement;
using Network;
using UIModule;
using Gameplay.Login.View;
using Serilog;
using Serilog.Sinks.Unity3D;

namespace Gameplay
{
    public class GameManager
    {
        private static GameManager m_Instance;
        public static GameManager Instance
        {
            get
            {
                if (m_Instance == null)
                    m_Instance = new GameManager();
                return m_Instance;
            }
        }

        [RuntimeInitializeOnLoadMethod(RuntimeInitializeLoadType.AfterAssembliesLoaded)]
        public static void InitializeLogger()
        {
            // TODO 日志文件在Editor下和Build下分开
            // 初始化日志
            Log.Logger = new LoggerConfiguration()
                .WriteTo.Unity3D()
                .WriteTo.File("log.txt",
                    rollingInterval: RollingInterval.Day,
                    rollOnFileSizeLimit: true)
                .CreateLogger();
            Log.Information("Logger Initialized");
        }

        public async Task Launch(string serverHost)
        {
            CreateMainBehaviour();

            await NetworkManager.Instance.ConnectAsync(serverHost, 15770);

            try
            {
                await NetworkManager.Instance.Handshake();
            }
            catch (Exception ex)
            {
                Log.Error(ex, "握手失败");
                return;
            }

            Log.Information("Start Loading Main Scene");
            SceneManager.LoadSceneAsync("MainScene", LoadSceneMode.Single).completed += (AsyncOperation obj) =>
            {
                Log.Information("Main Scene Loaded");
                UIManager.Instance.ShowUI<UILogin>();
            };
        }

        public async Task LaunchPlayground(string serverHost)
        {

            CreateMainBehaviour();

            await NetworkManager.Instance.ConnectAsync(serverHost, 15771);
        }

        public void Exit()
        {
            NetworkManager.Instance.Disconnect();
            Log.CloseAndFlush();
        }

        private void CreateMainBehaviour()
        {
            GameObject mainBehaviour = new GameObject("MainBehaviour");
            mainBehaviour.AddComponent<MainBehaviour>();
            GameObject.DontDestroyOnLoad(mainBehaviour);
        }
    }
}