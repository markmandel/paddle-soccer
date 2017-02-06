using UnityEngine;
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