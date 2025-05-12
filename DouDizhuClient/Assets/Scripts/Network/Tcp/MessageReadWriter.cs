using System.Net.Sockets;
using System.Threading.Tasks;

namespace Network.Tcp
{
    public interface MessageReadWriter
    {
        Task<byte[]> ReadFrom(NetworkStream networkStream);
        Task WriteTo(NetworkStream networkStream, byte[] messageBytes);
    }
}