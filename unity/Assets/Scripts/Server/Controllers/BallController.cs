using System;
using Game;
using Game.Common;
using UnityEngine;
using UnityEngine.Networking;

namespace Server.Controllers
{
    /// <summary>
    /// Create a ball. Once there is a Goal, create a new
    /// ball at the given time frame.
    /// </summary>
    public class BallController : NetworkBehaviour
    {
        [SerializeField]
        [Tooltip("The soccer ball prefab")]
        private GameObject prefabBall;

        /// <summary>
        /// The current instance of the ball.
        /// </summary>
        private GameObject currentBall;

        /// <summary>
        /// property to ensure we don't try and create a ball while creating a ball
        /// </summary>
        private bool isGoal;

        // --- Messages ---

        /// <summary>
        /// Make sure there is a ball prefab
        /// </summary>
        /// <exception cref="Exception">If the prefab is null, throws an exception</exception>
        private void OnValidate()
        {
            if (prefabBall == null)
            {
                throw new Exception("[Ball Controller] Ball prefab needs to be populated");
            }
        }

        /// <summary>
        /// Call when two players have joined the game
        /// </summary>
        private void Start()
        {
            Debug.Log("[Ball Controller] Initialising...");
            isGoal = false;
            var p1Goal = Goals.FindPlayerOneGoal().GetComponent<TriggerObservable>();
            var p2Goal = Goals.FindPlayerTwoGoal().GetComponent<TriggerObservable>();

            p1Goal.TriggerEnter += OnGoal;
            p2Goal.TriggerEnter += OnGoal;

            GameServer.OnGameReady += CreateBall;
        }

        // --- Functions ---

        /// <summary>
        /// Create a ball after 5 seconds. Removes the old one if there is one.
        /// </summary>
        private void OnGoal(Collider _)
        {
            if (!isGoal)
            {
                isGoal = true;
                Invoke("CreateBall", 5);
            }
        }

        /// <summary>
        /// Creates the ball
        /// </summary>
        private void CreateBall()
        {
            Debug.Log("[Ball Controller] Creating a ball");

            if (currentBall != null)
            {
                Destroy(currentBall);
            }
            currentBall = Instantiate(prefabBall);
            currentBall.name = Ball.Name;

            NetworkServer.Spawn(currentBall);

            isGoal = false;
        }
    }
}