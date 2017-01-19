using UnityEngine;
using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Hooks up the loca player to the camera
    /// </summary>
    public class LocalPlayerCamera : NetworkBehaviour
    {
        private static readonly string cameraTag = "MainCamera";

        // --- Messages ---

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