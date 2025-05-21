using System.Threading.Tasks;
using UnityEngine;
using UnityEngine.SceneManagement;
using Network;
using UIModule;
using Gameplay.Login.View;
using Serilog;

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
            // 初始化日志
            Log.Logger = new LoggerConfiguration()
                .WriteTo.Console()
                .WriteTo.File("log.txt",
                    rollingInterval: RollingInterval.Day,
                    rollOnFileSizeLimit: true)
                .CreateLogger();

            CreateMainBehaviour();
            await NetworkManager.Instance.ConnectAsync(serverHost);

            Debug.Log("GameManager Launch");

            SceneManager.LoadSceneAsync("MainScene", LoadSceneMode.Single).completed += (AsyncOperation obj) =>
            {
                UIManager.Instance.ShowUI<UILogin>();
            };
        }

        public void Exit()
        {
            NetworkManager.Instance.Disconnect();
        }

        private void CreateMainBehaviour()
        {
            GameObject mainBehaviour = new GameObject("MainBehaviour");
            mainBehaviour.AddComponent<MainBehaviour>();
            UnityEngine.Object.DontDestroyOnLoad(mainBehaviour);
        }
    }
}