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

﻿using Game;
using UnityEngine;
using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Moves the paddle around!
    /// </summary>
    [RequireComponent(typeof(Rigidbody))]
    [RequireComponent(typeof(BoxCollider))]
    public class PaddleInput : NetworkBehaviour
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

        /// <summary>
        /// Gab teh rigidbody and box collider
        /// </summary>
        private void Start()
        {
            rb = GetComponent<Rigidbody>();
            box = GetComponent<BoxCollider>();
        }

        /// <summary>
        /// Handles kicks, if local player
        /// </summary>
        private void Update()
        {
            if (isLocalPlayer)
            {
                if(Input.GetKeyDown(KeyCode.Space) || Input.GetButtonDown("Fire1"))
                {
                    KickBall();
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
            if(!(deltaX == 0 && deltaY == 0))
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
            if(Input.GetKey(KeyCode.LeftBracket))
            {
                PlayerRotation(-keyboardRotation);
            }
            if(Input.GetKey(KeyCode.RightBracket))
            {
                PlayerRotation(keyboardRotation);
            }
        }

        /// <summary>
        /// Kicks the ball. (Currently not working)
        /// TODO: Needs to be moved to a RPC as the ball is server side.
        /// </summary>
        private void KickBall()
        {
            var diff = new Vector3(0, box.size.y / kickAngle, 0);
            var origin = transform.position - transform.TransformVector(diff);

            RaycastHit hit;
            if(Physics.Raycast(origin, transform.forward, out hit, kickDistance))
            {
                if(hit.collider.name == Ball.Name)
                {
                    var crb = hit.collider.GetComponent<Rigidbody>();
                    var force = -kickForce * hit.normal;
                    crb.AddForceAtPosition(force, hit.point, ForceMode.Impulse);
                }
            }
        }
    }
}