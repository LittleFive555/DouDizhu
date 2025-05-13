using UnityEngine;

namespace UIModule
{
    public class UILayer : MonoBehaviour
    {
        [SerializeField]
        private EnumUILayer m_Layer;
        public EnumUILayer Layer => m_Layer;
    }
}
