using UnityEngine;

namespace UIModule
{
    public abstract class UIComponentBase : MonoBehaviour
    {
        public struct EmptyArgs { }

        public virtual void OnShowBegin<TArgs>(TArgs? args) where TArgs : struct
        {

        }

        public virtual void OnShowFinish<TArgs>(TArgs? args) where TArgs : struct
        {

        }

        public virtual void OnHideBegin()
        {

        }

        public virtual void OnHideFinish()
        {

        }
    }
}
