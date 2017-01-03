using System;
using UnityEngine;

namespace Assets.Client.Player
{
    // Moves the paddle around!
    [RequireComponent(typeof(Rigidbody))]
    public class PaddleInput : MonoBehaviour
    {
        // speed of movement, in all directions
        [SerializeField] private float speed = 5f;

        private Rigidbody rb;

        // --- Messages ---

        void Start()
        {
            Debug.Log("Starting Paddle Input!");
            rb = GetComponent<Rigidbody>();
        }

        // Handle forward, left and right
        void FixedUpdate()
        {
            KeyboardHorizontalInput();
        }

        // --- Functions ---

        private void KeyboardHorizontalInput()
        {
            float deltaX = Input.GetAxis("Horizontal") * speed;
            float deltaZ = Input.GetAxis("Vertical") * speed;

            if(Math.Abs(deltaZ + deltaX) > 0)
            {
                Vector3 translate = new Vector3(deltaX, 0, deltaZ);
                rb.MovePosition(rb.position + translate * Time.deltaTime);
            }
        }
    }
}