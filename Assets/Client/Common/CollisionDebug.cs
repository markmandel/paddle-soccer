using UnityEngine;

// Simply writes out to console when a collision
// occurs, and tells you with what.
public class CollisionDebug : MonoBehaviour
{
    // --- Handlers ---

    void OnCollisionEnter(Collision collision)
    {
        Log(collision, "Enter");
    }

    void OnCollisionExit(Collision collision)
    {
        Log(collision, "Exit");
    }

    void OnCollisionStay(Collision collision)
    {
        Log(collision, "Stay");
    }

    // --- Functions ---

    private void Log(Collision collision, string category)
    {
        Debug.Log(string.Format("[{0}] Collision {1}: {2}", name, category, collision));
    }
}