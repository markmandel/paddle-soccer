// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

﻿using System;
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