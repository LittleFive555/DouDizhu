using Config;
using EdenMeng.AssetManager;
using Gameplay.World;
using UnityEngine;

namespace Gameplay
{
    public class PlaygroundLauncher : MonoBehaviour
    {
        [SerializeField]
        private string m_ServerHost;

        [SerializeField]
        private bool m_UseAB = false;

        private async void Awake()
        {
            InitAssetManager();
            ConfigsManager.Instance.LoadConfigs();
            await GameManager.Instance.LaunchPlayground(m_ServerHost);
            InitializeWorld();
        }

        private void InitAssetManager()
        {
#if UNITY_EDITOR
            if (m_UseAB)
                AssetManager.InitWithAssetBundle(new AssetBundlePath());
            else
                AssetManager.InitWithDatabase();
#else
            AssetManager.InitWithAssetBundle(new AssetBundlePath());
#endif
        }

        private void InitializeWorld()
        {
            RoomWorld world = FindObjectOfType<RoomWorld>();
            world.Initialize();
        }

        private class AssetBundlePath : IAssetBundlePath
        {
            public string Path => Application.streamingAssetsPath;
        }
    }
}
