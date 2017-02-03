using System;
using System.Collections;
using System.Text;
using Client;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.Networking;
using Server;

namespace Network
{
    public class GameNetwork : NetworkManager, IUnityServer, IUnityClient
    {
        public readonly string Version = "0.3";

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

            if (PlayerInfo.IsHeadless())
            {
                Debug.Log("[GameNetwork] Starting Server");
                GameServer.Start(this);
            }
            else
            {
                Debug.Log("[GameNetwork] Starting Client");
                GameClient.Start(this, System.Environment.GetCommandLineArgs());
            }
        }

        // --- Server Commands ---
        public override void OnServerConnect(NetworkConnection conn)
        {
            base.OnServerConnect(conn);
            GameServer.OnServerConnect(conn);
        }

        public override void OnStopServer()
        {
            base.OnStopServer();
            GameServer.Stop();
        }

        public override void OnServerAddPlayer(NetworkConnection conn, short playerControllerId)
        {
            base.OnServerAddPlayer(conn, playerControllerId);
            GameServer.OnServerAddPlayer(numPlayers);
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

        public void PollGetHTTP(string host, Func<UnityWebRequest,bool> lambda)
        {
            StartCoroutine(AsyncPollGetHTTP(host, lambda));
        }

        private IEnumerator AsyncPollGetHTTP(string host, Func<UnityWebRequest,bool> action)
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
                    var success = action(get);

                    if (!success)
                    {
                        yield return new WaitForSeconds(2);
                        StartCoroutine(AsyncPollGetHTTP(host, action));
                    }
                }
            }
        }

        // --- Client & Server Commands ---

        public void PostHTTP(string host, string body, Action<UnityWebRequest> action = null)
        {
            StartCoroutine(AsyncPostHTTP(host, body, action));
        }

        private IEnumerator AsyncPostHTTP(string host, string body, Action<UnityWebRequest> action)
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
                    if (action != null)
                    {
                        action(post);
                    }
                }
            }
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