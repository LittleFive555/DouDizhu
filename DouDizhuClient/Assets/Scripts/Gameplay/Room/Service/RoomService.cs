using System.Collections.Generic;
using System.Threading.Tasks;
using Network;
using Network.Proto;

namespace Gameplay.Room.Service
{
    public class RoomService
    {
        public static async Task<IReadOnlyList<PRoom>> GetRoomList()
        {
            var response = await NetworkManager.Instance.RequestAsync<PGetRoomListRequest, PGetRoomListResponse>(PMsgId.GetRoomList, new PGetRoomListRequest());
            if (response.IsSuccess)
            {
                return response.Data.Rooms;
            }
            return null;
        }

        public static async Task<PRoom> CreateRoom(string roomName)
        {
            var response = await NetworkManager.Instance.RequestAsync<PCreateRoomRequest, PCreateRoomResponse>(PMsgId.CreateRoom, new PCreateRoomRequest()
            {
                RoomName = roomName
            });
            if (response.IsSuccess)
            {
                return response.Data.Room;
            }
            return null;
        }

        public static async Task<PRoom> EnterRoom(uint roomId)
        {
            var response = await NetworkManager.Instance.RequestAsync<PEnterRoomRequest, PEnterRoomResponse>(PMsgId.EnterRoom, new PEnterRoomRequest()
            {
                RoomId = roomId
            });
            if (response.IsSuccess)
            {
                return response.Data.Room;
            }
            return null;
        }

        public static async Task<bool> LeaveRoom()
        {
            var response = await NetworkManager.Instance.RequestAsync(PMsgId.LeaveRoom, new PLeaveRoomRequest());
            if (response.IsSuccess)
            {
                return true;
            }
            return false;
        }
    }
}