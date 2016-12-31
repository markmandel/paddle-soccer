using UnityEngine;

public class PaddleCollision : MonoBehaviour
{
    void OnCollisionEnter(Collision col)
    {
        Debug.Log("Collion! " + col.gameObject.ToString());
    }
}