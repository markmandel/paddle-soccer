using System;
using Assets.Client.Common;
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

        private TriggerEventObservable triggerEventObservable;

        // if this true, then player isn't going to move around
        private Boolean stop = false;

        // --- Handlers ---


        void Start()
        {
            triggerEventObservable = GetComponentInChildren<TriggerEventObservable>();
            if(triggerEventObservable == null)
            {
                throw new Exception("One of the player input children should have a TriggerEventObservable");
            }

            triggerEventObservable.TriggerEnter += TriggerEnter;
            triggerEventObservable.TriggerExit += TriggerExit;

            LockCursor();
        }

        // remove listeners
        void OnDestroy()
        {
            triggerEventObservable.TriggerEnter -= TriggerEnter;
            triggerEventObservable.TriggerExit -= TriggerExit;
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

        // Lock down the cursor (if not in editor)
        private static void LockCursor()
        {
            // locking inside the editor seems to ruin input
            if(!Application.isEditor)
            {
                Cursor.lockState = CursorLockMode.Locked;
                Cursor.visible = false;
            }
        }

        // Keyboard input, for forward, back, left and right
        private void KeyboardHorizontalInput()
        {
            if(!stop)
            {
                float deltaX = Input.GetAxis("Horizontal") * speed;
                float deltaZ = Input.GetAxis("Vertical") * speed;

                transform.Translate(deltaX, 0, deltaZ);
            }
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

        private void TriggerEnter(Collider other)
        {
            if(other.tag == "StopPlayer")
            {
                Debug.Log(String.Format("Player input has entered: {0} / {1}", other.name, other.tag));
                stop = true;
            }
        }

        private void TriggerExit(Collider other)
        {
            if(other.tag == "StopPlayer")
            {
                Debug.Log(String.Format("Player input has exit: {0} / {1}", other.name, other.tag));
                stop = false;
            }
        }
    }
}