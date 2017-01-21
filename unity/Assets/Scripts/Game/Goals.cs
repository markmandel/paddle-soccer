using System;
using Game.Common;
using UnityEngine;

namespace Game
{
    public class Goals
    {
        public static readonly string PlayerOneGoal = "/Soccerfield/PlayerGoal.1";
        public static readonly string PlayerTwoGoal = "/Soccerfield/PlayerGoal.2";

        public static GameObject FindPlayerTwoGoal()
        {
            return GameObject.Find(PlayerTwoGoal);
        }

        public static GameObject FindPlayerOneGoal()
        {
            return GameObject.Find(PlayerOneGoal);
        }

        // returns a event handler for the TriggerObservable that
        // only fires when the ball goes into the goal.
        public static TriggerObservable.Triggered OnBallGoal(Action<Collider> action)
        {
            return collider =>
            {
                if(collider.name == Ball.Name)
                {
                    action(collider);
                }
            };
        }
    }
}