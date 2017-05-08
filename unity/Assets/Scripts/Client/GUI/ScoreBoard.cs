using System;
using Game;
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
        [SerializeField]
        [Tooltip("Your score field")]
        private Text yourScore;

        [SerializeField]
        [Tooltip("Your opponent's score field")]
        private Text opponentScore;

        // --- Messages ---

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
        }

        // --- Functions ---

        public void OnScore(int score, bool isPlayerLocal)
        {
            Debug.LogFormat("[ScoreBoard] I gots a score! {0}, is local player: {1}", score, isPlayerLocal);
            if (isPlayerLocal)
            {
                yourScore.text = string.Format("You: {0}/{1}", score, PlayerScore.WinningScore);
            }
            else
            {
                opponentScore.text = string.Format("Them: {0}/{1}", score, PlayerScore.WinningScore);
            }
        }
    }
}