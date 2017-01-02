using UnityEngine;

namespace Assets.Client.Common
{
    // Zero outs the angle on each update, to make sure
    // the GameObject always faces the same direction
    public class ZeroAngle : MonoBehaviour
    {
        // Update is called once per frame
        void Update()
        {
            // fixes a weird issue I was having with paddle rotation
            transform.localEulerAngles = Vector3.zero;
        }
    }
}