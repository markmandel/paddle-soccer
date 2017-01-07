using System;
using Client.Common;
using UnityEngine;

namespace Client.Controllers
{
    public class ScoreController : MonoBehaviour
    {
        private static readonly string PlayerOneGoal = "/Soccerfield/PlayerGoal.1";
        private static readonly string PlayerTwoGoal = "/Soccerfield/PlayerGoal.2";
        private static readonly string Ball = "Ball";

        private int playerOneScore;
        private int playerTwoScore;

        private void Start()
        {
            playerOneScore = 0;
            playerTwoScore = 0;

            var p1Observable = GameObject.Find(PlayerOneGoal).GetComponent<TriggerObservable>();
            p1Observable.TriggerEnter += OnBallGoal(_ => playerOneScore += 1);
            p1Observable.TriggerEnter += OnBallGoal(OnGoal);

            var p2Observable = GameObject.Find(PlayerTwoGoal).GetComponent<TriggerObservable>();
            p2Observable.TriggerEnter += OnBallGoal(_ => playerOneScore += 1);
            p2Observable.TriggerEnter += OnBallGoal(OnGoal);
        }

        // returns a event handler for the TriggerObservable that
        //only fires when the ball goes into the goal.
        private static TriggerObservable.Triggered OnBallGoal(Action<Collider> action)
        {
            return collider =>
            {
                if(collider.name == Ball)
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