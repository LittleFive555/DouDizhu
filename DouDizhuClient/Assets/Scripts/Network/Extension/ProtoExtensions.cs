using UnityEngine;

namespace Network.Proto
{
    internal static class ProtoExtensions
    {
        public static Vector3 ToVector3(this PVector3 pVector3)
        {
            return new Vector3(pVector3.X, pVector3.Y, pVector3.Z);
        }

        public static PVector3 ToPVector3(this Vector3 vector3)
        {
            return new PVector3() { X = vector3.x, Y = vector3.y, Z = vector3.z };
        }
    }
}
