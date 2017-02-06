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