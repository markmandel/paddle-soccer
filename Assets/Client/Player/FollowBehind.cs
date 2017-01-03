using UnityEngine;

namespace Client.Player
{
    // Follows at a distance behind a specified gameobject
    public class FollowBehind : MonoBehaviour
    {
        [Tooltip("The transform to follow")]
        public Transform target;

        [SerializeField]
        [Tooltip("Distance to follow from")]
        private float distance = 2.2f;

        void Update()
        {
            if(target != null)
            {
                float y = transform.position.y;
                Vector3 diff = new Vector3(0, 0, distance);
                diff = target.TransformDirection(diff);
                diff = target.position - diff;
                diff.Set(diff.x, y, diff.z);
                transform.position = diff;
            }
        }
    }
}