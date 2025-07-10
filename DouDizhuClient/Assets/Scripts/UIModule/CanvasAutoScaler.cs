using UnityEngine;
using UnityEngine.UI;

namespace UIModule
{
    [RequireComponent(typeof(CanvasScaler))]
    public class CanvasAutoScaler : MonoBehaviour
    {
        private CanvasScaler m_CanvasScaler;

        private const float ASPECT_RATIO_MATCH_WIDTH = 4f / 3f;
        private const float ASPECT_RATIO_MATCH_HEIGHT = 16f / 9f;

        private float m_LastAspectRatio = 0;

        private void Awake()
        {
            m_CanvasScaler = GetComponent<CanvasScaler>();
        }
        
        private void LateUpdate()
        {
            float aspectRatio = (float)Screen.width / Screen.height;
            if (m_LastAspectRatio == aspectRatio)
                return;
                
            m_LastAspectRatio = aspectRatio;
            if (aspectRatio > ASPECT_RATIO_MATCH_HEIGHT)
            {
                m_CanvasScaler.matchWidthOrHeight = 1;
            }
            else if (aspectRatio < ASPECT_RATIO_MATCH_WIDTH)
            {
                m_CanvasScaler.matchWidthOrHeight = 0;
            }
            else
            {
                float range = ASPECT_RATIO_MATCH_HEIGHT - ASPECT_RATIO_MATCH_WIDTH;
                float offset = aspectRatio - ASPECT_RATIO_MATCH_WIDTH;
                m_CanvasScaler.matchWidthOrHeight = offset / range;
            }
        }
    }
}