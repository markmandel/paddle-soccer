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

/// <summary>
/// Simply writes out to console when a collision
/// occurs, and tells you with what.
/// </summary>
public class CollisionDebug : MonoBehaviour
{
    // --- Handlers ---

    void OnCollisionEnter(Collision collision)
    {
        Log(collision.gameObject, "Collision Enter");
    }

    void OnCollisionExit(Collision collision)
    {
        Log(collision.gameObject, "Collision Exit");
    }

    void OnCollisionStay(Collision collision)
    {
        Log(collision.gameObject, "Collision Stay");
    }

    void OnTriggerEnter(Collider other)
    {
        Log(other.gameObject, "Trigger Enter");
    }

    void OnTriggerExit(Collider other)
    {
        Log(other.gameObject, "Trigger Exit");
    }

    void OnTriggerStay(Collider other)
    {
        Log(other.gameObject, "Trigger Stay");
    }

    // --- Functions ---

    private void Log(GameObject obj, string category)
    {
        Debug.Log(string.Format("[{0}] {1} Event: {2}", name, category, obj.name));
    }
}