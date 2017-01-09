using System;
using Client.Common;
using Client.Game;
using Game;
using Game.Common;
using UnityEngine;

namespace Client.Controllers
{
    // Manages the score for both players
    public class ScoreController : MonoBehaviour
    {
        private int playerOneScore;
        private int playerTwoScore;

        // --- Messages ---

        private void Start()
        {
            playerOneScore = 0;
            playerTwoScore = 0;

            var p1Goal = Goals.FindPlayerOneGoal().GetComponent<TriggerObservable>();
            p1Goal.TriggerEnter += Goals.OnBallGoal(_ => playerOneScore += 1);
            p1Goal.TriggerEnter += Goals.OnBallGoal(OnGoal);

            var p2Goal = Goals.FindPlayerTwoGoal().GetComponent<TriggerObservable>();
            p2Goal.TriggerEnter += Goals.OnBallGoal(_ => playerOneScore += 1);
            p2Goal.TriggerEnter += Goals.OnBallGoal(OnGoal);
        }

        // --- Functions --

        private void OnGoal(Collider _)
        {
            Debug.Log("GOOOAAAAALL!!!");
            Debug.Log(string.Format("Player One Score: {0}", playerOneScore));
            Debug.Log(string.Format("Player Two Score: {0}", playerTwoScore));
        }
    }
}