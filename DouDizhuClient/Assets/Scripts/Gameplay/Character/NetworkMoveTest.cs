using System;
using System.Collections.Generic;
using UnityEngine;
using Network;
using Network.Proto;
using Serilog;

namespace Gameplay.Character
{
    public class NetworkMoveTest : MonoBehaviour
    {
        [SerializeField]
        private float m_Speed = 2.0f;

        private StarterAssetsInputs m_Input;

        private string m_CharacterId;

        private List<InputData> m_InputHistory = new List<InputData>();

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
                // 预测移动
                var lastInput = PeekLastInput();
                if (lastInput != null && lastInput.Move != Vector2.zero)
                {
                    var now = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds();
                    transform.position = GetPositionAtTime(now);
                }

                // 在帧尾收集输入变化，主要是为了计算过程与服务器同步
                var input = GetInputChange();
                if (input != null)
                {
                    m_InputHistory.Add(input);

                    // 发送移动请求
                    SendMoveRequest(input);
                }
            }
        }

        private InputData GetInputChange()
        {
            var move = m_Input.move;
            InputData lastInput = null;
            if (m_InputHistory.Count > 0)
                lastInput = m_InputHistory[m_InputHistory.Count - 1];
            // 输入有变化，或者从静止变为移动
            if ((lastInput == null && move != Vector2.zero) || (lastInput != null && lastInput.Move != move))
            {
                var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds();
                var inputData = new InputData()
                {
                    Position = transform.position,
                    Move = move,
                    Timestamp = timestamp
                };
                return inputData;
            }
            return null;
        }

        private InputData PeekLastInput()
        {
            lock (m_InputHistory)
            {
                if (m_InputHistory.Count > 0)
                    return m_InputHistory[m_InputHistory.Count - 1];
                return null;
            }
        }

        private void DiscardInputsBefore(long timestamp)
        {
            lock (m_InputHistory)
            {
                // 将该时间点前的输入历史移除，但保留最后一条以便继续预测
                int index = m_InputHistory.FindLastIndex(input => input.Timestamp <= timestamp);
                if (index > 0)
                    m_InputHistory.RemoveRange(0, index);
            }
        }

        private void SendMoveRequest(InputData input)
        {
            PCharacterMove msg = new()
            {
                CharacterId = m_CharacterId,
                Move = new PVector3()
                {
                    X = input.Move.x,
                    Y = input.Move.y
                },
                Timestamp = input.Timestamp,
            };
            _ = NetworkManager.Instance.RequestAsync(PMsgId.CharacterMove, msg);
        }

        private Vector3 GetPositionAtTime(long timestamp)
        {
            // 从输入历史中该时间戳之前的输入
            InputData lastInput = null;
            int index = m_InputHistory.FindLastIndex(input => input.Timestamp <= timestamp);
            if (index != -1)
                lastInput = m_InputHistory[index];
            if (lastInput == null)
            {
                // 没有找到对应的输入，返回当前位置
                return transform.position;
            }
            return lastInput.Position + new Vector3(lastInput.Move.x, lastInput.Move.y, 0) * (int)(timestamp - lastInput.Timestamp) / 1000f * m_Speed;
        }

        private void VerifyInputHistory(PCharacterState state)
        {
            if (m_InputHistory.Count == 0)
            {
                // 没有输入历史，直接设置位置
                transform.position = new Vector3(state.Pos.X, state.Pos.Y, 0);
                return;
            }
            // NOTE 基于服务器消息是有序的
            // 丢弃之前的输入历史
            DiscardInputsBefore(state.Timestamp);

            var predictPos = GetPositionAtTime(state.Timestamp);
            var serverPos = state.Pos.ToVector3();
            // 预测位置与服务器位置接近，无需调整
            if ((predictPos - serverPos).sqrMagnitude <= 0.01f)
                return;

            Log.Information("[Move] Position corrected from {PredictPos} to {ServerPos} at {Time}", predictPos, serverPos, state.Timestamp);
            // 预测位置与服务器位置差距较大
            if (m_InputHistory.Count == 0)
            {
                transform.position = serverPos;
            }
            else
            {
                // 更新历史输入的预测位置
                m_InputHistory[0].Position = serverPos;
                m_InputHistory[0].Timestamp = state.Timestamp;
                for (int i = 1; i < m_InputHistory.Count; i++)
                {
                    var input = m_InputHistory[i];
                    var prevInput = m_InputHistory[i - 1];
                    float deltaTime = input.Timestamp - prevInput.Timestamp;
                    input.Position = prevInput.Position + new Vector3(input.Move.x, input.Move.y, 0) * deltaTime * m_Speed;
                }
                var lastInput = PeekLastInput();
                // 基于最新的输入预测当前位置
                if (lastInput != null)
                {
                    var now = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds();
                    transform.position = GetPositionAtTime(now);
                }
            }
        }

        public void OnReceiveServerUpdate(PCharacterState state)
        {
            if (m_Input != null)
            {
                VerifyInputHistory(state);
            }
            else
            {
                // 非本地角色，直接设置位置
                transform.position = new Vector3(state.Pos.X, state.Pos.Y, 0);
            }
        }

        private class InputData
        {
            public Vector3 Position;
            public Vector2 Move;
            public long Timestamp;
        }
    }
}