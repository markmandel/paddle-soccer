using Client.GUI;
using Game;
using UnityEngine;
using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Component whose job it is to wire the player up to
    /// things that are dependent on knowing when the player shows
    /// up on screen. Primarily used for Client side management.
    /// </summary>
    [RequireComponent(typeof(PlayerScore))]
    public class PlayerDependencies : NetworkBehaviour
    {
        private const string ClientControllerName = "Client Controller";

        // --- Messages ---

        /// <summary>
        /// Wires up dependencies.
        /// </summary>
        private void Start()
        {
            if (isClient)
            {
                var playerScore = GetComponent<PlayerScore>();
                // I know people don't like this approach, but can't find a better way?
                var scoreBoard = GameObject.Find(ClientControllerName).GetComponent<ScoreBoard>();
                playerScore.ScoreIncrease += scoreBoard.OnScore;
            }
        }

        /// <summary>
        /// On destroy, make sure to remove yourself from dependencies
        /// </summary>
        private void OnDestroy()
        {
            if (isClient)
            {
                var playerScore = GetComponent<PlayerScore>();
                var scoreBoard = GameObject.Find(ClientControllerName).GetComponent<ScoreBoard>();
                playerScore.ScoreIncrease -= scoreBoard.OnScore;
            }
        }

        // --- Functions ---
    }
}