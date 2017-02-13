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

﻿using UnityEngine;

namespace Client
{
    /// <summary>
    /// Follows at a distance behind a specified gameobject
    /// </summary>
    public class FollowBehind : MonoBehaviour
    {
        [Tooltip("The transform to follow")]
        public Transform target;

        [SerializeField]
        [Tooltip("Distance to follow from")]
        private float distance = 2.2f;

        /// <summary>
        /// Set the position of the camera behind the player on each frame
        /// </summary>
        private void Update()
        {
            if(target != null)
            {
                // maintain the y position
                var yPosition = transform.position.y;
                // maintain the x rotation
                var xRotation = transform.localEulerAngles.x;

                Vector3 diff = target.forward * distance;
                diff = target.position - diff;
                diff.Set(diff.x, yPosition, diff.z);
                transform.position = diff;

                transform.LookAt(target);
                transform.rotation = Quaternion.Euler(xRotation, transform.localEulerAngles.y, transform.localEulerAngles.z);
            }
        }
    }
}