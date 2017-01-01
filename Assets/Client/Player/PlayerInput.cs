using UnityEngine;

// Moves the Player (paddle and camera) around
// when using the keyboard
public class PlayerInput : MonoBehaviour
{
    // speed of movement, in all directions
    [SerializeField] private float speed = 0.1f;

    // sensitivity of the mouse on the Y axis
    [SerializeField] private float mouseSensitivity = 5.0f;

    // --- Handlers ---

    // Lock down the cursor
    void Start()
    {
        // locking inside the editor seems to ruin input
        if(!Application.isEditor)
        {
            Cursor.lockState = CursorLockMode.Locked;
            Cursor.visible = false;
        }
    }

    // Handle rotational mouse input
    void Update()
    {
        float yAngle = Input.GetAxis("Mouse X") * mouseSensitivity;
        RotatePlayer(yAngle);

        // Disable Mouse cursor locking
        if(Input.GetKeyDown(KeyCode.Escape))
        {
            Cursor.lockState = CursorLockMode.None;
            Cursor.visible = true;
        }

        // Test out rotation
        if(Input.GetKey(KeyCode.LeftBracket))
        {
            RotatePlayer(transform.rotation.y - 5f);
        }
        else if(Input.GetKey(KeyCode.RightBracket))
        {
            RotatePlayer(transform.rotation.y + 5f);
        }
    }

    // Handle forward, left and right
    void FixedUpdate()
    {
        float deltaX = Input.GetAxis("Horizontal") * speed;
        float deltaZ = Input.GetAxis("Vertical") * speed;

        transform.Translate(deltaX, 0, deltaZ);
    }

    // --- Custom Functions ---

    private void RotatePlayer(float yAngle)
    {
        Debug.Log("Rotate Player: " + yAngle);
        transform.Rotate(0, yAngle, 0);
    }
}