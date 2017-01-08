using System;
using Client.Common;
using Client.Game;
using UnityEngine;

namespace Client.Controllers
{
    // Manages the score for both players
    public class ScoreController : MonoBehaviour
    {
        public static readonly string PlayerOneGoal = "/Soccerfield/PlayerGoal.1";
        public static readonly string PlayerTwoGoal = "/Soccerfield/PlayerGoal.2";

        private int playerOneScore;
        private int playerTwoScore;

        // --- Messages ---

        private void Start()
        {
            playerOneScore = 0;
            playerTwoScore = 0;

            var p1Goal = GameObject.Find(PlayerOneGoal).GetComponent<TriggerObservable>();
            p1Goal.TriggerEnter += OnBallGoal(_ => playerOneScore += 1);
            p1Goal.TriggerEnter += OnBallGoal(OnGoal);

            var p2Goal = GameObject.Find(PlayerTwoGoal).GetComponent<TriggerObservable>();
            p2Goal.TriggerEnter += OnBallGoal(_ => playerOneScore += 1);
            p2Goal.TriggerEnter += OnBallGoal(OnGoal);
        }

        // --- Functions --

        // returns a event handler for the TriggerObservable that
        // only fires when the ball goes into the goal.
        private static TriggerObservable.Triggered OnBallGoal(Action<Collider> action)
        {
            return collider =>
            {
                if(collider.name == Ball.Name)
                {
                    action(collider);
                }
            };
        }

        private void OnGoal(Collider _)
        {
            Debug.Log("GOOOAAAAALL!!!");
            Debug.Log(string.Format("Player One Score: {0}", playerOneScore));
            Debug.Log(string.Format("Player Two Score: {0}", playerTwoScore));
        }
    }
}