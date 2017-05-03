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

 using UnityEngine;

/// <summary>
/// Observable event for when collision triggers occur
/// on isTrigger Collisionable Game Objects
/// </summary>
[RequireComponent(typeof(Collider))]
public class TriggerObservable : MonoBehaviour
{
    /// <summary>
    /// Delegate for when a trigger event occurs
    /// </summary>
    /// <param name="other">The Collider we collide with</param>
    public delegate void Triggered(Collider other);

    /// <summary>
    /// Fired when the trigger enters
    /// </summary>
    public event Triggered TriggerEnter;

    /// <summary>
    /// Fires when the trigger exits
    /// </summary>
    public event Triggered TriggerExit;

    /// <summary>
    /// Fires when the trigger stays
    /// </summary>
    public event Triggered TriggerStay;

    /// <summary>
    /// Fires TriggerEnter event
    /// </summary>
    /// <param name="other"></param>
    private void OnTriggerEnter(Collider other)
    {
        if (TriggerEnter != null)
        {
            TriggerEnter(other);
        }
    }

    /// <summary>
    /// Fires TriggerExit event
    /// </summary>
    /// <param name="other"></param>
    private void OnTriggerExit(Collider other)
    {
        if (TriggerExit != null)
        {
            TriggerExit(other);
        }
    }

    /// <summary>
    /// Fires TriggerStay event
    /// </summary>
    /// <param name="other"></param>
    private void OnTriggerStay(Collider other)
    {
        if (TriggerStay != null)
        {
            TriggerStay(other);
        }
    }
}