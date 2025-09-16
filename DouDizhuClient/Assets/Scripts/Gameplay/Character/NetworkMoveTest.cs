using Network;
using Network.Proto;
using UnityEngine;

namespace Gameplay.Character
{
    public class NetworkMoveTest : MonoBehaviour
    {
        [SerializeField]
        private StarterAssetsInputs m_Input;

        public float Speed = 2f;

        private Vector2 m_MoveInput = Vector2.zero;

        private void Update()
        {
            var move = m_Input.move;
            if (m_MoveInput != move)
            {
                m_MoveInput = move;
                _ = NetworkManager.Instance.RequestAsync(PMsgId.CharacterMove, new PCharacterMove()
                {
                    MoveX = m_MoveInput.x
                });
            }

            float deltaX = m_MoveInput.x * Speed * Time.deltaTime;
            if (transform.position.x + deltaX > 5)
                deltaX = 5 - transform.position.x;
            else if (transform.position.x + deltaX < -5)
                deltaX = -5 - transform.position.x;
            transform.Translate(Vector3.right * deltaX);
        }
    }
}