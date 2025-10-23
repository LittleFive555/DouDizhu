using System.Collections.Generic;
using System.Threading.Tasks;
using Gameplay.Room.Service;
using Gameplay.Room.View;
using Network.Proto;
using UIModule;

namespace Gameplay.Room.Controller
{
    public class RoomController
    {
        public static async Task CreateRoom(string roomName)
        {
            var room = await RoomService.CreateRoom(roomName);
            if (room != null)
                UIManager.Instance.ShowUI<UIRoom, UIRoom.Args>(new UIRoom.Args() { Room = room });
        }

        public static async Task EnterRoom(uint roomId)
        {
            var result = await RoomService.EnterRoom(roomId);
            if (result != null)
                UIManager.Instance.ShowUI<UIRoom, UIRoom.Args>(new UIRoom.Args() { Room = result });
        }

        public static async Task<IReadOnlyList<PRoom>> RefreshRoomList()
        {
            return await RoomService.GetRoomList();
        }

        public static async Task LeaveRoom()
        {
            var result = await RoomService.LeaveRoom();
            if (result)
                UIManager.Instance.HideUI<UIRoom>();
        }
    }
}