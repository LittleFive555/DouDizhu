using System.Collections.Generic;
using UnityEngine;

namespace UIModule
{
    public class UIRoot : MonoBehaviour
    {
        private Dictionary<EnumUILayer, UILayer> m_UILayers = new Dictionary<EnumUILayer, UILayer>();

        private void Awake()
        {
            foreach (var uiLayer in GetComponentsInChildren<UILayer>())
                m_UILayers.Add(uiLayer.Layer, uiLayer);
        }

        public void AppendToLayer(EnumUILayer layer, GameObject gameObject)
        {
            gameObject.transform.SetParent(m_UILayers[layer].transform);
        }
    }
}
