using Gameplay.Character;
using Network;
using Network.Proto;
using System.Collections.Generic;
using UnityEngine;

namespace Gameplay.World
{
    public class RoomWorld : MonoBehaviour
    {
        [SerializeField]
        private NetworkMoveTest m_MoveTestProto;

        private Dictionary<string, NetworkMoveTest> m_CharacterMap = new();

        private NetworkMoveTest m_Player;

        public async void Initialize()
        {
            NetworkManager.Instance.RegisterNotificationHandler<PWorldState>(PMsgId.WorldState, OnWorldChange);

            var response = await NetworkManager.Instance.RequestAsync<PEnterWorldRequest, PEnterWorldResponse>(PMsgId.EnterWorld, new PEnterWorldRequest() { WorldId = "playground-111999555" });
            if (response.IsSuccess)
            {
                foreach (var character in response.Data.WorldState.Characters)
                {
                    var characterObj = Instantiate(m_MoveTestProto, transform);
                    var characterId = character.Key;
                    characterObj.name = characterId;
                    characterObj.OnReceiveServerUpdate(character.Value.States[0]);
                    m_CharacterMap[characterId] = characterObj;

                    if (characterId == response.Data.CharacterId)
                    {
                        m_Player = characterObj;
                        characterObj.Initialize(characterId, FindObjectOfType<StarterAssetsInputs>());
                    }
                    else
                    {
                        characterObj.Initialize(characterId);
                    }
                }
            }
        }

        private void OnWorldChange(PWorldState worldState)
        {
            Debug.Log("World state changed: " + worldState);
            // Handle the world state change here
            foreach (var character in worldState.Characters)
            {
                if (m_CharacterMap.TryGetValue(character.Key, out var characterObj))
                {
                    foreach (var state in character.Value.States)
                        characterObj.OnReceiveServerUpdate(state);
                }
            }
        }
    }
}
