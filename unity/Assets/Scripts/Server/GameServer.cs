using System;
using UnityEngine;
using UnityEngine.Networking;
using Random = System.Random;

namespace Server
{
    /// <summary>
    /// Game erver logic for the standalone server.
    /// You really only want one of these.
    /// </summary>
    public class GameServer
    {
        private static readonly int defaultMinPort = 7000;
        private static readonly int defaultMaxPort = 8000;

        /// <summary>
        /// Minimum port in range server can start on
        /// </summary>
        public static readonly string MinPortEnv = "MIN_PORT";

        /// <summary>
        /// Maximum port in range the server can start on
        /// </summary>
        public static readonly string MaxPortEnv = "MAX_PORT";

        /// <summary>
        /// Environment variable for where the session services
        /// is. Set by Kubernetes.
        /// </summary>
        public static readonly string SessionsServiceEnv = "SESSIONS_SERVICE_HOST";

        /// <summary>
        /// Environment or the session name. Provided by K8s downward API
        /// </summary>
        public static readonly string SessionNameEnv = "SESSION_NAME";


        /// <summary>
        /// How many times to retry for an open port
        /// </summary>
        private static readonly int maxStartRetries = 10;

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

        private Random rnd;
        private int port;

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
        /// Constructor, sets the server
        /// </summary>
        /// <param name="server"></param>
        public GameServer(IUnityServer server)
        {
            rnd = new Random();
            this.server = server;
        }

        /// <summary>
        /// Starts the server with the dependencies this singleton needs. This should only ever be called once,
        /// as it will throw an exception if it get called again.
        /// </summary>
        public static void Start(IUnityServer server)
        {
            if (instance == null)
            {
                instance = new GameServer(server);
                for (var i = 0; i < maxStartRetries; i++)
                {
                    instance.SelectPort();
                    if (Instance.server.StartServer())
                    {
                        instance.Register();
                        return;
                    }
                }

                Stop();
                throw new Exception(string.Format("Error starting server after {0} retries", maxStartRetries));
            }
            else
            {
                throw new Exception(string.Format("{0} Can only be started once!",
                    typeof(GameServer).FullName));
            }
        }

        /// <summary>
        /// Sets the port to a random
        /// </summary>
        private void SelectPort()
        {
            var minPort = defaultMinPort;
            var maxPort = defaultMaxPort;
            var minPortStr = Environment.GetEnvironmentVariable(MinPortEnv);
            var maxPortStr = Environment.GetEnvironmentVariable(MaxPortEnv);

            if (minPortStr != null)
            {
                minPort = int.Parse(minPortStr);
            }

            if (maxPortStr != null)
            {
                maxPort = int.Parse(maxPortStr);
            }

            port = rnd.Next(minPort, maxPort);
            Debug.LogFormat("[GameServer] Attempting to start server on port: {0}", port);
            server.SetPort(port);
        }

        /// <summary>
        /// Register this server
        /// </summary>
        private void Register()
        {
            var registry = Environment.GetEnvironmentVariable(SessionsServiceEnv);

            if (registry == null)
            {
                Debug.LogFormat("[GameNetwork] No Session Registry environment variable set. Skipping.");
                return;
            }

            Debug.Log("Registering Server...");

            var session = new Session
            {
                id = Environment.GetEnvironmentVariable(SessionNameEnv),
                port = port
            };

            var host = "http://" + registry + "/register";
            server.PostHTTP(host, JsonUtility.ToJson(session));
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