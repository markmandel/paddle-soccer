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
using System.Collections;
using Game;
using UnityEngine;

namespace Server.Controllers
{
    /// <summary>
    /// Manages the score and goals for both players
    /// </summary>
    public class ScoreController : MonoBehaviour, IScoreController
    {
        [SerializeField]
        [Tooltip("Goal #1")]
        private GameObject goal1;

        [SerializeField]
        [Tooltip("Goal #2")]
        private GameObject goal2;

        private ScoreControllerLogic logic;

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
            logic = new ScoreControllerLogic(this);
            GameServer.OnGameReady += ConnectPlayerGoals;
            GameServer.OnGameReady += ConnectDisconnectOnWin;
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

        /// <summary>
        /// Connect up to each player score, so that the server will disconnect on winning
        /// </summary>
        private void ConnectDisconnectOnWin()
        {
            GameServer.GetPlayers().ForEach(x => x.GetComponent<PlayerScore>().ScoreChange += logic.DisconnectOnWin);
        }

        /// <summary>
        /// Stop the server after a number of seconds
        /// </summary>
        /// <param name="seconds"></param>
        public void DelayedStopServer(int seconds)
        {
            StartCoroutine(StopServer(seconds));
        }

        /// <summary>
        /// Coroutine for stopping the server, after a 5 second wait
        /// </summary>
        /// <param name="seconds"></param>
        /// <returns></returns>
        private IEnumerator StopServer(int seconds)
        {
            yield return new WaitForSeconds(seconds);
            GameServer.Stop();
        }
    }
}