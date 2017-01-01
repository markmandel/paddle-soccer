using UnityEngine;


// Zero outs the angle on each update, to make sure
// the GameObject always faces the same direction
public class ZeroAngle : MonoBehaviour {

	// Update is called once per frame
    void Update()
    {
        // fixes some weird bug in the rotation
        transform.localEulerAngles = Vector3.zero;
    }
}
