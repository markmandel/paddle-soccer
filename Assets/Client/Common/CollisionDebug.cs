using UnityEngine;

namespace Assets.Client.Common
{
    // Simply writes out to console when a collision
    // occurs, and tells you with what.
    public class CollisionDebug : MonoBehaviour
    {
        // --- Handlers ---

        void OnCollisionEnter(Collision collision)
        {
            Log(collision.gameObject, "Collision Enter");
        }

        void OnCollisionExit(Collision collision)
        {
            Log(collision.gameObject, "Collision Exit");
        }

        void OnCollisionStay(Collision collision)
        {
            Log(collision.gameObject, "Collision Stay");
        }

        void OnTriggerEnter(Collider other)
        {
            Log(other.gameObject, "Trigger Enter");
        }

        void OnTriggerExit(Collider other)
        {
            Log(other.gameObject, "Trigger Exit");
        }

        void OnTriggerStay(Collider other)
        {
            Log(other.gameObject, "Trigger Stay");
        }

        // --- Functions ---

        private void Log(GameObject obj, string category)
        {
            Debug.Log(string.Format("[{0}] {1} Event: {2}", name, category, obj.name));
        }
    }
}