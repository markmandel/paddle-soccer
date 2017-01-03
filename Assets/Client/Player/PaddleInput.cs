using System;
using System.Security.Cryptography;
using UnityEngine;

namespace Assets.Client.Player
{
    // Moves the paddle around!
    [RequireComponent(typeof(Rigidbody))]
    public class PaddleInput : MonoBehaviour
    {
        [SerializeField]
        [Tooltip("Force to be applied at movement")]
        private float force = 50f;

        [SerializeField]
        [Tooltip("Drag to slow down to when not using input")]
        private float slowDrag = 5f;

        [SerializeField]
        [Tooltip("Maximum velocity")]
        private float maxSpeed = 10f;

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

        // This is to deal with some issues with the
        // Mesh renderer. - if you hit it, bounce it back!
        void OnCollisionEnter(Collision collision)
        {
            if(collision.gameObject.tag == "StopPlayer")
            {
                Debug.Log("Bouncing!!!");
                rb.AddForce(-1 * rb.mass * rb.velocity, ForceMode.Impulse);
            }
        }

        // --- Functions ---

        private void KeyboardHorizontalInput()
        {
            float deltaX = Input.GetAxis("Horizontal") * force;
            float deltaZ = Input.GetAxis("Vertical") * force;

            if(deltaX != 0f || deltaZ != 0f)
            {
                rb.drag = 0;
                if(rb.velocity.magnitude < maxSpeed)
                {
                    Vector3 translate = new Vector3(deltaX, 0, deltaZ);
                    rb.AddRelativeForce(translate, ForceMode.Force);
                }
            }
            else
            {
                rb.drag = slowDrag;
                //rb.velocity = Vector3.zero;
            }
        }
    }
}