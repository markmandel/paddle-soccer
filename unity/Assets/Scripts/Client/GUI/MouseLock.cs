using UnityEngine;

namespace Client.GUI
{
    /// <summary>
    /// Manages the mouse locked into the game
    /// </summary>
    public class MouseLock : MonoBehaviour
    {
        // --- Messages ---

        /// <summary>
        /// Lock and set the mouse to not visible
        /// </summary>
        private void Start()
        {
            LockMouse();
        }

        /// <summary>
        /// Escaping and capture of the mouse
        /// </summary>
        private void Update()
        {
            if(Input.GetButtonDown("Fire1") && Cursor.visible)
            {
                LockMouse();
            }
            else if(Input.GetKeyDown(KeyCode.Escape))
            {
                Cursor.lockState = CursorLockMode.None;
                Cursor.visible = true;
            }
        }

        // --- Functions ---

        /// <summary>
        /// Lock the mouse
        /// </summary>
        private static void LockMouse()
        {
            if(!Application.isEditor)
            {
                Cursor.lockState = CursorLockMode.Locked;
                Cursor.visible = false;
            }
        }
    }
}