using UnityEngine;
using Network;
using Network.Proto;
using System;

namespace Gameplay.Character
{
    public class NetworkMoveTest : MonoBehaviour
    {
        [SerializeField]
        private float m_Speed = 2.0f;

        private StarterAssetsInputs m_Input;

        private string m_CharacterId;

        private Vector2 m_LastMove = Vector2.zero;

        public void Initialize(string characterId, StarterAssetsInputs input)
        {
            m_CharacterId = characterId;
            m_Input = input;
        }

        public void Initialize(string characterId)
        {
            m_CharacterId = characterId;
        }

        private void Update()
        {
            if (m_Input != null)
            {
                var move = m_Input.move;
                if (m_LastMove != move)
                {
                    m_LastMove = move;
                    // 发送移动消息
                    PCharacterMove msg = new()
                    {
                        CharacterId = m_CharacterId,
                        Move = new PVector3()
                        {
                            X = move.x,
                            Y = move.y
                        },
                        Timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds()
                    };
                    _ = NetworkManager.Instance.RequestAsync(PMsgId.CharacterMove, msg);
                }
                if (move == Vector2.zero)
                    return;
                transform.position += new Vector3(move.x, move.y, 0) * m_Speed * Time.deltaTime;
            }
        }
    }
}