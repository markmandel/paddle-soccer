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
using System.Text;
using Client;
using Server;
using UnityEngine;
using UnityEngine.Networking;

namespace Network
{
    /// <summary>
    /// The overall network manager. Delegates as much work as possible to
    /// GameClient and GameServer.
    /// </summary>
    public class GameNetwork : NetworkManager, IUnityServer, IUnityClient
    {
        public readonly string Version = "0.1";

        /// <summary>
        /// How many players have joined the game?
        /// </summary>
        private int connCount;

        /// <summary>
        /// Starts eiter a client or a server - depending on if headless or not.
        /// </summary>
        private void Start()
        {
            Debug.LogFormat("[GameNetwork] Starting Client or Server? {0}", Version);

            if (UnityInfo.IsHeadless())
            {
                Debug.Log("[GameNetwork] Starting Server");
                GameServer.Start(this);
            }
            else
            {
                Debug.Log("[GameNetwork] Starting Client");
                GameClient.Start(this, Environment.GetCommandLineArgs());
            }
        }

        // --- Server Commands ---

        /// <summary>
        /// Run when the server is started. Delegates to GameServer.OnServerConnect
        /// </summary>
        /// <param name="conn"></param>
        public override void OnServerConnect(NetworkConnection conn)
        {
            base.OnServerConnect(conn);
            GameServer.OnServerConnect(conn);
        }

        /// <summary>
        /// Run when the server recieves a new player. Delegates to GameServer.OnServerAddPlayer
        /// </summary>
        /// <param name="conn"></param>
        /// <param name="playerControllerId"></param>
        public override void OnServerAddPlayer(NetworkConnection conn, short playerControllerId)
        {
            base.OnServerAddPlayer(conn, playerControllerId);
            var playerController = conn.playerControllers[playerControllerId];
            GameServer.OnServerAddPlayer(playerController.gameObject);
        }

        // --- Client Commands ---

        /// <summary>
        /// Change the Server Host settings from the default
        /// as set in the Unity editor.
        /// </summary>
        /// <param name="host">The server host</param>
        public void SetHost(string host)
        {
            networkAddress = host;
        }

        /// <summary>
        /// Asyncronously polls a HTTP endpoint every 2 seconds, until lambda returns true
        /// UnityWebRequest is passed to lambda on each invocation
        /// </summary>
        /// <param name="host">The host url to post to</param>
        /// <param name="lambda">lambda called on each request</param>
        public void PollGetHTTP(string host, Func<UnityWebRequest, bool> lambda)
        {
            StartCoroutine(AsyncPollGetHTTP(host, lambda));
        }

        /// <summary>
        /// Implementation of asyncronous polling of a HTTP endpoint
        /// </summary>
        /// <param name="host">the host url to make the GET Request</param>
        /// <param name="lambda">The lambda that is called on completion</param>
        private IEnumerator AsyncPollGetHTTP(string host, Func<UnityWebRequest, bool> lambda)
        {
            Debug.LogFormat("[GameNetwork] Getting data: {0}", host);

            using (var get = UnityWebRequest.Get(host))
            {
                yield return get.Send();

                if (get.isError)
                {
                    Debug.Log(get.error);
                }
                else
                {
                    Debug.Log("[GameNetwork] Get Complete");
                    var success = lambda(get);

                    if (!success)
                    {
                        yield return new WaitForSeconds(2);
                        StartCoroutine(AsyncPollGetHTTP(host, lambda));
                    }
                }
            }
        }

        // --- Client & Server Commands ---

        /// <summary>
        /// Asyncronously calls the HTTP POST host with the attached body.
        /// Calls lambda on successful post
        /// </summary>
        /// <param name="host">The host url to call</param>
        /// <param name="body">Body string to send as the POST body</param>
        /// <param name="lambda">Lambda called on successful post</param>
        public void PostHTTP(string host, string body, Action<UnityWebRequest> lambda = null)
        {
            StartCoroutine(AsyncPostHTTP(host, body, lambda));
        }

        /// <summary>
        /// Implementation of asyncronous call to POST HTTP
        /// </summary>
        /// <param name="host">The host url to call</param>
        /// <param name="body">Body string to send as the POST body</param>
        /// <param name="lambda">Lambda called on successful post</param>
        private IEnumerator AsyncPostHTTP(string host, string body, Action<UnityWebRequest> lambda)
        {
            Debug.LogFormat("[GameNetwork] Posting to: {0} data: {1}", host, body);

            // give it a null payload - because easy.
            if (body == null)
            {
                body = "{}";
            }

            var bytes = Encoding.UTF8.GetBytes(body);

            using (var post = UnityWebRequest.Put(host, bytes))
            {
                //switch it back to being post
                post.method = UnityWebRequest.kHttpVerbPOST;

                yield return post.Send();

                if (post.isError)
                {
                    Debug.Log(post.error);
                }
                else
                {
                    Debug.Log("[GameNetwork] Post complete!");
                    if (lambda != null)
                    {
                        lambda(post);
                    }
                }
            }
        }

        /// <summary>
        /// Close all connections, and then exits the servers
        /// </summary>
        public void Shutdown()
        {
            NetworkManager.Shutdown();
            Application.Quit();
        }

        /// <summary>
        /// Change the Server/Client port settings from the default
        /// as set in the Unity editor.
        /// </summary>
        /// <param name="port">The port to use</param>
        public void SetPort(int port)
        {
            networkPort = port;
        }
    }
}