using UnityEngine;
using UnityEngine.UI;

namespace Gameplay.Character
{
    public class NetworkMoveTest : MonoBehaviour
    {
        [SerializeField]
        private StarterAssetsInputs m_Input;

        private float Speed = 2.0f;

        private bool m_IsClient = true;

        private void Update()
        {
            if (m_IsClient)
            {
                var move = m_Input.move;
                if (move == Vector2.zero)
                    return;
                transform.position += new Vector3(move.x, move.y, 0) * Speed * Time.deltaTime;
            }
        }
    }
}