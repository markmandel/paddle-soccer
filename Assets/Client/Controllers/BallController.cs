using System;
using Client.Game;
using UnityEngine;

namespace Client.Controllers
{
    // Create a ball. Once there is a Goal, create a new
    // ball at the given time frame.
    public class BallController : MonoBehaviour
    {
        [SerializeField]
        [Tooltip("The soccer ball prefab")]
        private GameObject prefabBall;

        private GameObject currentBall;

        // --- Messages ---
        private void Start()
        {
            if(prefabBall == null)
            {
                throw new Exception("Prefab should not be null!");
            }

            CreateBall();
        }


        // --- Functions ---

        // Create a ball. Removes the old one if there is one.
        private void CreateBall()
        {
            if(currentBall != null)
            {
                Destroy(currentBall);
            }

            currentBall = Instantiate(prefabBall);
            currentBall.name = Ball.Name;
        }
    }
}