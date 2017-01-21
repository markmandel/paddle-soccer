using Client;
using UnityEngine;
using UnityEngine.Networking;
using Server;

namespace Network
{
    public class GameNetwork : NetworkManager, IUnityServer, IUnityClient
    {
        public readonly string Version = "0.2";

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
        /// <param name="port">The port to use</param>
        public void SetHost(string host, int port)
        {
            networkAddress = host;
            networkPort = port;
        }
    }
}