using UnityEngine;
using UnityEngine.Networking;

namespace Game
{
    /// <summary>
    /// Component to manage the score of a goal going
    /// State is coordinated via a syncvar
    /// </summary>
    public class PlayerScore : NetworkBehaviour
    {
        /// <summary>
        /// Maximum score a player can get to
        /// </summary>
        public const int WinningScore = 3;

        /// <summary>
        /// Synced field for storing the current score of a goal
        /// </summary>
        [SyncVar(hook = "ScoreHook")]
        private int _score;

        /// <summary>
        /// Player name - makes debugging easier.
        /// </summary>
        public string Name;

        /// <summary>
        /// Delegate for when the score has changed
        /// </summary>
        /// <param name="score"></param>
        /// <param name="isLocalPlayer">Whether this is the local player or not</param>
        public delegate void Scored(int score, bool isLocalPlayer);

        /// <summary>
        /// Fired when a score has changed
        /// </summary>
        public event Scored ScoreChange;

        /// <summary>
        /// Current score. Only the server can set the value
        /// Synced via Syncvar internally from server->client
        /// Will fire ScoreChange on the server side
        /// </summary>
        public int Score
        {
            get { return _score; }
            set
            {
                if (isServer)
                {
                    _score = value;
                    ScoreHook(_score);
                }
            }
        }

        // --- Messages ---

        /// <summary>
        /// Start sets the score to 0
        /// </summary>
        private void Start()
        {
            _score = 0;
        }

        // --- Functions ---

        /// <summary>
        /// Tell this Player what their target goal is
        /// </summary>
        /// <param name="observable">The TriggerObservable of the goal that the player should aim for</param>
        public void TargetGoal(TriggerObservable observable)
        {
            observable.TriggerEnter += Goals.OnBallGoal(_ => Score += 1);
            observable.TriggerEnter += Goals.OnBallGoal(
                _ => Debug.LogFormat("[PlayerScore] GOAL!!! {0}, Score: {1}", Name, Score));
        }

        /// <summary>
        /// Hook for when the score changes
        /// </summary>
        /// <param name="score"></param>
        private void ScoreHook(int score)
        {
            Debug.LogFormat("[PlayerScore] Score Hook: {0}", score);
            if (!isServer)
            {
                _score = score;
            }

            if (ScoreChange != null)
            {
                ScoreChange(_score, isLocalPlayer);
            }
        }
    }
}