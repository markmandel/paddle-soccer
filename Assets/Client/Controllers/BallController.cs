using System;
using Client.Common;
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

            var p1Goal = Goals.FindPlayerOneGoal().GetComponent<TriggerObservable>();
            var p2Goal = Goals.FindPlayerTwoGoal().GetComponent<TriggerObservable>();

            p1Goal.TriggerEnter += OnGoal;
            p2Goal.TriggerEnter += OnGoal;
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

        private void OnGoal(Collider _)
        {
            Invoke("CreateBall", 5);
        }

    }
}