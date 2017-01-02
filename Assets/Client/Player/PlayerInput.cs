using UnityEngine;

namespace Assets.Client.Player
{
    // Moves the Player (paddle and camera) around
    // when using the keyboard
    public class PlayerInput : MonoBehaviour
    {
        // speed of movement, in all directions
        [SerializeField] private float speed = 0.1f;

        // sensitivity of the mouse on the Y axis
        [SerializeField] private float mouseSensitivity = 5.0f;

        // --- Handlers ---

        // Lock down the cursor (if not in editor)
        void Start()
        {
            // locking inside the editor seems to ruin input
            if(!Application.isEditor)
            {
                Cursor.lockState = CursorLockMode.Locked;
                Cursor.visible = false;
            }
        }

        // Handle rotational input
        void Update()
        {
            MouseRotationInput();
            KeyboardRotationInput();
        }

        // Handle forward, left and right
        void FixedUpdate()
        {
            KeyboardHorizontalInput();
        }

        // --- Functions ---

        // Keyboard input, for forward, back, left and right
        private void KeyboardHorizontalInput()
        {
            float deltaX = Input.GetAxis("Horizontal") * speed;
            float deltaZ = Input.GetAxis("Vertical") * speed;

            transform.Translate(deltaX, 0, deltaZ);
        }

        // Handles rotational movement with the mouse
        private void MouseRotationInput()
        {
            float yAngle = Input.GetAxis("Mouse X") * mouseSensitivity;
            RotatePlayer(yAngle);

            // Disable Mouse cursor locking
            if(Input.GetKeyDown(KeyCode.Escape))
            {
                Cursor.lockState = CursorLockMode.None;
                Cursor.visible = true;
            }
        }

        // Handle Keyboard Input for rotation: [ and ]
        private void KeyboardRotationInput()
        {
            if(Input.GetKey(KeyCode.LeftBracket))
            {
                RotatePlayer(transform.rotation.y - 5f);
            }
            else if(Input.GetKey(KeyCode.RightBracket))
            {
                RotatePlayer(transform.rotation.y + 5f);
            }
        }

        private void RotatePlayer(float yAngle)
        {
            transform.Rotate(0, yAngle, 0);
        }
    }
}