using System.Net.Sockets;

namespace Network.Tcp
{
    public interface MessageReadWriter
    {
        byte[] ReadFrom(NetworkStream networkStream);
        void WriteTo(NetworkStream networkStream, byte[] messageBytes);
    }
}