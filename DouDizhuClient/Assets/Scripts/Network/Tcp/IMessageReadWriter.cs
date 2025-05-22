using System.Net.Sockets;
using System.Threading.Tasks;

namespace Network.Tcp
{
    public interface IMessageReadWriter
    {
        Task<byte[]> ReadFrom(NetworkStream networkStream);
        Task WriteTo(NetworkStream networkStream, byte[] messageBytes);
    }
}