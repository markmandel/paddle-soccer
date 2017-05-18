using Game;
using UnityEngine;

namespace Server.Controllers
{
    /// <summary>
    /// Logic for ScoreController
    /// </summary>
    public class ScoreControllerLogic
    {
        private IScoreController controller;

        public ScoreControllerLogic(IScoreController controller)
        {
            this.controller = controller;
        }

        /// <summary>
        /// PlayerScore.Scored delegate such that when the scores reach
        /// 3, then shut down the server after 5 seconds
        /// </summary>
        /// <param name="score"></param>
        /// <param name="isLocalPlayer"></param>
        public void DisconnectOnWin(int score, bool isLocalPlayer)
        {
            Debug.LogFormat("[ScoreController] DisconnectOnWin: {0}, {1}, {2}", score, isLocalPlayer,
                PlayerScore.WinningScore);
            
            if (score >= PlayerScore.WinningScore)
            {
                Debug.LogFormat(
                    "[ScoreController] score of {0} is greater than the {1} winnier score - disconnecting after 5 seconds",
                    score, PlayerScore.WinningScore);
                controller.DelayedStopServer(5);
            }
        }
    }
}