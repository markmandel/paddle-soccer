using Microsoft.Win32;
using UnityEngine;

namespace Client.Player
{
    // Moves the paddle around!
    [RequireComponent(typeof(Rigidbody))]
    [RequireComponent(typeof(BoxCollider))]
    public class PaddleInput : MonoBehaviour
    {
        [SerializeField]
        [Tooltip("Speed of horizontal movement")]
        private float horizontalSpeed = 10f;

        [SerializeField]
        [Tooltip("Speed of rotation by mouse")]
        private float rotationalSpeed = 5f;

        [SerializeField]
        [Tooltip("How hard to kick the ball")]
        private float kickForce = 3f;

        [SerializeField]
        [Tooltip("Distance the paddle can kick from")]
        private float kickDistance = 0.7f;

        [SerializeField]
        [Tooltip("Clamp for kick magnitude")]
        private float kickMagnitude = 50f;

        private Rigidbody rb;
        private BoxCollider box;

        // --- Messages ---

        private void Start()
        {
            rb = GetComponent<Rigidbody>();
            box = GetComponent<BoxCollider>();

            if(!Application.isEditor)
            {
                Cursor.lockState = CursorLockMode.Locked;
                Cursor.visible = false;
            }
        }

        private void Update()
        {
            if(Input.GetKeyDown(KeyCode.Space) || Input.GetButtonDown("Fire1"))
            {
                Debug.Log("Attempting Kick!");
                Vector3 bottom = new Vector3(transform.position.x, transform.position.y - (box.size.y / 1.5f),
                    transform.position.z);

                RaycastHit hit;
                if(Physics.SphereCast(bottom, (box.size.x / 2), transform.forward, out hit, kickDistance))
                {
                    Debug.Log("Hit something?");
                    if(hit.collider.name == "Ball")
                    {
                        Debug.Log("KICK BALL!");
                        Rigidbody crb = hit.collider.GetComponent<Rigidbody>();
                        Vector3 force = (1f / hit.distance) * -kickForce * hit.normal;
                        force = Vector3.ClampMagnitude(force, kickMagnitude);
                        Debug.Log(string.Format("Force: {0}", force));
                        crb.AddForceAtPosition(force, hit.point, ForceMode.Impulse);
                    }
                }
            }
        }

        // Handle forward, left and right
        private void FixedUpdate()
        {
            KeyboardHorizontalInput();
            MouseRotation();
        }

        // This is to deal with some issues with the
        // Mesh renderer. - if you hit it, bounce it back!
        private void OnCollisionEnter(Collision collision)
        {
            if(collision.gameObject.CompareTag("StopPlayer"))
            {
                rb.AddForce(-1 * rb.mass * rb.velocity, ForceMode.Impulse);
            }
        }

        // --- Functions ---

        // Basically work out what speed and direction we *want* to be
        // going in, and provide a force that will do such a thing
        // Credit and inspiration from: http://wiki.unity3d.com/index.php?title=RigidbodyFPSWalker
        private void KeyboardHorizontalInput()
        {
            float deltaX = Input.GetAxis("Horizontal");
            float deltaY = Input.GetAxis("Vertical");

            // skip this whole thing, if there is no input
            if(!(deltaX == 0 && deltaY == 0))
            {
                Vector3 targetVelocity = new Vector3(deltaX, 0, deltaY) * horizontalSpeed;
                // convert from local to world
                targetVelocity = transform.TransformDirection(targetVelocity);

                // Apply a force, to reach the target velocity
                Vector3 currentVelocity = rb.velocity;
                Vector3 delta = targetVelocity - currentVelocity;
                rb.AddForce(delta, ForceMode.VelocityChange);
            }
        }

        private void MouseRotation()
        {
            float axis = Input.GetAxis("Mouse X") * rotationalSpeed;
            Quaternion rotation = Quaternion.Euler(0, transform.localEulerAngles.y + axis, 0);
            rb.MoveRotation(rotation);
        }
    }
}