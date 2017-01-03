using UnityEngine;

namespace Client.Player
{
    // Moves the paddle around!
    [RequireComponent(typeof(Rigidbody))]
    public class PaddleInput : MonoBehaviour
    {
        [SerializeField]
        [Tooltip("Speed of movement")]
        private float speed = 10f;

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
        }

        // This is to deal with some issues with the
        // Mesh renderer. - if you hit it, bounce it back!
        void OnCollisionEnter(Collision collision)
        {
            if(collision.gameObject.CompareTag("StopPlayer"))
            {
                Debug.Log("Bouncing!!!");
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
                Vector3 targetVelocity = new Vector3(deltaX, 0, deltaY);
                Debug.Log(string.Format("[1] Local Input Velocity: {0}", targetVelocity));

                // convert from local to world
                targetVelocity = transform.TransformDirection(targetVelocity);
                Debug.Log(string.Format("[2] World Input Velocity: {0}", targetVelocity));

                targetVelocity *= speed;

                Debug.Log(string.Format("[3] World Speed Input Velocity: {0}", targetVelocity));

                // Apply a force, to reach the target velocity
                Vector3 currentVelocity = rb.velocity;
                Vector3 delta = targetVelocity - currentVelocity;
                rb.AddForce(delta, ForceMode.VelocityChange);
            }
        }
    }
}