using Serilog;
using UnityEngine;

namespace Gameplay
{
    public class Launcher : MonoBehaviour
    {
        [SerializeField]
        private string m_ServerHost;

        void Awake()
        {
            GameManager.Instance.Launch(m_ServerHost);
        }
    }
}
