using UnityEngine;

public class PaddleController : MonoBehaviour
{
    void OnCollisionEnter(Collision col)
    {
        Debug.Log("Collion! " + col.gameObject.ToString());
    }

    void Update()
    {
        // fixes some weird bug in the rotation
        transform.localEulerAngles = Vector3.zero;
    }
}