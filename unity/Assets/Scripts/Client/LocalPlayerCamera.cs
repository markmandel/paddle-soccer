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
using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Hooks up the local player to the camera
    /// </summary>
    public class LocalPlayerCamera : NetworkBehaviour
    {

        /// <summary>
        /// Tag for the camera, for lookup
        /// </summary>
        private static readonly string cameraTag = "MainCamera";

        // --- Messages ---

        /// <summary>
        /// Hook up the camera when the local player joins
        /// </summary>
        public override void OnStartLocalPlayer()
        {
            base.OnStartLocalPlayer();

            Debug.Log("Setting camera to follow current player");
            var camera = GameObject.FindGameObjectWithTag(cameraTag);
            var follow = camera.GetComponent<FollowBehind>();
            follow.target = transform;
        }

        // --- Functions ---
    }
}