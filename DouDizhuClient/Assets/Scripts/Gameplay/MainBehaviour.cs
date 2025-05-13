using UnityEngine;

namespace Gameplay
{
    public class MainBehaviour : MonoBehaviour
    {
        private void Update()
        {
            
        }

        private void LateUpdate()
        {
            
        }

        private void FixedUpdate()
        {

        }
        
        private void OnApplicationQuit()
        {
            GameManager.Instance.Exit();
        }
    }
}
