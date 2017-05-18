using System;
using UnityEngine;
using UnityEngine.UI;

namespace Client.GUI
{
    /// <summary>
    /// Components to take information from players
    /// and display it on the HUD score
    /// </summary>
    public class ScoreBoard : MonoBehaviour
    {
        private ScoreBoardLogic logic;

        [SerializeField]
        [Tooltip("Your score field")]
        private Text yourScore;

        [SerializeField]
        [Tooltip("Your opponent's score field")]
        private Text opponentScore;

        [SerializeField]
        [Tooltip("The center message box")]
        private Text centerMessage;

        // --- Messages ---

        private void Start()
        {
            logic = new ScoreBoardLogic(yourScore, opponentScore, centerMessage);
        }

        private void OnValidate()
        {
            if (yourScore == null)
            {
                throw new Exception("Your score is null. It should not be");
            }
            if (opponentScore == null)
            {
                throw new Exception("Opponent Score is null. It should not be");
            }
            if (centerMessage == null)
            {
                throw new Exception("Center Message is null. It should not be");
            }
        }

        // --- Functions ---
        
        /// <summary>
        /// Delegate to ScoreBoardLogic.OnScore
        /// </summary>
        /// <param name="score"></param>
        /// <param name="isPlayerLocal"></param>
        public void OnScore(int score, bool isPlayerLocal)
        {
            logic.OnScore(score, isPlayerLocal);
        }
    }
}