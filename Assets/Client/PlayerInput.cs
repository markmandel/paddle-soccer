using UnityEngine;

// Moves the Player (paddle and camera) around
// when using the keyboard
public class PlayerInput : MonoBehaviour
{
    // speed of movement, in all directions
    [SerializeField] private float speed = 0.1f;

    void FixedUpdate()
    {
        float deltaX = Input.GetAxis("Horizontal") * speed;
        float deltaZ = Input.GetAxis("Vertical") * speed;

        transform.Translate(deltaX, 0, deltaZ);
    }
}