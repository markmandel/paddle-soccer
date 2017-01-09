using UnityEngine;

namespace Client.GUI
{
    // Manages the mouse locked into the game
    public class MouseLock : MonoBehaviour
    {
        // --- Messages ---

        // Lock and set the mouse to not visible
        private void Start()
        {
            LockMouse();
        }

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