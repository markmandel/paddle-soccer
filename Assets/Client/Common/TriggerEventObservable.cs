using UnityEngine;

namespace Client.Common
{
    // Observable event for when collision triggers occur
    // on isTrigger Collisionable Game Objects
    public class TriggerEventObservable : MonoBehaviour
    {
        public delegate void Triggered(Collider other);

        public event Triggered TriggerEnter;
        public event Triggered TriggerExit;
        public event Triggered TriggerStay;

        // Fires TriggerEnter event
        void OnTriggerEnter(Collider other)
        {
            if(TriggerEnter != null)
            {
                TriggerEnter(other);
            }
        }

        // Fires TriggerExit event
        void OnTriggerExit(Collider other)
        {
            if(TriggerExit != null)
            {
                TriggerExit(other);
            }
        }

        // Fires TriggerStay event
        void OnTriggerStay(Collider other)
        {
            if(TriggerStay != null)
            {
                TriggerStay(other);
            }
        }
    }
}