// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using Server;
using UnityEngine;
using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Moves the paddle around!
    /// </summary>
    [RequireComponent(typeof(Rigidbody))]
    [RequireComponent(typeof(BoxCollider))]
    [RequireComponent(typeof(PlayerAction))]
    public class PaddleInput : NetworkBehaviour
    {
        [SerializeField]
        [Tooltip("Speed of horizontal movement")]
        private float horizontalSpeed = 10f;

        [SerializeField]
        [Tooltip("Speed of rotation by mouse")]
        private float rotationalSpeed = 5f;

        [SerializeField]
        [Tooltip("How much to rotate when using the keyboard input")]
        private float keyboardRotation = 0.5f;

        private Rigidbody rb;
        private PlayerAction pa;

        // --- Messages ---

        /// <summary>
        /// Gab teh rigidbody and box collider
        /// </summary>
        private void Start()
        {
            rb = GetComponent<Rigidbody>();
            pa = GetComponent<PlayerAction>();
        }

        /// <summary>
        /// Handles kicks, if local player
        /// </summary>
        private void Update()
        {
            if (isLocalPlayer)
            {
                if (Input.GetKeyDown(KeyCode.Space) || Input.GetButtonDown("Fire1"))
                {
                    pa.CmdKickBall();
                }
            }
        }

        /// <summary>
        /// Handels moving the rigidbody around via input
        /// </summary>
        private void FixedUpdate()
        {
            if (isLocalPlayer)
            {
                KeyboardHorizontalInput();
                PlayerRotation(Input.GetAxis("Mouse X"));
                KeyboardRotation();
            }
        }

        // --- Functions ---

        /// <summary>
        /// Basically work out what speed and direction we *want* to be
        /// going in, and provide a force that will do such a thing
        /// Credit and inspiration from: http://wiki.unity3d.com/index.php?title=RigidbodyFPSWalker
        /// </summary>
        private void KeyboardHorizontalInput()
        {
            var deltaX = Input.GetAxis("Horizontal");
            var deltaY = Input.GetAxis("Vertical");

            // skip this whole thing, if there is no input
            if (!(deltaX == 0 && deltaY == 0))
            {
                var targetVelocity = new Vector3(deltaX, 0, deltaY) * horizontalSpeed;
                // convert from local to world
                targetVelocity = transform.TransformDirection(targetVelocity);

                // Apply a force, to reach the target velocity
                var currentVelocity = rb.velocity;
                var delta = targetVelocity - currentVelocity;
                rb.AddForce(delta, ForceMode.VelocityChange);
            }
        }

        /// <summary>
        /// Manages the rotation of the player
        /// </summary>
        /// <param name="axis"></param>
        private void PlayerRotation(float axis)
        {
            axis = axis * rotationalSpeed;
            var rotation = Quaternion.Euler(0, transform.localEulerAngles.y + axis, 0);
            rb.MoveRotation(rotation);
        }

        /// <summary>
        /// Optional keyboard rotation. Handy when using a trackpad.
        /// </summary>
        private void KeyboardRotation()
        {
            if (Input.GetKey(KeyCode.LeftBracket))
            {
                PlayerRotation(-keyboardRotation);
            }
            if (Input.GetKey(KeyCode.RightBracket))
            {
                PlayerRotation(keyboardRotation);
            }
        }
    }
}