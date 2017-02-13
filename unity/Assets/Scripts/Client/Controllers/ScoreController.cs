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

﻿using Game;
using UnityEngine;

namespace Client.Controllers
{
    /// <summary>
    /// Manages the score for both players
    /// </summary>
    public class ScoreController : MonoBehaviour
    {
        private int playerOneScore;
        private int playerTwoScore;

        // --- Messages ---

        /// <summary>
        /// Attaches handlers to the goals
        /// </summary>
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

        /// <summary>
        /// Debug the score. Currently only client side.
        /// TODO: Move this to server side.
        /// </summary>
        /// <param name="_"></param>
        private void OnGoal(Collider _)
        {
            Debug.Log("GOOOAAAAALL!!!");
            Debug.Log(string.Format("Player One Score: {0}", playerOneScore));
            Debug.Log(string.Format("Player Two Score: {0}", playerTwoScore));
        }
    }
}