using UnityEngine;

namespace Gameplay
{
    public class Launcher : MonoBehaviour
    {
        void Awake()
        {
            GameManager.Instance.Launch();
        }
    }
}
