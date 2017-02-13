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

﻿using System;
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
        private static readonly string matchArg = "-match";


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
            }

            var host = defaultHost;
            var port = defaultPort;

            for (var i = 0; i < commandLineArgs.Count(); i++)
            {
                var arg = commandLineArgs[i];

                // Matchmaker takes precedence
                if (arg == matchArg)
                {
                    MatchMake(commandLineArgs[i + 1]);
                    return;
                }

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
            Instance.client.SetHost(host);
            Instance.client.SetPort(port);
            Instance.client.StartClient();
        }

        private static void MatchMake(string matchHost)
        {
            var path = matchHost + "/game";
            Debug.LogFormat("[GameClient] Invoking the MatchMaker! {0}", path);
            Instance.client.PostHTTP(path, null,
                request =>
                {
                    Debug.LogFormat("[GameClient] Matcher response: {0}:{1}", request.responseCode,
                        request.downloadHandler.text);

                    var game = JsonUtility.FromJson<Game>(request.downloadHandler.text);

                    if (request.responseCode == 201) //created and in queue
                    {
                        PollMatchMake(matchHost, game);
                    }
                    else if (request.responseCode == 200)
                    {
                        Instance.client.SetHost(game.ip);
                        Instance.client.SetPort(game.port);
                        Instance.client.StartClient();
                    }
                    else
                    {
                        throw new Exception("Error with matchmaker service");
                    }
                });
        }

        private static void PollMatchMake(string matchHost, Game game)
        {
            var path = matchHost + "/game/" + WWW.EscapeURL(game.id);

            Debug.LogFormat("[GameClient] Polling: {0}", path);

            Instance.client.PollGetHTTP(path, request =>
            {
                Debug.LogFormat("[GameClient] Polling Complete: {0}", request.downloadHandler.text);
                game = JsonUtility.FromJson<Game>(request.downloadHandler.text);

                // repeat if the game is still open
                if (game.status == 0)
                {
                    return false;
                }

                Instance.client.SetHost(game.ip);
                Instance.client.SetPort(game.port);
                Instance.client.StartClient();
                return true;

            });
        }

        /// <summary>
        /// Stop.
        /// </summary>
        public static void Stop()
        {
            instance = null;
        }

        /// <summary>
        /// Class for javascript deserialisation
        /// </summary>
        [Serializable]
        private class Game
        {
            public string id;
            public int status;
            public string sessionID;
            public int port;
            public string ip;
        }
    }
}