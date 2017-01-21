using System;
using Network;
using UnityEngine;
using UnityEngine.Networking;

namespace Server
{
    /// <summary>
    /// Game erver logic for the standalone server.
    /// You really only want one of these.
    /// </summary>
    public class GameServer
    {
        private static readonly int playersNeededForGame = 2;

        /// <summary>
        /// Delegate for when two players have joined the game
        /// </summary>
        public delegate void Ready();

        /// <summary>
        /// Event when two players have joined the game
        /// and we are ready to start
        /// </summary>
        public static event Ready OnGameReady;

        private static GameServer instance;

        // local member variables
        private int connCount;

        private IUnityServer server;

        /// <summary>
        /// Returns the singleton instances if it exists.
        /// </summary>
        /// <exception cref="Exception">When Start() has not been called yet, this is thrown</exception>
        private static GameServer Instance
        {
            get
            {
                if (instance == null)
                {
                    throw new Exception(string.Format("{0} has not had Start() called.",
                        typeof(GameServer).FullName));
                }
                return instance;
            }
        }

        /// <summary>
        /// Starts the server with the dependencies this singleton needs. This should only ever be called once,
        /// as it will throw an exception if it get called again.
        /// </summary>
        public static void Start(IUnityServer server)
        {
            if (instance == null)
            {
                instance = new GameServer {server = server};
                if (!Instance.server.StartServer())
                {
                    instance = null;
                    throw new Exception("Error starting server");
                }
            }
            else
            {
                throw new Exception(string.Format("{0} Can only be started once!",
                    typeof(GameServer).FullName));
            }
        }

        /// <summary>
        /// Stops the server, and lets go of all the resources.
        /// </summary>
        public static void Stop()
        {
            OnGameReady = null;
            instance = null;
        }

        /// <summary>
        /// Should be called when the server recieves a connection
        /// </summary>
        /// <param name="conn">The connection</param>
        public static void OnServerConnect(NetworkConnection conn)
        {
            Instance.connCount++;
            Debug.LogFormat("[GameNetwork] Client #{0} Connected", Instance.connCount);

            // only two players are allowed
            if (Instance.connCount > playersNeededForGame)
            {
                conn.Disconnect();
            }
        }

        /// <summary>
        /// Should be called when the server has a player added
        /// </summary>
        /// <param name="playerCount">Current number of players</param>
        public static void OnServerAddPlayer(int playerCount)
        {
            Debug.LogFormat("[GameNetwork] Adding Player {0}", playerCount);

            if (playerCount == playersNeededForGame && OnGameReady != null)
            {
                Debug.Log("[GameNetwork] Firing on game ready!");
                OnGameReady();
            }
        }
    }
}