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

using System;
using Game;
using UnityEngine;

namespace Server.Controllers
{
    /// <summary>
    /// Manages the score and goals for both players
    /// </summary>
    public class ScoreController : MonoBehaviour
    {
        [SerializeField]
        [Tooltip("Goal #1")]
        private GameObject goal1;

        [SerializeField]
        [Tooltip("Goal #2")]
        private GameObject goal2;

        // --- Messages ---

        private void OnValidate()
        {
            if (goal1 == null)
            {
                throw new Exception("[ScoreController] goal1 is null. This is bad");
            }

            if (goal2 == null)
            {
                throw new Exception("[ScoreController] goal2 is null. This is bad");
            }
        }

        /// <summary>
        /// Connect OnGameReady to connect players to their respective goals
        /// </summary>
        private void Start()
        {
            GameServer.OnGameReady += ConnectPlayerGoals;
        }

        // --- Functions --

        /// <summary>
        /// Connect Players to their respective Goals.
        /// </summary>
        private void ConnectPlayerGoals()
        {
            Debug.Log("Connecting Players to Goals!");
            var players = GameServer.GetPlayers();

            var p1Score = players[0].GetComponent<PlayerScore>();
            p1Score.Name = "Player 1";
            p1Score.TargetGoal(goal2.GetComponent<TriggerObservable>());

            var p2Score = players[1].GetComponent<PlayerScore>();
            p2Score.Name = "Player 2";
            p2Score.TargetGoal(goal1.GetComponent<TriggerObservable>());
        }
    }
}