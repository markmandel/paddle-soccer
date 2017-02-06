using System;
using UnityEngine;

namespace Game
{
    /// <summary>
    /// Utility for managing what happens to Goals.
    /// </summary>
    public static class Goals
    {
        private static readonly string PlayerOneGoal = "/Soccerfield/PlayerGoal.1";
        private static readonly string PlayerTwoGoal = "/Soccerfield/PlayerGoal.2";

        public static GameObject FindPlayerTwoGoal()
        {
            return GameObject.Find(PlayerTwoGoal);
        }

        public static GameObject FindPlayerOneGoal()
        {
            return GameObject.Find(PlayerOneGoal);
        }

        /// <summary>
        /// Rseturns a event handler for the TriggerObservable that
        /// only fires when the ball goes into the goal.
        /// </summary>
        /// <param name="action"></param>
        /// <returns></returns>
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