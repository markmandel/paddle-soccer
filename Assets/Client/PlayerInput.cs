using UnityEngine;

// Moves the Player (paddle and camera) around
// when using the keyboard
public class PlayerInput : MonoBehaviour
{
    // speed of movement, in all directions
    [SerializeField] private float speed = 0.1f;

    // sensitivity of the mouse on the Y axis
    [SerializeField] private float mouseSensitivity = 5.0f;

    // Lock down the cursor
    void Start()
    {
        // locking inside the editor seems to ruin input
        if (!Application.isEditor)
        {
            Cursor.lockState = CursorLockMode.Locked;
            Cursor.visible = false;
        }
    }

    // Handle rotational mouse input
    void Update()
    {
        float axis = Input.GetAxis("Mouse X");
        transform.Rotate(0, (axis * mouseSensitivity), 0);

        // Disable Mouse cursor locking
        if (Input.GetKeyDown(KeyCode.Escape))
        {
            Cursor.lockState = CursorLockMode.None;
            Cursor.visible = true;
        }
    }

    // Handle forward, left and right
    void FixedUpdate()
    {
        float deltaX = Input.GetAxis("Horizontal") * speed;
        float deltaZ = Input.GetAxis("Vertical") * speed;

        transform.Translate(deltaX, 0, deltaZ);
    }
}