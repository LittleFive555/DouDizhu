using Serilog;
using UnityEngine;

namespace Gameplay
{
    public class Launcher : MonoBehaviour
    {
        [SerializeField]
        private string m_ServerHost;

        private async void Awake()
        {
            await GameManager.Instance.Launch(m_ServerHost);
        }
    }
}
