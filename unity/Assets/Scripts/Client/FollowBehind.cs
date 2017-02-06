using UnityEngine;

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