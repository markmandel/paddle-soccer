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
        private float kickForce = 20f;

        [SerializeField]
        [Tooltip("Distance the paddle can kick from")]
        private float kickDistance = 1.5f;

        [SerializeField]
        [Tooltip("How much to rotate when using the keyboard input")]
        private float keyboardRotation = 0.5f;

        [SerializeField]
        [Tooltip("How far down to the bottom to kick. 2f is the bottom.")]
        private float kickAngle = 2.7f;

        private Rigidbody rb;
        private BoxCollider box;

        // --- Messages ---

        private void Start()
        {
            rb = GetComponent<Rigidbody>();
            box = GetComponent<BoxCollider>();
        }

        private void Update()
        {
            if(Input.GetKeyDown(KeyCode.Space) || Input.GetButtonDown("Fire1"))
            {
                KickBall();
            }
        }

        // Handle forward, left and right
        private void FixedUpdate()
        {
            KeyboardHorizontalInput();
            PlayerRotation(Input.GetAxis("Mouse X"));
            KeyboardRotation();
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

        private void PlayerRotation(float axis)
        {
            axis = axis * rotationalSpeed;
            Quaternion rotation = Quaternion.Euler(0, transform.localEulerAngles.y + axis, 0);
            rb.MoveRotation(rotation);
        }

        private void KeyboardRotation()
        {
            if(Input.GetKey(KeyCode.LeftBracket))
            {
                PlayerRotation(-keyboardRotation);
            }
            if(Input.GetKey(KeyCode.RightBracket))
            {
                PlayerRotation(keyboardRotation);
            }
        }

        private void KickBall()
        {
            Vector3 diff = new Vector3(0, box.size.y / kickAngle, 0);
            Vector3 origin = transform.position - transform.TransformVector(diff);

            RaycastHit hit;
            if(Physics.Raycast(origin, transform.forward, out hit, kickDistance))
            {
                if(hit.collider.name == "Ball")
                {
                    Rigidbody crb = hit.collider.GetComponent<Rigidbody>();
                    Vector3 force = -kickForce * hit.normal;
                    crb.AddForceAtPosition(force, hit.point, ForceMode.Impulse);
                }
            }
        }
    }
}