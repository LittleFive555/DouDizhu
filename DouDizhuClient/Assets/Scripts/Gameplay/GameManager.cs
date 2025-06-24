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

        public async Task Launch(string serverHost)
        {
            InitializeLogger();
            Log.Information("Logger Initialized");

            CreateMainBehaviour();

            await NetworkManager.Instance.ConnectAsync(serverHost);

            // await NetworkManager.Instance.Handshake();

            Log.Information("Start Loading Main Scene");
            SceneManager.LoadSceneAsync("MainScene", LoadSceneMode.Single).completed += (AsyncOperation obj) =>
            {
                Log.Information("Main Scene Loaded");
                UIManager.Instance.ShowUI<UILogin>();
            };
        }

        public void Exit()
        {
            NetworkManager.Instance.Disconnect();
            Log.CloseAndFlush();
        }

        private void InitializeLogger()
        {
            // TODO 日志文件在Editor下和Build下分开
            // 初始化日志
            Log.Logger = new LoggerConfiguration()
                .WriteTo.Unity3D()
                .WriteTo.File("log.txt",
                    rollingInterval: RollingInterval.Day,
                    rollOnFileSizeLimit: true)
                .CreateLogger();
        }

        private void CreateMainBehaviour()
        {
            GameObject mainBehaviour = new GameObject("MainBehaviour");
            mainBehaviour.AddComponent<MainBehaviour>();
            Object.DontDestroyOnLoad(mainBehaviour);
        }
    }
}