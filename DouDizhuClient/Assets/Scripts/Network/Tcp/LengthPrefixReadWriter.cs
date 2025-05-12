using System;
using System.Net.Sockets;
using System.Threading.Tasks;

namespace Network.Tcp
{
    public class LengthPrefixReadWriter : MessageReadWriter
    {
        public async Task<byte[]> ReadFrom(NetworkStream networkStream)
        {
            // 读取4字节的长度信息
            byte[] lengthBytes = new byte[4];
            await networkStream.ReadAsync(lengthBytes, 0, 4);
            
            // 如果是小端序，先反转字节序
            if (BitConverter.IsLittleEndian)
                Array.Reverse(lengthBytes);
                
            // 转换为整数
            int length = BitConverter.ToInt32(lengthBytes, 0);
            
            // 读取消息内容
            byte[] messageBytes = new byte[length];
            await networkStream.ReadAsync(messageBytes, 0, length);
            return messageBytes;
        }

        public async Task WriteTo(NetworkStream networkStream, byte[] messageBytes)
        {
            byte[] lengthBytes = BitConverter.GetBytes(messageBytes.Length);
            if (BitConverter.IsLittleEndian) // 如果系统是小端序，则需要反转字节顺序
                Array.Reverse(lengthBytes);
            await networkStream.WriteAsync(lengthBytes, 0, lengthBytes.Length);
            await networkStream.WriteAsync(messageBytes, 0, messageBytes.Length);
        }
    }
}
