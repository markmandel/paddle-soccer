using UnityEngine;

namespace Client.Player
{
    // Moves the paddle around!
    [RequireComponent(typeof(Rigidbody))]
    public class PaddleInput : MonoBehaviour
    {
        [SerializeField]
        [Tooltip("Speed of horizontal movement")]
        private float horizontalSpeed = 10f;

        [SerializeField]
        [Tooltip("Speed of rotation by mouse")]
        private float rotationalSpeed = 1f;

        private Rigidbody rb;

        // --- Messages ---

        void Start()
        {
            rb = GetComponent<Rigidbody>();
        }

        // Handle forward, left and right
        void FixedUpdate()
        {
            KeyboardHorizontalInput();

            float axis = Input.GetAxis("Mouse X");
            Quaternion rotation = Quaternion.Euler(0, transform.localEulerAngles.y + axis, 0);
            rb.MoveRotation(rotation);
        }

        // This is to deal with some issues with the
        // Mesh renderer. - if you hit it, bounce it back!
        void OnCollisionEnter(Collision collision)
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
    }
}