using System;
using System.Linq;
using UnityEngine;

namespace Client
{
    public class GameClient
    {
        private static GameClient instance;
        private IUnityClient client;

        private static readonly string defaultHost = "localhost";
        private static readonly int defaultPort = 7777;
        private static readonly string hostArg = "-host";
        private static readonly string portArg = "-port";


        /// <summary>
        /// Returns the singleton instances if it exists.
        /// </summary>
        /// <exception cref="Exception">When Start() has not been called yet, this is thrown</exception>
        private static GameClient Instance
        {
            get
            {
                if (instance == null)
                {
                    throw new Exception(string.Format("{0} has not had Start() called.",
                        typeof(GameClient).FullName));
                }
                return instance;
            }
        }

        /// <summary>
        /// Start the network connection to the server.
        /// Looks at -host and -port arguments to override
        /// the default network server location.
        /// </summary>
        /// <param name="client">The client that will start the connection</param>
        /// <param name="commandLineArgs">The current command line arguments</param>
        /// <exception cref="Exception"></exception>
        public static void Start(IUnityClient client, string[] commandLineArgs)
        {
            if (instance == null)
            {
                instance = new GameClient {client = client};
                ProcessCommandLineArguments(commandLineArgs);
                Instance.client.StartClient();
            }
            else
            {
                throw new Exception(string.Format("{0} Can only be started once!",
                    typeof(GameClient).FullName));
            }
        }

        /// <summary>
        /// Processes the command line arguments, and applies the
        /// host and port settings
        /// </summary>
        /// <param name="commandLineArgs"></param>
        private static void ProcessCommandLineArguments(string[] commandLineArgs)
        {
            if (!commandLineArgs.Any())
            {
                Debug.Log("[GameClient] Default host and port");
                return;
            }

            var host = defaultHost;
            var port = defaultPort;

            for (var i = 0; i < commandLineArgs.Count(); i++)
            {
                var arg = commandLineArgs[i];
                if (arg == hostArg)
                {
                    host = commandLineArgs[i + 1];
                }
                else if (arg == portArg)
                {
                    port = int.Parse(commandLineArgs[i + 1]);
                }
            }

            Debug.LogFormat("[GameClient] Host: {0}, Port: {1}", host, port);
            Instance.client.SetHost(host, port);
        }

        /// <summary>
        /// Stop.
        /// </summary>
        public static void Stop()
        {
            instance = null;
        }
    }
}